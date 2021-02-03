// +build linux

package wgwrapper

type WireguardPeerIterator func(p WireguardPeer)

// WireguardWrapper is the main interface to work
// with wireguard interfaces
type WireguardWrapper interface {

	// AddInterface creates a wireguard interface from basic properties
	// (name, ip)
	AddInterface(intf WireguardInterface) error

	// AddInterfaceNoAddr is similar to AddInterface with the exception that no
	// IP address is added to the interface
	AddInterfaceNoAddr(intf WireguardInterface) error

	// DeleteInterface downs and deletes the interface
	DeleteInterface(intf WireguardInterface) error

	// SetInterfaceUp brings interface in UP state
	SetInterfaceUp(intf WireguardInterface) error

	// HasInterface checks if given interface exists (by name)
	HasInterface(intf WireguardInterface) (bool, error)

	// Configure makes sure that the wireguard interface
	// has a listen port configured and a keypair (by creatig one).
	// Needs endpoint ip and listen port from intf.
	// Extracts public key part and stores it in intf.
	Configure(intf *WireguardInterface) error

	// AddPeer adds a new peer to an existing interface
	AddPeer(intf WireguardInterface, peer WireguardPeer) (bool, error)

	// HasPeer check if a peer is present on an interface. Compares by public key only
	HasPeer(intf WireguardInterface, peer WireguardPeer) (bool, error)

	// RemoveAllPeers removes all peers on an existing interface
	RemoveAllPeers(intf WireguardInterface) error

	// RemovePeerByPubkey remove a single peer from an interface
	RemovePeerByPubkey(intf WireguardInterface, pubkey string) error

	// IteratePeers walks over the current list of peers of an interface
	IteratePeers(intf WireguardInterface, it WireguardPeerIterator) error
}

// New sets up a new WireguardWrapper
func New() WireguardWrapper {
	return wgwrapper{}
}
