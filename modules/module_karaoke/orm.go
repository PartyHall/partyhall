package module_karaoke

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/partyhall/partyhall/models"
	"github.com/partyhall/partyhall/orm"
)

func ormListSongs(query string, offset int, limit int) ([]PhkSong, error) {
	rows, err := orm.GET.DB.Queryx(`
		SELECT
			id,
			uuid,
			spotify_id,
			artist,
			title,
			hotspot,
			format,
			has_cover,
			has_vocals,
			has_full,
			filename
		FROM karaoke_song
		WHERE LENGTH($1) == 0
		OR (
			LOWER(filename) LIKE CONCAT('%', LOWER($1), '%')
			OR LOWER(artist) LIKE CONCAT('%', LOWER($1), '%')
			OR LOWER(title) LIKE CONCAT('%', LOWER($1), '%')
		)
		ORDER BY artist, title
		LIMIT $2
		OFFSET $3
	`, query, limit, offset)

	if err != nil {
		return nil, err
	}

	songs := []PhkSong{}

	for rows.Next() {
		var song PhkSong
		err = rows.StructScan(&song)
		if err != nil {
			return nil, err
		}

		songs = append(songs, song)
	}

	return songs, nil
}

func ormFindOneSongBy(condition string, expectedValue any) (*PhkSong, error) {
	row := orm.GET.DB.QueryRowx(`
		SELECT
			id,
			uuid,
			spotify_id,
			artist,
			title,
			hotspot,
			format,
			has_cover,
			has_vocals,
			has_full,
			filename
		FROM karaoke_song
		WHERE 
	`+condition, expectedValue)

	if row.Err() != nil {
		return nil, row.Err()
	}

	var song PhkSong
	err := row.StructScan(&song)
	if err != nil {
		return nil, err
	}

	return &song, nil
}

func ormLoadSongByFilename(filename string) (*PhkSong, error) {
	return ormFindOneSongBy("filename = $1", filename)
}

func ormLoadSongByUuid(uuid string) (*PhkSong, error) {
	return ormFindOneSongBy("uuid = $1", uuid)
}

func ormCreateSong(song PhkSong) (*PhkSong, error) {
	row := orm.GET.DB.QueryRowx(
		`
			INSERT INTO karaoke_song (uuid, spotify_id, artist, title, hotspot, format, has_cover, has_vocals, has_full, filename)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
			RETURNING id, uuid, spotify_id, artist, title, hotspot, format, has_cover, has_vocals, has_full, filename
		`,
		song.Uuid,
		song.SpotifyID,
		song.Artist,
		song.Title,
		song.Hotspot,
		song.Format,
		song.HasCover,
		song.HasVocals,
		song.HasFull,
		song.Filename,
	)

	if row.Err() != nil {
		return nil, row.Err()
	}

	err := row.StructScan(&song)

	if err != nil {
		return nil, err
	}

	return &song, nil
}

func ormDeleteSong(uuid string) {
	orm.GET.DB.Exec(`DELETE FROM karaoke_song WHERE uuid = $1`, uuid)
}

func ormFetchSongUUIDs() ([]string, error) {
	rows, err := orm.GET.DB.Query(`SELECT uuid FROM karaoke_song`)
	if err != nil {
		return nil, err
	}

	uuids := []string{}

	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			return nil, err
		}

		uuids = append(uuids, name)
	}

	return uuids, nil
}

func ormCountSongs(query string) (int, error) {
	row := orm.GET.DB.QueryRow(`
		SELECT COUNT(*)
		FROM karaoke_song
		WHERE LENGTH($1) == 0
		OR (
			LOWER(filename) LIKE CONCAT('%', LOWER($1), '%')
			OR LOWER(artist) LIKE CONCAT('%', LOWER($1), '%')
			OR LOWER(title) LIKE CONCAT('%', LOWER($1), '%')
		)
	`, query)
	if row.Err() != nil {
		return -1, row.Err()
	}

	var count int = -1
	err := row.Scan(&count)

	return count, err
}

/**
 * @TODO: This probably should not be an upsert but a create and a update separated
 * as eventId is not required when updating
 */
func ormUpsertSession(eventId int, session SongSession) (*SongSession, error) {
	var row *sqlx.Row

	if session.Id == 0 {
		currTime := time.Now().Unix()

		row = orm.GET.DB.QueryRowx(`
			INSERT INTO karaoke_song_session (
				song_id,
				event_id,
				sung_by,
				added_at
			) VALUES (
				$1,
				$2,
				$3,
				$4
			)
			RETURNING id, event_id, sung_by, added_at
		`, session.Song.Id, eventId, session.SungBy, currTime)
	} else {
		row = orm.GET.DB.QueryRowx(
			`
				UPDATE karaoke_song_session SET
					added_at = $1,
					started_at = $2,
					ended_at = $3,
					cancelled_at = $4
				WHERE id = $5
				RETURNING id, event_id, sung_by, added_at
			`,
			session.AddedAt,
			session.StartedAt,
			session.EndedAt,
			session.CancelledAt,
			session.Id,
		)
	}

	if row.Err() != nil {
		return nil, row.Err()
	}

	err := row.StructScan(&session)

	return &session, err
}

func ormSaveUnattendedPicture(session SongSession) (*SongImage, error) {
	row := orm.GET.DB.QueryRowx(`
		INSERT INTO karaoke_image(
			song_session_id,
			created_at
		) VALUES (
			$1,
			$2
		)
		RETURNING id, song_session_id, created_at
	`, session.Id, models.Timestamp(time.Now()))

	if row.Err() != nil {
		return nil, row.Err()
	}

	var img SongImage
	err := row.StructScan(&img)

	return &img, err
}
