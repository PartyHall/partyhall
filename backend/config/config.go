package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const AMT_RESULTS_PER_PAGE = 25

var GET Config

type Config struct {
	IsInDev   bool    `yaml:"-"`
	EventPath string  `yaml:"-"`
	Mercure   Mercure `yaml:"-"`

	ListeningAddr string `yaml:"listening_addr"`
	RootPath      string `yaml:"root_path"`
	SendTime      bool   `yaml:"send_time"`
	MosquittoAddr string `yaml:"mqtt_addr"`
	GuestsAllowed bool   `yaml:"guests_allowed"`

	UserSettings UserSettings `yaml:"-"`
}

type Mercure struct {
	SubscriberKey []byte
	PublisherKey  []byte
	ApplianceJWT  string
}

func Load(isInDev bool) error {
	GET = Config{}

	configPath := os.Getenv("PARTYHALL_CONFIG_FILE")
	if len(configPath) == 0 {
		configPath = "/etc/partyhall.yaml"
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, &GET)
	if err != nil {
		return err
	}

	userSettingsPath := filepath.Join(GET.RootPath, "user_settings.yaml")
	if _, err := os.Stat(userSettingsPath); !os.IsNotExist(err) {
		userData, err := os.ReadFile(userSettingsPath)
		if err != nil {
			return err
		}

		err = yaml.Unmarshal(userData, &GET.UserSettings)
		if err != nil {
			return err
		}
	}

	GET.IsInDev = isInDev

	GET.EventPath = filepath.Join(GET.RootPath, "events")
	if _, err := os.Stat(GET.EventPath); os.IsNotExist(err) {
		if err := os.MkdirAll(GET.EventPath, os.ModePerm); err != nil {
			return err
		}
	}

	GET.Mercure.SubscriberKey, err = loadMercureKey(GET.RootPath, "subscriber")
	if err != nil {
		return err
	}

	GET.Mercure.PublisherKey, err = loadMercureKey(GET.RootPath, "publisher")
	if err != nil {
		return err
	}

	return nil
}
