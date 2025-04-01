package dal

import (
	"github.com/partyhall/partyhall/models"
)

var BACKDROPS Backdrops

type Backdrops struct{}

func (b Backdrops) GetAllAlbums() ([]models.BackdropAlbum, error) {
	rows, err := DB.Queryx(`
		SELECT
			id,
			nexus_id,
			name,
			author,
			version
		FROM backdrop_album
		ORDER BY name ASC
	`)

	if err != nil {
		return nil, err
	}

	albums := []models.BackdropAlbum{}
	for rows.Next() {
		album := models.BackdropAlbum{}

		err := rows.StructScan(&album)
		if err != nil {
			return nil, err
		}

		albums = append(albums, album)
	}

	return albums, nil
}

func (b Backdrops) GetAlbumCollection(
	search string,
	amt,
	offset int,
) (*models.PaginatedResponse, error) {
	args := []any{search, search}

	resp := models.PaginatedResponse{}

	row := DB.QueryRow(`
		SELECT COUNT(DISTINCT ba.rowid)
		FROM backdrop_album ba
		WHERE (
			LENGTH(?) = 0
			OR ba.name LIKE '%' || ? || '%'
		)
	`, args...)
	if row.Err() != nil {
		return nil, row.Err()
	}

	err := row.Scan(&resp.TotalCount)
	if err != nil {
		return nil, err
	}

	resp.CalculateMaxPage()

	args = append(args, amt, offset)

	rows, err := DB.Queryx(`
		SELECT
		 ba.id,
		 ba.nexus_id,
		 ba.name,
		 ba.author,
		 ba.version
		FROM backdrop_album ba
		WHERE (
			LENGTH(?) = 0
			OR ba.name LIKE '%' || ? || '%'
		)
		ORDER BY ba.name ASC
		LIMIT ? OFFSET ?
	`, args...)

	if err != nil {
		return nil, err
	}

	albums := []models.BackdropAlbum{}

	for rows.Next() {
		album := models.BackdropAlbum{}

		err := rows.StructScan(&album)
		if err != nil {
			return nil, err
		}

		albums = append(albums, album)
	}

	resp.Results = albums

	return &resp, nil
}

func (b Backdrops) GetAlbum(albumId int64) (models.BackdropAlbum, error) {
	album := models.BackdropAlbum{}

	row := DB.QueryRowx(`
		SELECT
			id,
			nexus_id,
			name,
			author,
			version
		FROM backdrop_album
		WHERE id = ?
	`, albumId)

	if row.Err() != nil {
		return album, row.Err()
	}

	err := row.StructScan(&album)
	if err != nil {
		return album, err
	}

	backdrops, err := b.GetAll(albumId)
	if err != nil {
		return album, err
	}

	album.Backdrops = backdrops

	return album, nil
}

func (b Backdrops) CreateAlbum(backdropAlbum *models.BackdropAlbum) error {
	res, err := DB.NamedExec(`
		INSERT INTO backdrop_album (
			nexus_id,
			name,
			author,
			version
		)
		VALUES (:nexus_id, :name, :author, :version)
	`, backdropAlbum)

	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	backdropAlbum.Id = id

	return nil
}

func (b Backdrops) DeleteAlbum(id int64) ([]int64, error) {
	deletedBackdropIds := []int64{}

	backdrops, err := b.GetAll(id)
	if err != nil {
		return nil, err
	}

	for _, bd := range backdrops {
		if err := b.Delete(bd); err != nil {
			return deletedBackdropIds, err
		}

		deletedBackdropIds = append(deletedBackdropIds, bd.Id)
	}

	_, err = DB.Exec(`DELETE FROM backdrop_album WHERE id = ?`, id)
	if err != nil {
		return deletedBackdropIds, err
	}

	return deletedBackdropIds, nil
}

func (b Backdrops) GetAll(albumId int64) ([]models.Backdrop, error) {
	rows, err := DB.Queryx(`
		SELECT
			id,
			album_id,
			nexus_id,
			title,
			filename
		FROM backdrop
		WHERE album_id = ?
		ORDER BY id ASC
	`, albumId)

	if err != nil {
		return nil, err
	}

	backdrops := []models.Backdrop{}
	for rows.Next() {
		backdrop := models.Backdrop{}

		err := rows.StructScan(&backdrop)
		if err != nil {
			return nil, err
		}

		backdrops = append(backdrops, backdrop)
	}

	return backdrops, nil
}

func (ba Backdrops) GetCollection(
	albumId int64,
	search string,
	amt,
	offset int,
) (*models.PaginatedResponse, error) {
	args := []any{search, search, albumId}

	resp := models.PaginatedResponse{}

	row := DB.QueryRow(`
		SELECT COUNT(DISTINCT b.rowid)
		FROM backdrop b
		WHERE (
			LENGTH(?) = 0
			OR b.title LIKE '%' || ? || '%'
		)
		  AND b.album_id = ?
	`, args...)
	if row.Err() != nil {
		return nil, row.Err()
	}

	err := row.Scan(&resp.TotalCount)
	if err != nil {
		return nil, err
	}

	resp.CalculateMaxPage()

	args = append(args, amt, offset)

	rows, err := DB.Queryx(`
		SELECT
		 b.id,
		 b.album_id,
		 b.nexus_id,
		 b.title,
		 b.filename
		FROM backdrop b
		WHERE (
			LENGTH(?) = 0
			OR b.title LIKE '%' || ? || '%'
		)
		  AND b.album_id = ?
		ORDER BY b.name ASC
		LIMIT ? OFFSET ?
	`, args...)

	if err != nil {
		return nil, err
	}

	backdrops := []models.Backdrop{}

	for rows.Next() {
		backdrop := models.Backdrop{}

		err := rows.StructScan(&backdrop)
		if err != nil {
			return nil, err
		}

		backdrops = append(backdrops, backdrop)
	}

	resp.Results = backdrops

	return &resp, nil
}

func (b Backdrops) Get(backdropId int64) (models.Backdrop, error) {
	backdrop := models.Backdrop{}

	row := DB.QueryRowx(`
		SELECT
			id,
			album_id,
			nexus_id,
			title,
			filename
		FROM backdrop
		WHERE id = ?
	`, backdropId)

	if row.Err() != nil {
		return backdrop, row.Err()
	}

	err := row.StructScan(&backdrop)

	return backdrop, err
}

func (b Backdrops) Create(backdrop *models.Backdrop) error {
	res, err := DB.NamedExec(`
		INSERT INTO backdrop (
			album_id,
			nexus_id,
			title,
			filename
		)
		VALUES (:album_id, :nexus_id, :title, :filename)
	`, backdrop)

	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	backdrop.Id = id

	return nil
}

func (b Backdrops) Delete(bd models.Backdrop) error {
	_, err := DB.Exec(`DELETE FROM backdrop WHERE id = ?`, bd.Id)
	if err != nil {
		return err
	}

	return nil
}
