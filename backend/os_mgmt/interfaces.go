package os_mgmt

import (
	"fmt"
	"net"
	"strings"
)

func isUselessInterface(iface string) bool {
	// We only care about ethernet and wifi
	for _, usefulIface := range []string{"wl", "en"} {
		if strings.HasPrefix(iface, usefulIface) {
			return false
		}
	}

	return true
}

type Interface struct {
	FriendlyName string   `json:"friendly_name"`
	Name         string   `json:"name"`
	Wireless     bool     `json:"wireless"`
	Ips          []string `json:"ips"`
}

type Interfaces struct {
	Ethernet []Interface `json:"ethernet"`
	Wifi     []Interface `json:"wifi"`
}

func FindInterfaces() Interfaces {
	interfaces := Interfaces{
		Ethernet: []Interface{},
		Wifi:     []Interface{},
	}

	ifaces, _ := net.Interfaces()

	for _, iface := range ifaces {
		if isUselessInterface(iface.Name) {
			continue
		}

		ifaceObj := Interface{
			Name: iface.Name,
			Ips:  []string{},
		}

		addresses, _ := iface.Addrs()
		for _, ip := range addresses {
			ipObj := net.ParseIP(strings.Split(ip.String(), "/")[0])

			if ipObj == nil {
				continue
			}

			ip4 := ipObj.To4()
			if ip4 == nil {
				continue
			}

			ifaceObj.Ips = append(ifaceObj.Ips, ip4.String())
		}

		ifaceObj.Wireless = strings.HasPrefix(ifaceObj.Name, "wl")

		if ifaceObj.Wireless {
			ifaceObj.FriendlyName = fmt.Sprintf("Wifi [%s]", ifaceObj.Name)
			interfaces.Wifi = append(interfaces.Wifi, ifaceObj)
		} else {
			ifaceObj.FriendlyName = fmt.Sprintf("Ethernet [%s]", ifaceObj.Name)
			interfaces.Ethernet = append(interfaces.Ethernet, ifaceObj)
		}
	}

	return interfaces
}
