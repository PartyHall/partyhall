package utils

import (
	"net"
	"strings"
)

func GetIPs() map[string][]string {
	ipAddresses := map[string][]string{}

	interfaces, _ := net.Interfaces()
	for _, inter := range interfaces {
		shouldSkip := false
		for _, ignored := range []string{"lo", "br-", "docker", "vmnet", "veth"} { // Ignoring docker / vmware networks for in-dev purposes
			if strings.HasPrefix(inter.Name, ignored) {
				shouldSkip = true
				break
			}
		}

		if shouldSkip {
			continue
		}

		ipAddresses[inter.Name] = []string{}

		addrs, _ := inter.Addrs()
		for _, ip := range addrs {
			ipAddresses[inter.Name] = append(ipAddresses[inter.Name], ip.String())
		}
	}

	return ipAddresses
}
