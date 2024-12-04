package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ROLE_GUEST     = "GUEST"
	ROLE_USER      = "USER"
	ROLE_ADMIN     = "ADMIN"
	ROLE_APPLIANCE = "APPLIANCE"
)

type User struct {
	Id       int    `db:"id" json:"id"`
	Name     string `db:"name" json:"name"`
	Username string `db:"username" json:"username"`
	Password string `db:"password" json:"-"`
	Roles    Roles  `db:"roles" json:"roles"`
}

type RefreshToken struct {
	Id        int       `db:"id"`
	User      *User     `db:"user"`
	Token     string    `db:"token"`
	ExpiresAt time.Time `db:"expires_at"`
}

type MercureClaims struct {
	Subscribe []string       `json:"subscribe"`
	Publish   []string       `json:"publish"`
	Payload   map[string]any `json:"payload"`
}

type JwtCustomClaims struct {
	Name     string        `json:"name"`
	Username string        `json:"username"`
	Roles    Roles         `json:"roles"`
	Mercure  MercureClaims `json:"mercure"`
	jwt.RegisteredClaims
}
