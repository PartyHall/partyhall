package os_mgmt

import (
	"fmt"
	"os/exec"

	"github.com/partyhall/partyhall/log"
)

var (
	ErrInvalidParameters         = fmt.Errorf("the provided parameters are incorrect, make sure to provide an ethernet interface, a wifi interface, an ssid, and a password")
	ErrEthernetInterfaceNotFound = fmt.Errorf("the specified ethernet interface was not found, check the interface name")
	ErrWifiInterfaceNotFound     = fmt.Errorf("the specified wifi interface was not found, check the interface name")
	ErrInvalidHostapdConfig      = fmt.Errorf("the hostapd configuration is invalid, check the provided parameters")
	ErrInvalidSsid               = fmt.Errorf("the ssid is invalid, it must contain between 1 and 32 alphanumeric characters or spaces")
	ErrInvalidPassword           = fmt.Errorf("the password is invalid, it must contain between 12 and 63 alphanumeric characters")
	ErrHostapdRestartFailed      = fmt.Errorf("unable to restart the hostapd service, check the service logs for more details")
	ErrDnsmasqRestartFailed      = fmt.Errorf("unable to restart the dnsmasq service, check the service logs for more details")
	ErrAssignIpFailed            = fmt.Errorf("unable to assign an IP address to the wifi interface, check the network configuration")
	ErrUnknown                   = fmt.Errorf("an unknown error occurred while configuring the hotspot")
)

func SetHostapdConfig(
	ethIface string,
	wifiIface string,
	enabled bool,
	ssid string,
	password string,
) error {
	if !enabled {
		stopCmd := exec.Command("sudo", "systemctl", "stop", "hostapd")
		if err := stopCmd.Run(); err != nil {
			return fmt.Errorf("unable to stop the hostapd service, ensure the service is installed and active: %w", err)
		}
		return nil
	}

	log.Info("Setting new hotspot settings", "ethIface", ethIface, "wifiIface", wifiIface, "ssid", ssid, "password", password)

	setupCmd := exec.Command(
		"sudo",
		"-c",
		fmt.Sprintf("/usr/bin/setup-hotspot %s %s %s %s", ethIface, wifiIface, ssid, password),
	)

	// err := setupCmd.Run()

	out, err := setupCmd.CombinedOutput()

	if err != nil {
		log.Error("Failed to setup hotspot", "output", string(out))

		exitError, ok := err.(*exec.ExitError)
		if ok {
			switch exitError.ExitCode() {
			case 2:
				return ErrEthernetInterfaceNotFound
			case 3:
				return ErrWifiInterfaceNotFound
			case 4:
				return ErrInvalidHostapdConfig
			case 5:
				return ErrInvalidSsid
			case 6:
				return ErrInvalidPassword
			case 7:
				return ErrHostapdRestartFailed
			case 8:
				return ErrDnsmasqRestartFailed
			case 9:
				return ErrAssignIpFailed
			case 10:
				return ErrInvalidParameters
			default:
				return fmt.Errorf("%w: error code %d", ErrUnknown, exitError.ExitCode())
			}
		}

		return fmt.Errorf("error while executing the setup-hotspot script: %w", err)
	}

	return nil
}
