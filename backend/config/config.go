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
	MosquittoAddr string `yaml:"mqtt_addr"`
	RootPath      string `yaml:"root_path"`
	GuestsAllowed bool   `yaml:"guests_allowed"`
	SendTime      bool   `yaml:"send_time"`

	NexusURL       string `yaml:"nexus_url"`
	NexusIgnoreSSL bool   `yaml:"nexus_ignore_ssl"`
	HardwareID     string `yaml:"hardware_id"`
	ApiKey         string `yaml:"api_key"`

	ModulesSettings ModulesSettings `yaml:"settings"`
	HardwareHandler HardwareHandler `yaml:"hardware_handler"`
}

type Mercure struct {
	SubscriberKey []byte
	PublisherKey  []byte
	ApplianceJWT  string
}

type HardwareHandler struct {
	BaudRate int               `yaml:"baud_rate"`
	Mappings map[string]string `yaml:"mappings"`
}

type ModulesSettings struct {
	Photobooth struct {
		Countdown  int `yaml:"countdown" json:"countdown"`
		Resolution struct {
			Width  int `yaml:"width" json:"width"`
			Height int `yaml:"height" json:"height"`
		} `yaml:"resolution" json:"resolution"`
		Unattended struct {
			Enabled  bool `yaml:"enabled" json:"enabled"`
			Interval int  `yaml:"interval" json:"interval"` // In seconds
		} `yaml:"unattended" json:"unattended"`
	} `yaml:"photobooth" json:"photobooth"`
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
