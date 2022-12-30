package main

import (
	"embed"
	"io/fs"

	"github.com/partyhall/partyhall/config"
	"github.com/partyhall/partyhall/logs"
	"github.com/partyhall/partyhall/services"
)

//go:embed gui/dist
var webapp embed.FS

//go:embed sql
var dbScripts embed.FS

const DEBUG = true

func init() {
	if !config.IsInDev() || DEBUG {
		subfs, err := fs.Sub(webapp, "gui/dist")
		if err != nil {
			logs.Warn("Failed to get webapp path. Not loading the webapp", err)
		} else {
			services.WEBAPP_FS = &subfs
		}
	}

	services.DB_SCRIPTS_FS = dbScripts
}
