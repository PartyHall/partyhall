package module_photobooth

import (
	"os"
	"path/filepath"

	"github.com/partyhall/partyhall/services"
	"github.com/partyhall/partyhall/utils"
)

func getModuleEventDir() (string, error) {
	evt := services.GET.CurrentState.CurrentEvent
	eventId := -1

	if evt != nil {
		eventId = *evt
	}

	basePath, err := utils.GetEventFolder(eventId)
	if err != nil {
		return "", err
	}

	fullPath := filepath.Join(basePath, "photobooth")
	err = os.MkdirAll(fullPath, os.ModePerm)

	return fullPath, err
}
