// +build linux

package wgwrapper

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"os/exec"
	"strings"

	wgctrl "golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type wgwrapper struct {
	WireguardWrapper
}

// AddInterface adds a new wireguard interface
// by calling /sbin/ip. Also adds given address
func (wg wgwrapper) AddInterface(intf WireguardInterface) error {
	i, err := net.InterfaceByName(intf.InterfaceName)

	if i == nil || err != nil {
		// create wireguard interface
		var cmd *exec.Cmd

		cmd = exec.Command("/sbin/ip", "link", "add", "dev", intf.InterfaceName, "type", "wireguard")

		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			e := fmt.Sprintf("/sbin/ip reported: %s", err)
			return errors.New(e)
		}
		_, errStr := string(stdout.Bytes()), string(stderr.Bytes())
		if len(errStr) > 0 {
			e := fmt.Sprintf("/sbin/ip reported: %s", errStr)
			return errors.New(e)
		}
	}

	// try to get this intf via netlink
	i, err = net.InterfaceByName(intf.InterfaceName)
	if err != nil {
		return err
	}

	// Assign IP if desired and not yet present
	a, err := i.Addrs()
	if err != nil {
		return err
	}
	if len(a) == 0 {
		cmd := exec.Command("/sbin/ip", "address", "add", "dev", intf.InterfaceName, intf.IP.String())
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			return err
		}
		_, errStr := string(stdout.Bytes()), string(stderr.Bytes())
		if len(errStr) > 0 {
			e := fmt.Sprintf("/sbin/ip reported: %s", errStr)
			return errors.New(e)
		}
	}

	a, err = i.Addrs()
	if len(a) == 0 {
		e := fmt.Sprintf("unable to add ip address %s to interface %s: %s", intf.IP.String(), intf.InterfaceName, err)
		return errors.New(e)
	}

	return nil
}

// AddInterfaceNoAddr is similar to AddInterface with the exception that no
// IP address is added to the interface
func (wg wgwrapper) AddInterfaceNoAddr(intf WireguardInterface) error {
	i, err := net.InterfaceByName(intf.InterfaceName)

	if i == nil || err != nil {
		// create wireguard interface
		var cmd *exec.Cmd

		cmd = exec.Command("/sbin/ip", "link", "add", "dev", intf.InterfaceName, "type", "wireguard")

		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			e := fmt.Sprintf("/sbin/ip reported: %s", err)
			return errors.New(e)
		}
		_, errStr := string(stdout.Bytes()), string(stderr.Bytes())
		if len(errStr) > 0 {
			e := fmt.Sprintf("/sbin/ip reported: %s", errStr)
			return errors.New(e)
		}
	}

	// try to get this intf via netlink
	i, err = net.InterfaceByName(intf.InterfaceName)
	if err != nil {
		return err
	}

	return nil
}

func (wg wgwrapper) SetInterfaceUp(intf WireguardInterface) error {

	// check status
	cmd := exec.Command("/sbin/ip", "--br", "link", "show", "dev", intf.InterfaceName, "up", "type", "wireguard")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return err
	}
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	if len(errStr) > 0 {
		e := fmt.Sprintf("/sbin/ip reported: %s", errStr)
		return errors.New(e)
	}
	if len(outStr) > 0 {
		return nil // already up
	}

	// bring up wireguard interface
	cmd = exec.Command("/sbin/ip", "link", "set", "up", "dev", intf.InterfaceName)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		_, errStr = string(stdout.Bytes()), string(stderr.Bytes())
		if len(errStr) > 0 {
			e := fmt.Sprintf("/sbin/ip reported: %s", errStr)
			return errors.New(e)
		}

		return err
	}

	return nil

}

