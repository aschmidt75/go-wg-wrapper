// +build linux

package wgwrapper

import (
	"math/rand"
	"net"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestPeers(t *testing.T) {
	wg := New()
	wgi := newWGIntf()
	wgiNonEx := newWGIntf()

	err := wg.AddInterface(wgi)
	if err != nil {
		t.Fatalf("Unable to execute AddInterface:  %s", err)
	}

	wgi.ListenPort = 46534
	err = wg.Configure(&wgi)
	if err != nil {
		t.Errorf("Unable to execute Configure:  %s", err)
	}
	if wgi.PublicKey == "" {
		t.Error("Configure should set a public key but did not")
	}

	_, ipv4Net, err := net.ParseCIDR("127.0.0.1/32")

	wgp1 := WireguardPeer{
		RemoteEndpointIP: "10.1.2.3",
		ListenPort:       43210,
		Pubkey:           "9g4Eec+u+wBuMF06+qnsYl3G81l2PNCnG7nvtss9O2I=",
		AllowedIPs: []net.IPNet{
			*ipv4Net,
		},
		Psk: nil,
	}
	wgp2 := WireguardPeer{
		RemoteEndpointIP: "10.4.5.6",
		ListenPort:       32345,
		Pubkey:           "xqr+unDSDc5Fq0W9Zp2SJlzr+wOaFAquNdIMwPLHarw=",
		AllowedIPs: []net.IPNet{
			*ipv4Net,
		},
		Psk: nil,
	}

	ok, err := wg.HasPeer(wgi, wgp1)
	if err != nil {
		t.Errorf("Unable to execute HasPeer: %s", err)
	}
	if ok {
		t.Errorf("Was able to find test peer, but should not have been")
	}

	ok, err = wg.HasPeer(wgiNonEx, wgp1)
	if (err == nil) || ok {
		t.Errorf("HasPeer on nonexisting interface succeeds but should not.")
	}

	ok, err = wg.AddPeer(wgi, wgp1)
	if err != nil {
		t.Errorf("Unable to execute AddPeer: %s", err)
	}
	if !ok {
		t.Errorf("Was not able to add peer, but should have been")
	}

	ok, err = wg.AddPeer(wgi, wgp2)
	if err != nil {
		t.Errorf("Unable to execute AddPeer: %s", err)
	}
	if !ok {
		t.Errorf("Was not able to add peer, but should have been")
	}

	ok, err = wg.HasPeer(wgi, wgp1)
	if err != nil {
		t.Errorf("Unable to execute HasPeer: %s", err)
	}
	if !ok {
		t.Errorf("Was not able to find peer just added, but should have been")
	}

	ok, err = wg.HasPeer(wgi, wgp2)
	if err != nil {
		t.Errorf("Unable to execute HasPeer: %s", err)
	}
	if !ok {
		t.Errorf("Was not able to find peer just added, but should have been")
	}

	iterateCount := 0
	wg.IteratePeers(wgi, func(p WireguardPeer) {
		if p.Pubkey == "9g4Eec+u+wBuMF06+qnsYl3G81l2PNCnG7nvtss9O2I=" {
			iterateCount++
		} else {
			if p.Pubkey == "xqr+unDSDc5Fq0W9Zp2SJlzr+wOaFAquNdIMwPLHarw=" {
				iterateCount++
			} else {
				t.Errorf("Unexpected callback while iterating peers")
			}
		}
	})

	if iterateCount != 2 {
		t.Errorf("Expected to find two peers while iterating but did not.")
	}

	err = wg.RemovePeerByPubkey(wgi, "xqr+unDSDc5Fq0W9Zp2SJlzr+wOaFAquNdIMwPLHarw=")
	if err != nil {
		t.Errorf("Unable to execute RemovePeerByPubkey: %s", err)
	}

	ok, err = wg.HasPeer(wgi, wgp2)
	if err != nil {
		t.Errorf("Unable to execute HasPeer: %s", err)
	}
	if ok {
		t.Errorf("Was able to find 2nd peer just added, but should not have been")
	}

	err = wg.RemoveAllPeers(wgi)
	if err != nil {
		t.Errorf("Unable to execute RemoveAllPeers: %s", err)
	}

	ok, err = wg.HasPeer(wgi, wgp1)
	if err != nil {
		t.Errorf("Unable to execute HasPeer: %s", err)
	}
	if ok {
		t.Errorf("Was able to find test peer, but should not have been")
	}

	err = wg.DeleteInterface(wgi)
	if err != nil {
		t.Fatalf("Unable to execute DeleteInterface:  %s", err)
	}
}
