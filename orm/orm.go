package orm

import (
	"github.com/jmoiron/sqlx"
	"github.com/partyhall/partyhall/migrations"
	"github.com/partyhall/partyhall/utils"

	_ "github.com/mattn/go-sqlite3"
)

var GET *ORM

type ORM struct {
	DB       *sqlx.DB
	AppState AppState
	Events   Events
	Users    Users
}

func Load() error {
	db := sqlx.MustConnect("sqlite3", utils.GetPath("partyhall.db"))

	GET = &ORM{
		DB:       db,
		AppState: AppState{db},
		Events:   Events{db},
		Users:    Users{db},
	}

	err := migrations.DoMigrations(db)
	if err != nil {
		return err
	}

	return nil
}
