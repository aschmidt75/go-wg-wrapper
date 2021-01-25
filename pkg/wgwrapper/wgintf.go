// +build linux

package wgwrapper

import (
	"net"
)

// WireguardInterface wraps basic information about a wireguard interface
type WireguardInterface struct {
	InterfaceName string // name of the wireguard interface
	IP            net.IP // local ip of wg interface
	ListenPort    int    // UDP listening port
	PublicKey     string // public key of interface
}

// NewWireguardInterface creates a new WireguardInterface
func NewWireguardInterface(interfaceName string, ip net.IP) WireguardInterface {
	return WireguardInterface{
		InterfaceName: interfaceName,
		IP:            ip,
	}
}
