package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/partyhall/partyhall/utils"
	"go.bug.st/serial"
	"gopkg.in/yaml.v2"
)

var GET Config

const (
	MODE_PHOTOBOOTH = "PHOTOBOOTH"
	MODE_KARAOKE    = "KARAOKE"
	MODE_DISABLED   = "DISABLED"
)

var MODES = []string{
	MODE_PHOTOBOOTH,
	MODE_KARAOKE,
	MODE_DISABLED,
}

type MosquittoConfig struct {
	Address string `json:"address"`
}

type HardwareHandlerConfig struct {
	Mappings   map[string]string `yaml:"mappings"`
	BaudRate   int               `yaml:"baud_rate"`
	SerialMode *serial.Mode      `yaml:"-"`
}

type Config struct {
	Web struct {
		ListeningAddr string `yaml:"listening_addr"`
		AdminPassword string `yaml:"admin_password"`
	} `yaml:"web"`

	DebugMode bool `yaml:"debug_mode"`

	RootPath    string `yaml:"root_path"`
	DefaultMode string `yaml:"default_mode"`

	Mosquitto MosquittoConfig `yaml:"mosquitto"`

	HardwareHandler HardwareHandlerConfig `yaml:"hardware_handler"`

	SpotifyClientID     string `yaml:"spotify_client_id"`
	SpotifyClientSecret string `yaml:"spotify_client_secret"`

	Modules []string
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

	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return err
	}

	if len(cfg.DefaultMode) == 0 {
		cfg.DefaultMode = MODE_PHOTOBOOTH
	}

	cfg.HardwareHandler.SerialMode = &serial.Mode{
		BaudRate: cfg.HardwareHandler.BaudRate,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}

	GET = cfg

	utils.ROOT_PATH = cfg.RootPath

	return initializeFolders()
}

func initializeFolders() error {
	err := utils.MakeOrCreateFolder("")
	if err != nil {
		fmt.Println("Failed to create root folder")
		os.Exit(1)
	}

	for _, folder := range []string{"images", "config"} {
		if err := utils.MakeOrCreateFolder(folder); err != nil {
			return err
		}
	}

	return nil
}

func GetMqttTopic(module string, topic string) string {
	module = strings.ToLower(module)
	topic = strings.ToLower(topic)

	if len(module) == 0 {
		return fmt.Sprintf("partyhall/%v", topic)
	}

	return fmt.Sprintf("partyhall/%v/%v", module, topic)
}

func IsInDev() bool {
	return strings.HasPrefix(os.Args[0], "/tmp")
}
