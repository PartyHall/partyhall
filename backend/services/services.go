package services

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/partyhall/partyhall/config"
	"github.com/partyhall/partyhall/log"
	"github.com/partyhall/partyhall/migrations"
	"github.com/partyhall/partyhall/models"
	"github.com/partyhall/partyhall/state"
)

var DB *sqlx.DB

func Load() error {
	// https://www.reddit.com/r/golang/comments/19b8eup/sqlite_and_timestamps/
	// The time_format make sqlite driver be not stupid
	db, err := sqlx.Connect("sqlite3", filepath.Join(config.GET.RootPath, "database.sqlite"))
	if err != nil {
		return err
	}

	// db.SetMaxOpenConns(1)
	DB = db
	log.DB = db

	err = migrations.ApplyMigrations(log.LOG, db, false)
	if err != nil {
		return err
	}

	return nil
}

func SetEvent(event *models.Event) error {
	if event != nil {
		folderPath := filepath.Join(config.GET.EventPath, fmt.Sprintf("%v", event.Id))

		for _, f := range []string{"karaoke", "photobooth"} {
			if err := os.MkdirAll(filepath.Join(folderPath, f), os.ModePerm); err != nil {
				return err
			}
		}
	}

	state.STATE.CurrentEvent = event

	return nil
}
