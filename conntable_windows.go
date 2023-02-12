//go:build windows

package main

import (
	"log"
	"net/netip"
	"os"
	"os/exec"
)

func redirectPipeExecute(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func PostTunnelSetup(localNet *netip.Prefix, devName, dnsSuffix string) error {
	if err := redirectPipeExecute("netsh", "interface", "ip", "set", "address", "name="+devName, "source=static", "addr="+localNet.String(), "gateway=none"); err != nil {
		log.Printf("Configure tunnel failed.")
	}
	if err := redirectPipeExecute("netsh", "interface", "ipv4", "set", "dns", "name="+devName, "source=static", "addr="+localNet.Addr().String()); err != nil {
		log.Printf("Configure DNS failed.")
	}
	return nil
}

func PostTunnelCleanup(devName string) {
}
