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

func (u *Users) Delete(id int) error {
	res, err := DB.Exec(`DELETE FROM ph_user WHERE id = ?`, id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (u *Users) GetCollection(amt, offset int) (*models.PaginatedResponse, error) {
	resp := models.PaginatedResponse{}

	row := DB.QueryRowx(`SELECT COUNT(*) FROM ph_user`)
	if row.Err() != nil {
		return nil, row.Err()
	}

	err := row.Scan(&resp.TotalCount)
	if err != nil {
		return nil, err
	}

	resp.CalculateMaxPage()

	rows, err := DB.Queryx(`
        SELECT id, username, name, password, roles
        FROM ph_user
        LIMIT ? OFFSET ?
    `, amt, offset)

	if err != nil {
		return nil, err
	}

	users := []models.User{}
	for rows.Next() {
		usr := models.User{}
		if err := rows.StructScan(&usr); err != nil {
			return nil, err
		}
		users = append(users, usr)
	}

	resp.Results = users

	return &resp, nil
}

func (u *Users) HasAnAdmin() (bool, error) {
	var existsInt int

	row := DB.QueryRowx(`
        SELECT EXISTS(
            SELECT 1 FROM ph_user, json_each(roles)
            WHERE json_each.value = ?
        )
    `, models.ROLE_ADMIN)

	if err := row.Scan(&existsInt); err != nil {
		return false, err
	}

	return existsInt == 1, nil
}
