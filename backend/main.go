package main

import (
	"embed"
	"fmt"
	"os"
	"strings"

	"github.com/partyhall/partyhall/cmd"
	"github.com/partyhall/partyhall/config"
	"github.com/partyhall/partyhall/log"
	"github.com/partyhall/partyhall/state"
)

//go:embed frontend/app
var appFs embed.FS

//go:embed frontend/appliance
var applianceFs embed.FS

func main() {
	cmd.AppFS = appFs
	cmd.ApplianceFS = applianceFs

	isInDev := strings.TrimSpace(strings.ToLower(os.Getenv("PARTYHALL_ENV"))) == "dev"
	err := config.Load(isInDev)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	log.Load(isInDev)

	state.STATE.ModulesSettings = config.GET.ModulesSettings

	cmd.Execute()
}
