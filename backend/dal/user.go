package dal

import (
	"database/sql"

	"github.com/partyhall/partyhall/models"
)

var USERS Users

type Users struct{}

func (u *Users) Create(dto models.User) (*models.User, error) {
	row := DB.QueryRowx(`
		INSERT INTO ph_user (username, name, password, roles)
		VALUES (?, ?, ?, ?)
		RETURNING id, username, name, password, roles
	`, dto.Username, dto.Name, dto.Password, dto.Roles)

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

func (u *Users) Update(dto models.User) (*models.User, error) {
	row := DB.QueryRowx(`
		UPDATE ph_user
		SET username = ?,
		    password = ?,
			name = ?,
			roles = ?
		WHERE id = ?
		RETURNING id, username, name, password, roles
	`, dto.Username, dto.Password, dto.Name, dto.Roles, dto.Id)

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
	row := DB.QueryRowx(`
		SELECT id, username, name, password, roles
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
	row := DB.QueryRowx(`
		SELECT id, username, name, password, roles
		FROM ph_user
		WHERE LOWER(username) = LOWER(?)
	`, username)

	if row.Err() != nil {
		if row.Err() == sql.ErrNoRows {
			return nil, nil
		}

		return nil, row.Err()
	}

	dbUser := models.User{}
	err := row.StructScan(&dbUser)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &dbUser, err
}

func (u *Users) FindByRefreshToken(token string) (*models.User, error) {
	row := DB.QueryRowx(`
		SELECT
			u.id as id,
			u.username as username,
			u.name as name,
			u.password as password,
			u.roles as roles
		FROM ph_user u
		INNER JOIN refresh_token rt ON rt.user_id = u.id
		WHERE rt.token = ? AND strftime('%s', rt.expires_at) > strftime('%s', 'now')
	`, token)

	if row.Err() != nil {
		return nil, row.Err()
	}

	dbUser := models.User{}
	err := row.StructScan(&dbUser)

	return &dbUser, err
}

func (u *Users) Upsert(user *models.User) (*models.User, error) {
	if user.Id == 0 {
		return u.Create(*user)
	}

	return u.Update(*user)
}

func (u *Users) DeleteRefreshToken(token string) error {
	_, err := DB.Exec("DELETE FROM refresh_token WHERE token = ?", token)
	return err
}

func (u *Users) CreateRefreshToken(userId int, token string) (int, error) {
	row := DB.QueryRow(`
		INSERT INTO refresh_token(token, expires_at, user_id)
		VALUES (?, datetime('now', '+7 days'), ?)
		RETURNING id
	`, token, userId)

	if row.Err() != nil {
		return 0, row.Err()
	}

	var rtId int
	err := row.Scan(&rtId)

	return rtId, err
}
