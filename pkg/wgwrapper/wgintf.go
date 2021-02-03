// +build linux

package wgwrapper

import (
	"net"
)

// WireguardInterface wraps basic information about a wireguard interface
type WireguardInterface struct {
	InterfaceName string    // name of the wireguard interface
	IP            net.IPNet // local ip of wg interface
	ListenPort    int       // UDP listening port
	PublicKey     string    // public key of interface
}

// NewWireguardInterface creates a new WireguardInterface with a given name and ip
func NewWireguardInterface(interfaceName string, ipnet net.IPNet) WireguardInterface {
	return WireguardInterface{
		InterfaceName: interfaceName,
		IP:            ipnet,
	}
}

// NewWireguardInterface creates a new WireguardInterface with a given name but
// without an address assignement
func NewWireguardInterfaceNoAddr(interfaceName string) WireguardInterface {
	return WireguardInterface{
		InterfaceName: interfaceName,
	}
}
