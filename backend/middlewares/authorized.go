package middlewares

import (
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/partyhall/partyhall/api_errors"
	"github.com/partyhall/partyhall/config"
	"github.com/partyhall/partyhall/models"
)

func Authorized(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.Request.Header.Get("Authorization")
		if header == "" {
			c.Render(http.StatusUnauthorized, api_errors.NO_TOKEN)
			c.Abort()

			return
		}

		jwtToken := strings.Split(header, " ")
		if len(jwtToken) != 2 {
			c.Render(http.StatusUnauthorized, api_errors.NO_TOKEN)
			c.Abort()

			return
		}

		claims := &models.JwtCustomClaims{}

		token, err := jwt.ParseWithClaims(jwtToken[1], claims, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
			}

			return config.GET.Mercure.SubscriberKey, nil
		})

		if err != nil {
			c.Render(http.StatusForbidden, api_errors.INVALID_TOKEN.WithExtra(map[string]any{
				"err": err.Error(),
			}))

			c.Abort()

			return
		}

		if !token.Valid {
			c.Render(http.StatusForbidden, api_errors.INVALID_TOKEN.WithExtra(map[string]any{
				"err": "Token invalid",
			}))

			c.Abort()

			return
		}

		if len(roles) > 0 {
			hasRequiredRole := false
			for _, requiredRole := range roles {
				if slices.Contains(claims.Roles, requiredRole) {
					hasRequiredRole = true
					break
				}
			}

			if !hasRequiredRole {
				c.Status(http.StatusForbidden)
				c.Abort()
				return
			}
		}

		c.Set("TokenClaims", claims)

		c.Next()
	}
}
