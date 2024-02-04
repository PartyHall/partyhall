package migrations

import (
	"database/sql"
	"embed"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/partyhall/partyhall/logs"
	"github.com/partyhall/partyhall/utils"
)

var MIGRATIONS = []migration{}

type migration interface {
	Apply(*sqlx.DB) error
	Revert(*sqlx.DB) error
}

func CheckDbExists(scripts embed.FS) error {
	path := utils.GetPath("partyhall.db")
	if utils.FileExists(path) {
		return nil
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return err
	}

	initScriptData, err := scripts.ReadFile("sql/init.sql")
	if err != nil {
		return err
	}

	initScript := string(initScriptData)
	commentStrippedScript := ""
	for _, line := range strings.Split(initScript, "\n") {
		if strings.HasPrefix(line, "--") || len(line) == 0 {
			continue
		}

		commentStrippedScript += line + "\n"
	}

	for _, sqlCommand := range strings.Split(commentStrippedScript, ";\n") {
		_, err := db.Exec(sqlCommand + ";")
		if err != nil {
			return err
		}
	}

	envHwid := os.Getenv("PARTYHALL_HWID")
	envToken := os.Getenv("PARTYHALL_TOKEN")
	token := ""
	if len(envHwid) == 0 {
		envHwid = uuid.New().String()
	}

	if len(envToken) > 0 {
		token = envToken
	}

	_, err = db.Exec(`
			INSERT INTO app_state(hwid, token)
			VALUES (?, ?)
		`, envHwid, token)
	if err != nil {
		return err
	}

	return db.Close()
}

func DoMigrations(db *sqlx.DB) error {
	row := db.QueryRow(`SELECT last_applied_migration FROM app_state`)
	if row.Err() != nil {
		return row.Err()
	}

	var version int
	err := row.Scan(&version)
	if err != nil {
		return err
	}

	if len(MIGRATIONS) > version {
		logs.Info("Found newer migrations, applying them")
		currentVersion := version
		lastVersion := len(MIGRATIONS)
		logs.Infof("Current version: %v, Latest version: %v\n", currentVersion, lastVersion)
		for ; currentVersion < len(MIGRATIONS); currentVersion++ {
			logs.Infof("\t- Applying migration %v\n", currentVersion)
			err := MIGRATIONS[currentVersion].Apply(db)
			if err != nil {
				return err
			}

			_, err = db.Exec(`
				UPDATE app_state
				SET last_applied_migration = ?
			`, currentVersion+1)
			if err != nil {
				return err
			}
		}
	} else {
		logs.Info("Database is up to date !")
	}

	return nil
}
