package migrations

import (
	"database/sql"
	"embed"
	"fmt"
	"os"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/oxodao/photobooth/utils"
)

var MIGRATIONS = []migration{}

type migration interface {
	Apply(*sqlx.DB) error
	Revert(*sqlx.DB) error
}

func CheckDbExists(scripts embed.FS) error {
	path := utils.GetPath("photobooth.db")
	if _, err := os.Stat(path); !os.IsNotExist(err) {
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

	envHwid := os.Getenv("PHOTOBOOTH_HWID")
	envToken := os.Getenv("PHOTOBOOTH_TOKEN")
	token := ""
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
		fmt.Println("Found newer migrations, applying them")
		currentVersion := version
		lastVersion := len(MIGRATIONS)
		fmt.Printf("Current version: %v, Latest version: %v\n", currentVersion, lastVersion)
		for ; currentVersion < len(MIGRATIONS); currentVersion++ {
			fmt.Printf("\t- Applying migration %v\n", currentVersion)
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
		fmt.Println("Database is up to date !")
	}

	return nil
}
