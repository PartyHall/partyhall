package modules

import (
	"errors"
	"fmt"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gorilla/mux"
	"github.com/partyhall/easyws"
	"github.com/partyhall/partyhall/config"
	"github.com/partyhall/partyhall/logs"
	"github.com/partyhall/partyhall/modules/module_karaoke"
	"github.com/partyhall/partyhall/modules/module_photobooth"
	"github.com/partyhall/partyhall/remote"
	"github.com/partyhall/partyhall/utils"
)

var MODULES map[string]Module

type Module interface {
	GetModuleName() string
	LoadConfig(filename string) error

	Initialize() error
	GetMqttHandlers() map[string]mqtt.MessageHandler
	GetWebsocketHandlers() []easyws.MessageHandler
	UpdateFrontendSettings()
	RegisterApiRoutes(router *mux.Router)
}

func LoadModules() error {
	MODULES = map[string]Module{
		"photobooth": module_photobooth.ModulePhotobooth{},
		"karaoke":    module_karaoke.ModuleKaraoke{},
	}

	var loadedModules = 0

	for _, loadableModule := range config.GET.Modules {
		name := strings.ToLower(loadableModule)
		module, ok := MODULES[name]
		if !ok {
			return fmt.Errorf("failed to load module %v: does not exists", name)
		}

		logs.Info("Loading module " + module.GetModuleName())

		configFile := utils.GetPath("config", name+".yaml")

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

func InitializeModules() {
	for name, module := range MODULES {
		logs.Info("Initializing module " + name)
		err := module.Initialize()
		if err != nil {
			panic(err)
		}
	}
}

// @TODO do this better so that module updates their settings themselves
func BroadcastFrontendSettings() {
	for _, module := range MODULES {
		module.UpdateFrontendSettings()
	}

	remote.BroadcastState()
}

// @TODO remove spaces & not url valid characters
func NormalizeModuleName(m Module) string {
	return strings.ToLower(m.GetModuleName())
}

func RegisterRoutes(r *mux.Router) {
	for _, module := range MODULES {
		sr := r.PathPrefix("/" + NormalizeModuleName(module) + "/").Subrouter()
		module.RegisterApiRoutes(sr)
	}
}
