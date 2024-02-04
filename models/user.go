package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ROLE_USER  = "USER"
	ROLE_ADMIN = "ADMIN"
)

type Roles []string

func (roles *Roles) Scan(val interface{}) error {
	switch v := val.(type) {
	case []byte:
		json.Unmarshal(v, &roles)
		return nil
	case string:
		json.Unmarshal([]byte(v), &roles)
		return nil
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
}

func (roles Roles) Value() (driver.Value, error) {
	return json.Marshal(roles)
}

type User struct {
	Id       int    `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	Username string `json:"username" db:"username"`
	Password string `json:"-" db:"password"`
	Roles    Roles  `json:"roles" db:"roles"`
}

type RefreshToken struct {
	Id        int        `db:"id"`
	User      *User      `db:"user"`
	Token     string     `db:"token"`
	ExpiresAt *Timestamp `db:"expires_at"`
}

type JwtCustomClaims struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Roles    Roles  `json:"roles"`
	jwt.RegisteredClaims
}
