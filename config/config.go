package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

var GET Config

const (
	MODE_PHOTOBOOTH = "PHOTOBOOTH"
	MODE_QUIZ       = "QUIZ"
	MODE_DISABLED   = "DISABLED"
)

var MODES = []string{
	MODE_PHOTOBOOTH,
	MODE_QUIZ,
	MODE_DISABLED,
}

type WebcamResolution struct {
	Width  int `json:"width" yaml:"width"`
	Height int `json:"height" yaml:"height"`
}

type PhotoboothConfig struct {
	HardwareFlash      bool             `yaml:"hardware_flash" json:"hardware_flash"`
	DefaultTimer       int              `yaml:"default_timer" json:"-"`
	UnattendedInterval int              `yaml:"unattended_interval" json:"-"`
	WebcamResolution   WebcamResolution `yaml:"webcam_resolution" json:"webcam_resolution"`
}

type MosquittoConfig struct {
	Address string `json:"address"`
}

type Config struct {
	Web struct {
		ListeningAddr string `yaml:"listening_addr"`
		AdminPassword string `yaml:"admin_password"`
	} `yaml:"web"`

	DebugMode bool `yaml:"debug_mode"`

	RootPath    string `yaml:"root_path"`
	DefaultMode string `yaml:"default_mode"`

	Mosquitto MosquittoConfig  `yaml:"mosquitto"`
	PartyHall PhotoboothConfig `yaml:"partyhall"`
}

func (c *Config) GetImageFolder(eventId int64, unattended bool) (string, error) {
	subfolder := "pictures"
	if unattended {
		subfolder = "unattended"
	}

	folderName := fmt.Sprintf("%v", eventId)
	if eventId < 0 {
		folderName = "NO_EVENT"
	}

	path := filepath.Join(c.RootPath, "images", folderName, subfolder)
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return "", err
	}
	return path, nil
}

func Load() error {
	cfg := Config{}

	configPath := os.Getenv("PARTYHALL_CONFIG_PATH")
	if len(configPath) == 0 {
		configPath = "/etc/partyhall.yaml"
	}

	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return err
	}

	if cfg.PartyHall.DefaultTimer == 0 {
		cfg.PartyHall.DefaultTimer = 3
	}

	if len(cfg.DefaultMode) == 0 {
		cfg.DefaultMode = MODE_PHOTOBOOTH
	}

	GET = cfg

	return nil
}
