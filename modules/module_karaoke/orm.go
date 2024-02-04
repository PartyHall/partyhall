package module_karaoke

import (
	"github.com/partyhall/partyhall/models"
	"github.com/partyhall/partyhall/orm"
)

func ormSearchSong(song string) ([]models.Song, error) {
	rows, err := orm.GET.DB.Queryx(`
		SELECT id, filename, COALESCE(artist, '') artist, COALESCE(title, '') title, format
		FROM song
		LIMIT 20
	`, song)

	if err != nil {
		return nil, err
	}

	songs := []models.Song{}

	for rows.Next() {
		var song models.Song
		err = rows.StructScan(&song)
		if err != nil {
			return nil, err
		}

		songs = append(songs, song)
	}

	return songs, nil
}

func ormListSongs(query string, offset int, limit int) ([]models.Song, error) {
	rows, err := orm.GET.DB.Queryx(`
		SELECT id, filename, COALESCE(artist, '') artist, COALESCE(title, '') title, format
		FROM song
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

	songs := []models.Song{}

	for rows.Next() {
		var song models.Song
		err = rows.StructScan(&song)
		if err != nil {
			return nil, err
		}

		songs = append(songs, song)
	}

	return songs, nil
}

func ormLoadSongByFilename(filename string) (*models.Song, error) {
	row := orm.GET.DB.QueryRowx(`
		SELECT id, filename, COALESCE(artist, '') artist, COALESCE(title, '') title, format
		FROM song
		WHERE filename = $1
	`, filename)

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

func ormCreateSong(filename, artist, title, format string) (*models.Song, error) {
	row := orm.GET.DB.QueryRowx(`
		INSERT INTO song (filename, artist, title, format)
		VALUES ($1, $2, $3, $4)
		RETURNING id, filename, artist, title, format
	`, filename, artist, title, format)

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

func ormDeleteSong(filename string) {
	orm.GET.DB.Exec(`DELETE FROM song WHERE filename = $1`, filename)
}

func ormFetchSongFilenames() ([]string, error) {
	rows, err := orm.GET.DB.Query(`SELECT filename FROM song`)
	if err != nil {
		return nil, err
	}

	names := []string{}

	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			return nil, err
		}

		names = append(names, name)
	}

	return names, nil
}

func ormCountSongPlayed(filename string) {
	orm.GET.DB.Exec(`UPDATE song SET play_count = play_count + 1 WHERE filename = $1`, filename)
}

func ormCountSongs() (int, error) {
	row := orm.GET.DB.QueryRow(`SELECT COUNT(*) FROM song`)
	if row.Err() != nil {
		return -1, row.Err()
	}

	var count int = -1
	err := row.Scan(&count)

	return count, err
}
