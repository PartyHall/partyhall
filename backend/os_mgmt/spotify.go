package os_mgmt

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

/**
 * Maybe at some point we need to use coreos/go-systemd
 * to have a better handling of everything
 * and maybe to have realtime stats about the services
 * but meh, it goes through dbus and I'm fed up with it
 * since I tried to make pulseaudio work so not for now
 **/

func SetSpotifySettings(enabled bool, deviceName string) error {
	spotifyConf := fmt.Sprintf(`
[global]
device_type="s_t_b"
device_name="%s"
`, deviceName)

	homeDir, _ := os.UserHomeDir()
	configDir := filepath.Join(homeDir, ".config", "spotifyd")

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	configPath := filepath.Join(configDir, "spotifyd.conf")
	if err := os.WriteFile(configPath, []byte(spotifyConf), 0644); err != nil {
		return err
	}

	if enabled {
		enableCmd := exec.Command("systemctl", "--user", "enable", "spotifyd")
		if err := enableCmd.Run(); err != nil {
			return fmt.Errorf("failed to enable spotifyd: %w", err)
		}

		restartCmd := exec.Command("systemctl", "--user", "restart", "spotifyd")
		if err := restartCmd.Run(); err != nil {
			return fmt.Errorf("failed to restart spotifyd: %w", err)
		}
	} else {
		disableCmd := exec.Command("systemctl", "--user", "disable", "--now", "spotifyd")
		if err := disableCmd.Run(); err != nil {
			return fmt.Errorf("failed to disable spotifyd: %w", err)
		}
	}

	/**
	 * @TODO: We need to check that the service is running properly / is stopped properly
	 * We don't want the UI to say "ok everything fine" when it is not
	 **/

	return nil
}
