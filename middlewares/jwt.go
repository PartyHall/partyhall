package middlewares

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/partyhall/partyhall/models"
)

func RequireAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*models.JwtCustomClaims)

		for _, role := range claims.Roles {
			if role == "ADMIN" {
				return next(c)
			}
		}

		return echo.ErrUnauthorized
	}
}
