// +build linux

package wgwrapper

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

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

// DefaultRouteInterface returns the interface name of the default route.
func (wg wgwrapper) DefaultRouteInterface() (string, error) {
	//
	cmd := exec.Command("/sbin/ip", "route", "show", "default")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	if len(errStr) > 0 {
		e := fmt.Sprintf("/sbin/ip reported: %s", errStr)
		return "", errors.New(e)
	}
	a := strings.Split(outStr, " ")
	b := false
	for _, p := range a {
		if b {
			return p, nil
		}
		if p == "dev" {
			b = true
		}
	}
	return "", nil
}