// DeleteInterface takes down an existing wireguard interface
// by calling /sbin/ip
func (wg wgwrapper) DeleteInterface(intf WireguardInterface) error {
	i, err := net.InterfaceByName(intf.InterfaceName)

	if err != nil {
		return err
	}
	if i == nil {
		e := fmt.Sprintf("No network/interface by name %s", intf.InterfaceName)
		return errors.New(e)
	}

	// take down wireguard interface
	var cmd *exec.Cmd

	cmd = exec.Command("/sbin/ip", "link", "set", "down", "dev", intf.InterfaceName)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		e := fmt.Sprintf("/sbin/ip reported: %s", err)
		return errors.New(e)
	}
	_, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	if len(errStr) > 0 {
		e := fmt.Sprintf("/sbin/ip reported: %s", errStr)
		return errors.New(e)
	}

	// remove wireguard interface
	cmd = exec.Command("/sbin/ip", "link", "delete", "dev", intf.InterfaceName, "type", "wireguard")

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		return err
	}
	_, errStr = string(stdout.Bytes()), string(stderr.Bytes())
	if len(errStr) > 0 {
		e := fmt.Sprintf("/sbin/ip reported: %s", errStr)
		return errors.New(e)
	}

	return nil
}

// HasInterface checks if the interface is already present
func (wg wgwrapper) HasInterface(intf WireguardInterface) (bool, error) {
	wgClient, err := wgctrl.New()
	if err != nil {
		return false, err
	}
	defer wgClient.Close()

	device, err := wgClient.Device(intf.InterfaceName)
	if err != nil {
		return false, nil
	}

	return (device != nil), nil
}

// Configure makes sure that the wireguard interface
// has a keypair and a listen port configured. Extracts public key
// part and stores it in intf.
func (wg wgwrapper) Configure(intf *WireguardInterface) error {
	// wireguard: create private key, add device (listen-port)
	wgClient, err := wgctrl.New()
	if err != nil {
		return err
	}
	defer wgClient.Close()

	wgDevice, err := wgClient.Device(intf.InterfaceName)
	if err != nil {
		return err
	}

	// check if device already has key and Listen port set. If not, do so
	if bytes.Compare(wgDevice.PrivateKey[:], emptyBytes32) == 0 {
		newKey, err := wgtypes.GeneratePrivateKey()
		if err != nil {
			return err
		}

		newConfig := wgtypes.Config{
			PrivateKey: &newKey,
		}
		err = wgClient.ConfigureDevice(intf.InterfaceName, newConfig)
		if err != nil {
			return err
		}
	}

	if wgDevice.ListenPort == 0 {
		if intf.ListenPort == 0 {
			return errors.New("wg listenPort may not be 0")
		}

		newConfig := wgtypes.Config{
			ListenPort: &intf.ListenPort,
		}
		err = wgClient.ConfigureDevice(intf.InterfaceName, newConfig)
		if err != nil {
			return err
		}
	}

	// query again make sure stuff is present
	wgDevice, err = wgClient.Device(intf.InterfaceName)
	if err != nil {
		return err
	}
	if wgDevice == nil {
		return errors.New("error reading wg device configuration")
	}

	if bytes.Compare(wgDevice.PrivateKey[:], emptyBytes32) == 0 || bytes.Compare(wgDevice.PublicKey[:], emptyBytes32) == 0 || wgDevice.ListenPort == 0 {
		return errors.New("unable to set wireguard key configuration")
	}

	intf.PublicKey = base64.StdEncoding.EncodeToString(wgDevice.PublicKey[:])

	return nil
}

// SetRoute checks if there is a route on given interface to network. If not, adds it. all using /sbin/ip
func (wg wgwrapper) SetRoute(intf WireguardInterface, networkCIDR string) error {
	//
	cmd := exec.Command("/sbin/ip", "route", "show", "dev", intf.InterfaceName)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return err
	}
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	if len(errStr) > 0 {
		e := fmt.Sprintf("/sbin/ip reported: %s", errStr)
		return errors.New(e)
	}
	a := strings.Split(outStr, " ")
	if len(a) > 0 {
		if a[0] == networkCIDR {
			// route is already present
			return nil
		}
	}

	//
	cmd = exec.Command("/sbin/ip", "route", "add", networkCIDR, "dev", intf.InterfaceName)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		_, errStr = string(stdout.Bytes()), string(stderr.Bytes())
		if len(errStr) > 0 {
			e := fmt.Sprintf("/sbin/ip reported: %s", errStr)
			return errors.New(e)
		}

		return err
	}

	return nil
}

var (
	emptyBytes32 = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
)
