package modules

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/labstack/echo/v4"
	"github.com/partyhall/easyws"
	"github.com/partyhall/partyhall/config"
	"github.com/partyhall/partyhall/logs"
	"github.com/partyhall/partyhall/models"
	"github.com/partyhall/partyhall/modules/module_karaoke"
	"github.com/partyhall/partyhall/modules/module_photobooth"
	"github.com/partyhall/partyhall/remote"
	"github.com/partyhall/partyhall/services"
	"github.com/partyhall/partyhall/utils"
)

var MODULES map[string]Module

type Module interface {
	GetModuleName() string
	LoadConfig(filename string) error
	NewExporter(basePath string, event *models.Event) utils.Exporter

	PreInitialize() error
	Initialize() error
	GetMqttHandlers() map[string]mqtt.MessageHandler
	GetWebsocketHandlers() []easyws.MessageHandler
	UpdateFrontendSettings()
	RegisterApiRoutes(*echo.Group)
}

func LoadModules() error {
	// Ugly temporary hack to bypass golang's cycle issue
	services.RequestModuleExport = ExportEvent

	MODULES = map[string]Module{
		"photobooth": module_photobooth.ModulePhotobooth{},
		"karaoke": module_karaoke.ModuleKaraoke{
			VolumeInstru: 1,
			VolumeVocals: .2,
			VolumeFull:   .2,
		},
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

func PreInitializeModules() {
	for name, module := range MODULES {
		logs.Info("Pre-initializing module " + name)
		err := module.PreInitialize()
		if err != nil {
			panic(err)
		}
	}
}

func InitializeModules() {
	for name, module := range MODULES {
		logs.Info("Initializing module " + name)
		err := module.Initialize()
		if err != nil {
			panic(err)
		}

		module.UpdateFrontendSettings()
	}
}

// @TODO do this better so that module updates their settings themselves
func BroadcastFrontendSettings() {
	for _, module := range MODULES {
		module.UpdateFrontendSettings()
	}

	remote.BroadcastState()
}

func NormalizeModuleName(m Module) string {
	return strings.ToLower(m.GetModuleName())
}

func RegisterRoutes(g *echo.Group) {
	for _, module := range MODULES {
		sr := g.Group("/" + NormalizeModuleName(module))
		module.RegisterApiRoutes(sr)
	}
}

func ExportEvent(basePath string, event *models.Event) (map[string]any, error) {
	exportErrors := []error{}
	metadata := map[string]any{}

	for name, module := range MODULES {
		logs.Info("Exporting module " + name + "...")

		outpath := filepath.Join(basePath, name)
		if _, err := os.Stat(outpath); os.IsNotExist(err) {
			os.MkdirAll(outpath, os.ModePerm)
		}

		ex := module.NewExporter(outpath, event)
		m, err := ex.Export()
		if err != nil {
			exportErrors = append(exportErrors, err)
			continue
		}

		metadata[name] = m
	}

	if len(exportErrors) > 0 {
		return metadata, errors.Join(exportErrors...)
	}

	return metadata, nil
}
