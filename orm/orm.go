package orm

import (
	"github.com/jmoiron/sqlx"
	"github.com/oxodao/photobooth/migrations"
	"github.com/oxodao/photobooth/utils"

	_ "github.com/mattn/go-sqlite3"
)

var GET *ORM

type ORM struct {
	DB       *sqlx.DB
	AppState AppState
	Events   Events
}

func Load() error {
	db := sqlx.MustConnect("sqlite3", utils.GetPath("photobooth.db"))

	GET = &ORM{
		DB:       db,
		AppState: AppState{db},
		Events:   Events{db},
	}

	err := migrations.DoMigrations(db)
	if err != nil {
		return err
	}

	return nil
}
