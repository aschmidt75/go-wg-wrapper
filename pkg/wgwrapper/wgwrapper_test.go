// +build linux

package wgwrapper

import (
	"fmt"
	"math/rand"
	"net"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")

func newWGIntf() WireguardInterface {
	b := make([]rune, 4)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return NewWireguardInterface(fmt.Sprintf("wgtst-%s", string(b)), net.IPNet{
		IP:   net.IPv4(10, 99, 99, 99),
		Mask: net.CIDRMask(32, 32),
	})
}

func TestAddDeleteHasInterface(t *testing.T) {
	wg := New()
	wgi := newWGIntf()

	ex, err := wg.HasInterface(wgi)
	if err != nil {
		t.Errorf("Unable to execute HasInterface:  %s", err)
	}
	if ex {
		t.Errorf("Interface already exists but should not")
	}

	err = wg.AddInterface(wgi)
	if err != nil {
		t.Fatalf("Unable to execute AddInterface:  %s", err)
	}

	ex, err = wg.HasInterface(wgi)
	if err != nil {
		t.Errorf("Unable to execute HasInterface:  %s", err)
	}
	if !ex {
		t.Errorf("Interface does not exist but shoulds")
	}

	err = wg.DeleteInterface(wgi)
	if err != nil {
		t.Fatalf("Unable to execute DeleteInterface:  %s", err)
	}

	ex, err = wg.HasInterface(wgi)
	if err != nil {
		t.Errorf("Unable to execute HasInterface:  %s", err)
	}
	if ex {
		t.Errorf("Interface still exists but should not")
	}

}

func TestConfigure(t *testing.T) {
	wg := New()
	wgi := newWGIntf()
	wgiNonEx := newWGIntf()

	err := wg.AddInterface(wgi)
	if err != nil {
		t.Fatalf("Unable to execute AddInterface:  %s", err)
	}

	err = wg.Configure(&wgiNonEx)
	if err == nil {
		t.Error("Configure should fail on a nonexisting interface")
	}

	err = wg.Configure(&wgi)
	if err == nil {
		t.Error("Configure should fail without a listening port")
	}

	wgi.ListenPort = 46534
	err = wg.Configure(&wgi)
	if err != nil {
		t.Errorf("Unable to execute Configure:  %s", err)
	}
	if wgi.PublicKey == "" {
		t.Error("Configure should set a public key but did not")
	}

	err = wg.SetInterfaceUp(wgi)
	if err != nil {
		t.Fatalf("Unable to execute SetInterfaceUp:  %s", err)
	}

	err = wg.DeleteInterface(wgi)
	if err != nil {
		t.Fatalf("Unable to execute DeleteInterface:  %s", err)
	}

	err = wg.DeleteInterface(wgiNonEx)
	if err == nil {
		t.Fatalf("DeleteInterface on nonexisting interface succeeds but should not")
	}
}
