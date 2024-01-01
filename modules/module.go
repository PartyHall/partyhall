package modules

import (
	"errors"
	"fmt"
	"os"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/partyhall/easyws"
	"github.com/partyhall/partyhall/config"
	"github.com/partyhall/partyhall/modules/module_photobooth"
	"github.com/partyhall/partyhall/services"
	"github.com/partyhall/partyhall/utils"
)

var MODULES map[string]Module

type Module interface {
	GetModuleName() string
	LoadConfig(filename string) error

	GetMqttHandlers() map[string]mqtt.MessageHandler
	GetWebsocketHandlers() []easyws.MessageHandler
	GetState() map[string]interface{}
	GetFrontendSettings() map[string]interface{}
}

func LoadModules() error {
	MODULES = map[string]Module{
		"photobooth": module_photobooth.ModulePhotobooth{},
	}

	var loadedModules = 0

	for _, loadableModule := range config.GET.Modules {
		name := strings.ToLower(loadableModule)
		module, ok := MODULES[name]
		if !ok {
			return fmt.Errorf("failed to load module %v: does not exists", name)
		}

		configFile := utils.GetPath("config", name+".yaml")
		fmt.Println(configFile)
		if _, err := os.Stat(configFile); os.IsNotExist(err) {
			return fmt.Errorf("failed to load module %v: config file does not exists at path %v", name, configFile)
		}

		if err := module.LoadConfig(configFile); err != nil {
			return fmt.Errorf("failed to load module %v: %v", name, err)
		}

		loadedModules++
	}

	if loadedModules == 0 {
		return errors.New("no modules loaded, please add at least one")
	}

	return nil
}

func UpdateFrontendModuleSettings() {
	moduleSettings := map[string]interface{}{}

	for name, module := range MODULES {
		moduleSettings[name] = module.GetFrontendSettings()
	}

	services.GET.ModuleSettings = moduleSettings
}
