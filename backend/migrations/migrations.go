package migrations

import (
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Migration struct {
	Name        string    `db:"name"`
	Up          string    `db:"-"`
	Down        string    `db:"down_script"`
	MigrationTS int       `db:"migration_ts"`
	AppliedAt   time.Time `db:"applied_at"`
}

// @TODO: This should be in a separate library and the software using it
// should then just provide the migrations array

// @TODO: migration_ts should be PK

var migrations = []Migration{
	{
		Name:        "Migration table",
		Up:          `CREATE TABLE migrations(name VARCHAR(255), down_script TEXT, migration_ts BIGINT, applied_at TIMESTAMP);`,
		Down:        "",
		MigrationTS: 1725887251, // 2024-09-09 @ 15:07
	},
	{
		Name: "User",
		Up: `
			CREATE TABLE ph_user (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				username TEXT NOT NULL UNIQUE,
				name TEXT NULL,
				password TEXT NOT NULL,
				roles TEXT NOT NULL
			);

			CREATE TABLE refresh_token(
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				user_id INTEGER NOT NULL,
				token TEXT NOT NULL,
				expires_at TIMESTAMP NOT NULL,
				FOREIGN KEY (user_id) REFERENCES ph_user(id)
			);
		`,
		Down:        "",
		MigrationTS: 1726592411, // 2024-09-17 @ 19:00
	},
	{
		Name: "Events",
		Up: `
			CREATE TABLE event (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				name TEXT NOT NULL,
				author TEXT NULL,
				date DATETIME NULL,
				location TEXT NULL,
				amt_images_handtaken INT NOT NULL DEFAULT 0,
				amt_images_unattended INT NOT NULL DEFAULT 0,
				nexus_id VARCHAR(255) NULL DEFAULT NULL
			);

			CREATE TABLE picture (
				id TEXT PRIMARY KEY,
				taken_at DATETIME NOT NULL,
				unattended BOOLEAN NOT NULL,
				filename TEXT NOT NULL,
				event_id INTEGER NOT NULL REFERENCES event(id),
				nexus_id VARCHAR(255) NULL DEFAULT NULL
			);
		`,
		Down:        "",
		MigrationTS: 1726656973, // 2024-09-18 @ 12:56
	},
	{
		Name: "AppState",
		Up: `
			CREATE TABLE app_state (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				current_event INTEGER NULL REFERENCES event(id)
			);
		
			INSERT INTO app_state(current_event) VALUES (NULL);
		`,
		Down:        "",
		MigrationTS: 1726656998, // 2024-09-18 @ 12:57
	},
	{
		Name: "Logs",
		Up: `
			CREATE TABLE logs (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				type TEXT NOT NULL,
				text TEXT NOT NULL,
				timestamp DATETIME NOT NULL
			);
		`,
		Down:        "",
		MigrationTS: 1727716505, // 2024-09-30 @ 19:15
	},
	{
		Name: "Song",
		Up: `
			CREATE TABLE song (
				nexus_id VARCHAR(255) NOT NULL UNIQUE,
				title TEXT NOT NULL,
				artist TEXT NOT NULL,
				format VARCHAR(32) NOT NULL DEFAULT 'cdg',
				spotify_id VARCHAR(255) NULL,
				hotspot INTEGER NULL,

				has_cover BOOL NOT NULL DEFAULT FALSE,
				has_vocals BOOL NOT NULL DEFAULT FALSE,
				has_combined BOOL NOT NULL DEFAULT FALSE
			);
		
			CREATE TABLE song_session (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				event_id INTEGER NOT NULL REFERENCES event(id),
				nexus_id VARCHAR(255) NULL,
				title TEXT NOT NULL,
				artist TEXT NOT NULL,
				sung_by VARCHAR(255) NOT NULL,
				added_at datetime NULL DEFAULT NULL,
				started_at datetime NULL DEFAULT NULL,
				ended_at datetime NULL DEFAULT NULL,
				cancelled_at datetime NULL DEFAULT NULL
			);
		`,
		Down:        "",
		MigrationTS: 1731341120, // 2024-11-11 @ 17:05
	},
	{
		Name: "Song duration",
		Up: `
			ALTER TABLE song
			ADD COLUMN duration INTEGER NOT NULL DEFAULT 0
		`,
		Down:        "",
		MigrationTS: 1731871468, // 2024-11-17 @ 20:24
	},
	{
		Name: "SongSession Nexus ID",
		Up: `
			ALTER TABLE song_session
			ADD COLUMN session_nexus_id INTEGER NULL DEFAULT NULL
		`,
		Down:        "",
		MigrationTS: 1732467353, // 2024-11-24 @ 17:55
	},
	{
		Name: "Full text search for songs",
		Up: `
			CREATE VIRTUAL TABLE songs_fts USING fts5(
				title,
				artist,
				content='song',
				content_rowid='rowid',
				tokenize='unicode61 remove_diacritics 1'
			);

			INSERT INTO songs_fts(rowid, title, artist)
			SELECT rowid, title, artist FROM song;

			CREATE TRIGGER songs_ai AFTER INSERT ON song BEGIN
				INSERT INTO songs_fts(rowid, title, artist) VALUES (new.rowid, new.title, new.artist);
			END;

			CREATE TRIGGER songs_ad AFTER DELETE ON song BEGIN
				INSERT INTO songs_fts(songs_fts, rowid, title, artist) VALUES('delete', old.rowid, old.title, old.artist);
			END;

			CREATE TRIGGER songs_au AFTER UPDATE ON song BEGIN
				INSERT INTO songs_fts(songs_fts, rowid, title, artist) VALUES('delete', old.rowid, old.title, old.artist);
				INSERT INTO songs_fts(rowid, title, artist) VALUES (new.rowid, new.title, new.artist);
			END;
		`,
		Down:        "",
		MigrationTS: 1735405547, // 2024-12-28 @ 18:05
	},
	{
		Name: "Adding backdrops",
		Up: `
			CREATE TABLE backdrop_album (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				nexus_id INTEGER NULL DEFAULT NULL,
				name TEXT NOT NULL,
				author TEXT NOT NULL,
				version INTEGER NULL DEFAULT NULL
			);

			CREATE TABLE backdrop (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				album_id INTEGER NOT NULL REFERENCES backdrop_album(id),
				nexus_id INTEGER NULL DEFAULT NULL,
				title TEXT NOT NULL,
				filename VARCHAR(255) NULL
			);
		`,
		Down:        "",
		MigrationTS: 1742846926, // 2025-03-24 @ 21:09
	},
	{
		Name: "Adding alternate file",
		Up: `
			ALTER TABLE picture ADD COLUMN alternate_filename TEXT NULL DEFAULT NULL;
		`,
		Down:        "",
		MigrationTS: 1742849939, // 2025-03-24 @ 21:58
	},
}

func reverse(s []Migration) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func doApplyMigrations(logger *zap.SugaredLogger, db *sqlx.DB, up bool) error {
	for _, migration := range migrations {
		way := "up"
		sql := migration.Up
		if !up {
			sql = migration.Down
			way = "down"
		}

		logger.Infof(`- Migrating %v (%v)...`, migration.Name, way)
		_, err := db.Exec(sql)
		if err != nil {
			return err
		}

		if up {
			_, err = db.Exec(
				`INSERT INTO migrations (name, down_script, migration_ts, applied_at) VALUES ($1, $2, $3, $4)`,
				migration.Name,
				migration.Down,
				migration.MigrationTS,
				time.Now(),
			)

			if err != nil {
				return err
			}
		} else {
			_, err = db.Exec(
				`DELETE FROM migrations WHERE migration_ts = $1`,
				migration.MigrationTS,
			)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func ApplyMigrations(logger *zap.SugaredLogger, db *sqlx.DB, allowDowngrades bool) error {
	logger.Info("Checking for available migrations...")

	// Checking if the table "migrations" exists
	row := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND tbl_name='migrations';`)

	if row.Err() != nil {
		return row.Err()
	}

	var tableName string
	err := row.Scan(&tableName)
	if err != nil {
		// If the table does not exists at all (brand new DB)
		// We pass all the migrations
		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}

		logger.Warn("This database do not seems to have been initialized. Applying all migrations!")

		return doApplyMigrations(logger, db, true)
	}

	//#region Fetching currently applied migrations
	rows, err := db.Queryx(`
		SELECT name, down_script, migration_ts, applied_at
		FROM migrations
		ORDER BY migration_ts ASC
	`)

	if err != nil {
		return err
	}

	appliedMigrations := []Migration{}
	for rows.Next() {
		m := Migration{}
		err := rows.StructScan(&m)
		if err != nil {
			return err
		}

		appliedMigrations = append(appliedMigrations, m)
	}
	//#endregion

	amtAppliedMigrations := len(appliedMigrations)
	latestAvailableMigration := migrations[len(migrations)-1].MigrationTS

	if amtAppliedMigrations > 0 && appliedMigrations[amtAppliedMigrations-1].MigrationTS > migrations[len(migrations)-1].MigrationTS {
		logger.Warn("This database was made with a newer version of the software")

		if allowDowngrades {
			tmpMigrations := []Migration{}

			for _, m := range appliedMigrations {
				if m.MigrationTS > latestAvailableMigration {
					tmpMigrations = append(tmpMigrations, m)
				}
			}

			reverse(tmpMigrations)

			// return doApplyMigrations(logger, db, false)
			return errors.New("downgrade db is not implemented yet! (Well it is but it doesn't work @TODO)")
		} else {
			return errors.New("won't downgrade the database automatically as it could result in data loss. If you are sure, pass `--allow-downgrades` to the start command")
		}
	}

	// Filter already applied migrations and put the other in the migrations array
	tmpMigrations := []Migration{}
	latestMigration := appliedMigrations[amtAppliedMigrations-1].MigrationTS

	for _, m := range migrations {
		if m.MigrationTS > latestMigration {
			tmpMigrations = append(tmpMigrations, m)
		}
	}

	if len(tmpMigrations) == 0 {
		logger.Info("The database is up to date!")
	}

	migrations = tmpMigrations

	return doApplyMigrations(logger, db, true)
}
