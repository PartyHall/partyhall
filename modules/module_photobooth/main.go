package module_photobooth

import (
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"gopkg.in/yaml.v2"
)

var (
	INSTANCE = &ModulePhotobooth{
		Actions: Action{},
	}

	CONFIG = Config{}
)

type ModulePhotobooth struct {
	Actions Action
}

func (m ModulePhotobooth) GetModuleName() string {
	return "Photoobooth"
}

func (m ModulePhotobooth) LoadConfig(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	cfg := Config{}
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return err
	}

	if cfg.DefaultTimer == 0 {
		cfg.DefaultTimer = 3
	}

	CONFIG = cfg

	return nil
}

func (m ModulePhotobooth) GetMqttHandlers() map[string]mqtt.MessageHandler {
	return map[string]mqtt.MessageHandler{
		"take_picture": OnTakePicture,
	}
}

// func (m ModulePhotobooth) GetWebsocketHandlers() map[string] {
// 	return []message_handler.MessageHandler{}
// }

func (m ModulePhotobooth) GetFrontendSettings() map[string]interface{} {
	return map[string]interface{}{
		"unattended_interval": CONFIG.UnattendedInterval,
		"default_timer":       CONFIG.DefaultTimer,
		"hardware_flash":      CONFIG.HardwareFlash,
		"webcam_resolution":   CONFIG.WebcamResolution,
	}
}

func (m ModulePhotobooth) GetState() map[string]interface{} {
	return map[string]interface{}{}
}
