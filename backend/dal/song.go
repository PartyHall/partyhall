package dal

import (
	"fmt"
	"strings"

	"github.com/partyhall/partyhall/log"
	"github.com/partyhall/partyhall/models"
)

var SONGS Songs

type Songs struct{}

func (s Songs) GetAll() ([]models.Song, error) {
	rows, err := DB.Queryx(`
		SELECT
			nexus_id,
			title,
			artist,
			format,
			spotify_id,
			duration,
			hotspot,
			has_cover,
			has_vocals,
			has_combined
		FROM song
	`)

	if err != nil {
		return nil, err
	}

	songs := []models.Song{}

	for rows.Next() {
		event := models.Song{}

		err := rows.StructScan(&event)
		if err != nil {
			return nil, err
		}

		songs = append(songs, event)
	}

	return songs, nil
}

func (s Songs) GetCollection(search string, amt, offset int) (*models.PaginatedResponse, error) {
	resp := models.PaginatedResponse{}

	row := DB.QueryRow(`
		SELECT COUNT(*)
		FROM song
		WHERE LENGTH(?) = 0
		   OR LOWER(title) LIKE '%' || LOWER(?) || '%'
		   OR LOWER(artist) LIKE '%' || LOWER(?) || '%'
	`, search, search, search)
	if row.Err() != nil {
		return nil, row.Err()
	}

	err := row.Scan(&resp.TotalCount)
	if err != nil {
		return nil, err
	}

	resp.CalculateMaxPage()

	rows, err := DB.Queryx(`
		SELECT
			nexus_id,
			title,
			artist,
			format,
			spotify_id,
			duration,
			hotspot,
			has_cover,
			has_vocals,
			has_combined
		FROM song
		WHERE LENGTH(?) = 0
		   OR LOWER(title) LIKE '%' || LOWER(?) || '%'
		   OR LOWER(artist) LIKE '%' || LOWER(?) || '%'
		ORDER BY artist ASC, title ASC
		LIMIT ?
		OFFSET ?
	`, search, search, search, amt, offset)

	if err != nil {
		return nil, err
	}

	songs := []models.Song{}

	for rows.Next() {
		song := models.Song{}

		err := rows.StructScan(&song)
		if err != nil {
			return nil, err
		}

		songs = append(songs, song)
	}

	resp.Results = songs

	return &resp, nil
}

func (s Songs) Get(id string) (*models.Song, error) {
	row := DB.QueryRowx(`
		SELECT
			nexus_id,
			title,
			artist,
			format,
			spotify_id,
			duration,
			hotspot,
			has_cover,
			has_vocals,
			has_combined
		FROM song
		WHERE nexus_id = ?
	`, id)

	if row.Err() != nil {
		return nil, row.Err()
	}

	var song models.Song
	err := row.StructScan(&song)
	if err != nil {
		return nil, err
	}

	return &song, nil
}

func (s Songs) Create(song *models.Song) error {
	song.Format = strings.ToLower(song.Format)

	_, err := DB.NamedExec(`
		INSERT INTO song (
			nexus_id,
			title,
			artist,
			format,
			spotify_id,
			duration,
			hotspot,
			has_cover,
			has_vocals,
			has_combined
		)
		VALUES (:nexus_id, :title, :artist, :format, :spotify_id, :duration, :hotspot, :has_cover, :has_vocals, :has_combined);
	`, song)

	if err != nil {
		return fmt.Errorf("failed to create song: %w", err)
	}

	return nil
}

func (s Songs) GetSession(sessionId string) (*models.SongSession, error) {
	row := DB.QueryRowx(`
		SELECT id, event_id, nexus_id, session_nexus_id, title, artist, sung_by, added_at, started_at, ended_at, cancelled_at
		FROM song_session
		WHERE id = ?
	`, sessionId)

	if row.Err() != nil {
		return nil, row.Err()
	}

	session := models.SongSession{}
	err := row.StructScan(&session)

	return &session, err
}

func (s Songs) Delete(song models.Song) error {
	_, err := DB.Exec(`DELETE FROM song WHERE nexus_id = ?`, song.NexusId)
	if err != nil {
		return err
	}

	return nil
}

func (s Songs) WipeInvalidSessions() error {
	_, err := DB.Exec(`
		DELETE FROM song_session
		WHERE ended_at IS NULL
		  AND cancelled_at IS NULL
	`)
	if err != nil {
		return err
	}

	return nil
}

func (s Songs) CreateSession(session *models.SongSession) error {
	_, err := DB.Exec(`
		INSERT INTO song_session (event_id, nexus_id, session_nexus_id, title, artist, sung_by, added_at, started_at, ended_at, cancelled_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`,
		session.EventId,
		session.NexusId,
		session.SessionNexusId,
		session.Title,
		session.Artist,
		session.SungBy,
		session.AddedAt,
		session.StartedAt,
		session.EndedAt,
		session.CancelledAt,
	)

	if err != nil {
		log.Error("Failed to insert session", "err", err)
		return err
	}

	row := DB.QueryRow(`SELECT id FROM song_session WHERE rowid = last_insert_rowid();`)

	if row.Err() != nil {
		log.Error("Failed to fetch lastrowid from song_session", "err", row.Err())
		return row.Err()
	}

	var id int64
	err = row.Scan(&id)

	if err != nil {
		log.Error("Failed to scan lastrowid from song_session", "err", row.Err())
		return err
	}

	session.Id = id

	return nil
}

func (s Songs) UpdateSession(session *models.SongSession) error {
	_, err := DB.Exec(`
		UPDATE song_session
		SET
			session_nexus_id = ?,
			started_at = ?,
			ended_at = ?,
			cancelled_at = ?
		WHERE id = ?;
	`,
		session.SessionNexusId,
		session.StartedAt,
		session.EndedAt,
		session.CancelledAt,
		session.Id,
	)

	if err != nil {
		log.Error("Failed to insert session", "err", err)
		return err
	}

	return nil
}

func (s Songs) GetNotSyncedSessions(eventId int64) ([]models.SongSession, error) {
	rows, err := DB.Queryx(`
		SELECT id, event_id, nexus_id, session_nexus_id, title, artist, sung_by, added_at, started_at, ended_at, cancelled_at
		FROM song_session
		WHERE session_nexus_id IS NULL
		  AND ended_at IS NOT NULL
		  AND event_id = ?
	`, eventId)

	if err != nil {
		return nil, err
	}

	songSessions := []models.SongSession{}

	for rows.Next() {
		session := models.SongSession{}

		err = rows.StructScan(&session)
		if err != nil {
			return nil, err
		}

		songSessions = append(songSessions, session)
	}

	return songSessions, nil
}
