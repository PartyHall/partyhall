package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type UserSettings struct {
	Onboarded bool `yaml:"onboarded" json:"onboarded"`

	NexusURL       string `yaml:"nexus_url" json:"-"`
	NexusIgnoreSSL bool   `yaml:"nexus_ignore_ssl" json:"-"`
	HardwareID     string `yaml:"hardware_id" json:"hardware_id"`
	ApiKey         string `yaml:"api_key" json:"-"`

	Photobooth struct {
		Countdown       int `yaml:"countdown" json:"countdown"`
		FlashBrightness int `yaml:"flash_brightness" json:"flash_brightness"`
		Resolution      struct {
			Width  int `yaml:"width" json:"width"`
			Height int `yaml:"height" json:"height"`
		} `yaml:"resolution" json:"resolution"`
		Unattended struct {
			Enabled  bool `yaml:"enabled" json:"enabled"`
			Interval int  `yaml:"interval" json:"interval"` // In seconds
		} `yaml:"unattended" json:"unattended"`
	} `yaml:"photobooth" json:"photobooth"`

	Spotify struct {
		Enabled bool   `yaml:"enabled" json:"enabled"`
		Name    string `yaml:"name" json:"name"`
	} `yaml:"spotify" json:"-"`

	WirelessAp struct {
		Enabled  bool   `yaml:"enabled" json:"enabled"`
		Ssid     string `yaml:"ssid" json:"ssid"`
		Password string `yaml:"password" json:"password"`
	} `yaml:"wireless_ap" json:"-"`

	ButtonMappings map[string]string `yaml:"button_mappings" json:"-"`
}

func (us UserSettings) Save() error {
	bytes, err := yaml.Marshal(us)
	if err != nil {
		return err
	}

	return os.WriteFile(
		filepath.Join(GET.RootPath, "user_settings.yaml"),
		bytes,
		os.ModePerm,
	)
}
