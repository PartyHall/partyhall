package orm

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/partyhall/partyhall/models"
)

type Users struct {
	db *sqlx.DB
}

func (u *Users) Create(dto models.User) (*models.User, error) {
	row := u.db.QueryRowx(`
		INSERT INTO ph_user (username, password, roles)
		VALUES (?, ?, ?)
		RETURNING id, username, password, roles
	`, dto.Username, dto.Password, dto.Roles)

	if row.Err() != nil {
		return nil, row.Err()
	}

	var user models.User = models.User{}
	err := row.StructScan(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *Users) Get(id int) (*models.User, error) {
	row := u.db.QueryRowx(`
		SELECT id, username, password, roles
		FROM ph_user
		WHERE id = ?
	`, id)

	if row.Err() != nil {
		if row.Err() == sql.ErrNoRows {
			return nil, nil
		}

		return nil, row.Err()
	}

	dbUser := models.User{}
	err := row.StructScan(&dbUser)

	return &dbUser, err
}

func (u *Users) FindByUsername(username string) (*models.User, error) {
	row := u.db.QueryRowx(`
		SELECT id, username, password, roles
		FROM ph_user
		WHERE username = ?
	`, username)

	if row.Err() != nil {
		if row.Err() == sql.ErrNoRows {
			return nil, nil
		}

		return nil, row.Err()
	}

	dbUser := models.User{}
	err := row.StructScan(&dbUser)

	return &dbUser, err
}
