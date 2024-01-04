package module_karaoke

import (
	"github.com/partyhall/partyhall/models"
	"github.com/partyhall/partyhall/orm"
)

func ormSearchSong(song string) ([]models.Song, error) {
	rows, err := orm.GET.DB.Queryx(`
		SELECT id, filename, COALESCE(artist, '') artist, COALESCE(title, '') title, format
		FROM song
		WHERE LOWER(filename) LIKE CONCAT('%', LOWER($1), '%')
		   OR LOWER(artist) LIKE CONCAT('%', LOWER($1), '%')
		   OR LOWER(title) LIKE CONCAT('%', LOWER($1), '%')
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

func ormListSongs(offset int64, limit int64) ([]models.Song, error) {
	rows, err := orm.GET.DB.Queryx(`
		SELECT id, filename, COALESCE(artist, '') artist, COALESCE(title, '') title, format
		FROM song
		ORDER BY artist, title
		LIMIT $1
		OFFSET $2
	`, limit, offset)

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

func ormCreateSong(filename, artist, title, format string) error {
	_, err := orm.GET.DB.Exec(`
		INSERT INTO song (filename, artist, title, format)
		VALUES ($1, $2, $3, $4)
	`, filename, artist, title, format)

	return err
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
