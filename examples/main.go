package main

import (
	"fmt"
	"net"
	"os/exec"

	"github.com/aschmidt75/go-wg-wrapper/pkg/wgwrapper"
)

func main() {
	// get a new wrapper
	wg := wgwrapper.New()
	//fmt.Printf("%#v\n", wg)

	// set up a new wireguard interface struct w/ some defaults
	wgi := wgwrapper.NewWireguardInterface("wg-wrap-0", net.IPv4(10, 99, 99, 99))
	//fmt.Printf("%#v\n", wgi)

	// add the interface.
	err := wg.AddInterface(wgi)
	if err != nil {
		panic(err)
	}

	// Calling /sbin/ip link show yields an interface wg-wrap-0
	out, err := exec.Command("/sbin/ip", "-d", "link", "show", "dev", "wg-wrap-0").Output()
	println(string(out))

	// we're able to locate it
	ex, err := wg.HasInterface(wgi)
	if err != nil {
		panic(err)
	}

	// Configure the interface by creating a keypair and make it
	// listen on a port
	wgi.ListenPort = 46534
	err = wg.Configure(&wgi)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("%#v\n", wgi) // must have public key here and endpoint/listenport assigned

	err = wg.SetInterfaceUp(wgi)
	if err != nil {
		panic(err)
	}

	out, err = exec.Command("/sbin/ip", "addr", "show", "wg-wrap-0").Output()
	println(string(out))

	// add a sample peering to nowhere
	_, ipv4Net, err := net.ParseCIDR("127.0.0.1/32")

	_, err = wg.AddPeer(wgi, wgwrapper.WireguardPeer{
		RemoteEndpointIP: "10.1.2.3",
		ListenPort:       43210,
		Pubkey:           "9g4Eec+u+wBuMF06+qnsYl3G81l2PNCnG7nvtss9O2I=",
		AllowedIPs: []net.IPNet{
			*ipv4Net,
		},
		Psk: nil,
	})

	if err != nil {
		panic(err)
	}

	// calling /usr/bin/wg yields some output, including the peer. Assumes wireguard-tools is installed
	out, err = exec.Command("wg").Output()
	println(string(out))

	// we can iterate peer here, also
	wg.IteratePeers(wgi, func(p wgwrapper.WireguardPeer) {
		fmt.Printf("Peer: %s %s:%d\n", p.Pubkey, p.RemoteEndpointIP, p.ListenPort)
	})

	// delete the interface
	err = wg.DeleteInterface(wgi)
	if err != nil {
		panic(err)
	}

	ex, _ = wg.HasInterface(wgi)
	if ex {
		panic("Error: Interface exists but should have been deleted.")
	}
}
