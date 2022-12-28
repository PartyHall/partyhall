package main

import (
	"embed"
	"io/fs"
	"os"
	"strings"

	"github.com/partyhall/partyhall/logs"
	"github.com/partyhall/partyhall/services"
)

//go:embed gui/dist
var webapp embed.FS

//go:embed gui_admin/dist
var adminapp embed.FS

//go:embed sql
var dbScripts embed.FS

func init() {
	execPath := os.Args[0]
	if !strings.HasPrefix(execPath, "/tmp/") {
		subfs, err := fs.Sub(webapp, "gui/dist")
		if err != nil {
			logs.Warn("Failed to get webapp path. Not loading the webapp", err)
		} else {
			services.WEBAPP_FS = &subfs
		}

		subfs2, err := fs.Sub(adminapp, "gui_admin/dist")
		if err != nil {
			logs.Warn("Failed to get adminapp path. Not loading the adminapp", err)
		} else {
			services.ADMIN_FS = &subfs2
		}
	}

	services.DB_SCRIPTS_FS = dbScripts
}
