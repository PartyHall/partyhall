package module_photobooth

import (
	"errors"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/labstack/echo/v4"
	"github.com/partyhall/easyws"
	"github.com/partyhall/partyhall/logs"
	"github.com/partyhall/partyhall/remote"
	"github.com/partyhall/partyhall/services"
	"github.com/partyhall/partyhall/utils"
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
	if !utils.FileExists(filename) {
		err := errors.New("Config file not found for module Photobooth")
		logs.Error(err.Error())

		return err
	}

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

func (m ModulePhotobooth) PreInitialize() error {
	remote.RegisterOnJoin("photobooth", func(socketType string, s *easyws.Socket) {
		if socketType == utils.SOCKET_TYPE_BOOTH {
			m.Actions.StartUnattended(s)
		}
	})

	return nil
}

func (m ModulePhotobooth) Initialize() error {
	return nil
}

func (m ModulePhotobooth) GetMqttHandlers() map[string]mqtt.MessageHandler {
	return map[string]mqtt.MessageHandler{
		"take_picture": OnTakePicture,
	}
}

func (m ModulePhotobooth) GetWebsocketHandlers() []easyws.MessageHandler {
	return []easyws.MessageHandler{
		TakePictureHandler{},
		RemoteTakePictureHandler{},
	}
}

func (m ModulePhotobooth) UpdateFrontendSettings() {
	services.GET.ModuleSettings["photobooth"] = map[string]interface{}{
		"unattended_interval": CONFIG.UnattendedInterval,
		"default_timer":       CONFIG.DefaultTimer,
		"hardware_flash":      CONFIG.HardwareFlash,
		"webcam_resolution":   CONFIG.WebcamResolution,
	}
}

func (m ModulePhotobooth) RegisterApiRoutes(g *echo.Group) {

}
