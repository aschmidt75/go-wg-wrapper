// +build linux

package wgwrapper

import (
	"encoding/base64"
	"fmt"
	"net"
	"time"

	wgctrl "golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

// WireguardPeer is a single wireguard peer
type WireguardPeer struct {
	RemoteEndpointIP            string
	ListenPort                  int
	Pubkey                      string
	AllowedIPs                  []net.IPNet
	Psk                         *string
	PersistentKeepaliveInterval time.Duration
}

// AddPeer adds a new peer to an existing interface
func (wg wgwrapper) AddPeer(intf WireguardInterface, peer WireguardPeer) (bool, error) {
	wgClient, err := wgctrl.New()
	if err != nil {
		return false, err
	}
	defer wgClient.Close()

	pk, err := wgtypes.ParseKey(peer.Pubkey)
	if err != nil {
		return false, err
	}

	wgDevice, err := wgClient.Device(intf.InterfaceName)
	if err != nil {
		return false, err
	}
	for _, p := range wgDevice.Peers {
		if p.PublicKey == pk {
			// Already present, skipping
			return false, nil
		}
	}

	var pskAsKey wgtypes.Key
	if peer.Psk != nil {
		pskAsKey, err = wgtypes.ParseKey(*peer.Psk)
		if err != nil {
			return false, err
		}
	}

	// process peer
	ep, err := net.ResolveUDPAddr("udp", net.JoinHostPort(peer.RemoteEndpointIP, fmt.Sprintf("%d", peer.ListenPort)))
	if err != nil {
		return false, err
	}

	newConfig := wgtypes.Config{
		ReplacePeers: false,
		Peers: []wgtypes.PeerConfig{
			wgtypes.PeerConfig{
				PublicKey:                   pk,
				Remove:                      false,
				PresharedKey:                &pskAsKey,
				Endpoint:                    ep,
				AllowedIPs:                  peer.AllowedIPs,
				PersistentKeepaliveInterval: &peer.PersistentKeepaliveInterval,
			},
		},
	}

	err = wgClient.ConfigureDevice(intf.InterfaceName, newConfig)
	if err != nil {
		return false, err
	}

	return true, nil
}

// HasPeer check if a peer is present on an interface. Compares by public key only
func (wg wgwrapper) HasPeer(intf WireguardInterface, peer WireguardPeer) (bool, error) {
	wgClient, err := wgctrl.New()
	if err != nil {
		return false, err
	}
	defer wgClient.Close()

	pk, err := wgtypes.ParseKey(peer.Pubkey)
	if err != nil {
		return false, err
	}

	wgDevice, err := wgClient.Device(intf.InterfaceName)
	if err != nil {
		return false, err
	}
	for _, p := range wgDevice.Peers {
		if p.PublicKey == pk {
			return true, nil
		}
	}

	return false, nil
}

// RemovePeerByPubkey remove a single peer from an interface
func (wg wgwrapper) RemovePeerByPubkey(intf WireguardInterface, pubkey string) error {
	wgClient, err := wgctrl.New()
	if err != nil {
		return err
	}
	defer wgClient.Close()

	// process peer
	pk, err := wgtypes.ParseKey(pubkey)
	if err != nil {
		return err
	}

	newConfig := wgtypes.Config{
		ReplacePeers: false,
		Peers: []wgtypes.PeerConfig{
			wgtypes.PeerConfig{
				PublicKey: pk,
				Remove:    true,
			},
		},
	}

	err = wgClient.ConfigureDevice(intf.InterfaceName, newConfig)
	if err != nil {
		return err
	}

	return nil
}

// RemoveAllPeers removes all peers on an existing interface
func (wg wgwrapper) RemoveAllPeers(intf WireguardInterface) error {
	wgClient, err := wgctrl.New()
	if err != nil {
		return err
	}
	defer wgClient.Close()

	newConfig := wgtypes.Config{
		ReplacePeers: true,
		Peers:        []wgtypes.PeerConfig{},
	}

	return wgClient.ConfigureDevice(intf.InterfaceName, newConfig)
}

// IteratePeers walks over the current list of peers of an interface
func (wg wgwrapper) IteratePeers(intf WireguardInterface, it WireguardPeerIterator) error {
	wgClient, err := wgctrl.New()
	if err != nil {
		return err
	}
	defer wgClient.Close()

	wgDevice, err := wgClient.Device(intf.InterfaceName)
	if err != nil {
		return err
	}
	for _, p := range wgDevice.Peers {
		it(WireguardPeer{
			RemoteEndpointIP: p.Endpoint.IP.String(),
			ListenPort:       p.Endpoint.Port,
			Pubkey:           base64.StdEncoding.EncodeToString(p.PublicKey[:]),
			AllowedIPs:       p.AllowedIPs,
			Psk:              nil, //p.PresharedKey.String(),
		})
	}

	return nil
}
