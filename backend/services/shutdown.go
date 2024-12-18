package services

import (
	"os/exec"

	"github.com/partyhall/partyhall/config"
	"github.com/partyhall/partyhall/log"
)

func Shutdown() error {
	if config.GET.IsInDev {
		log.Info("[IN DEV] Shutdown requested")

		return nil
	}

	err := DB.Close()
	if err != nil {
		return err
	}

	err = exec.Command("shutdown", "-h", "now").Run()

	return err
}
