package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/partyhall/partyhall/config"
	"github.com/partyhall/partyhall/models"
)

func GetClaimsFromUser(dbUser *models.User, isGuest bool) models.JwtCustomClaims {
	var duration time.Duration = 15
	if config.GET.IsInDev {
		duration = 2
	}

	topics := []string{
		"/time",
		"/event",
		"/flash",
		"/mode",
		"/sync-progress",
		"/snackbar",
		"/backdrop-state",
		"/karaoke",
		"/karaoke-queue",
		"/karaoke-timecode",
		"/user-settings",
	}

	// Admin can subscribe to everything
	for _, r := range dbUser.Roles {
		if r == "ADMIN" {
			topics = append(topics, []string{
				"/audio-devices",
				"/logs",
			}...)
			break
		}
	}

	claims := jwt.RegisteredClaims{
		Issuer:   "PARTYHALL",
		Subject:  fmt.Sprintf("%v", dbUser.Id),
		IssuedAt: jwt.NewNumericDate(time.Now()),
	}

	if !isGuest {
		claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute * duration))
	}

	return models.JwtCustomClaims{
		Name:             dbUser.Name,
		Username:         dbUser.Username,
		Roles:            dbUser.Roles,
		RegisteredClaims: claims,
		Mercure: models.MercureClaims{
			Subscribe: topics,
			Publish:   []string{},
			Payload: map[string]any{
				"user_type": "user",
				"user_id":   dbUser.Id,
				"username":  dbUser.Username,
				"name":      dbUser.Name,
			},
		},
	}
}

func GetClaimsFromAppliance() models.JwtCustomClaims {
	return models.JwtCustomClaims{
		Name:     "PartyHall",
		Username: "PartyHall",
		Roles:    models.Roles{models.ROLE_USER, models.ROLE_ADMIN, models.ROLE_APPLIANCE},
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:   "PARTYHALL",
			Subject:  "appliance",
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
		Mercure: models.MercureClaims{
			Subscribe: []string{"*"},
			Publish:   []string{},
			Payload: map[string]any{
				"user_type": "appliance",
			},
		},
	}
}

func GetClaimsPublisher() models.JwtCustomClaims {
	return models.JwtCustomClaims{
		Mercure: models.MercureClaims{
			Subscribe: []string{},
			Publish:   []string{"*"},
			Payload: map[string]any{
				"user_type": "publisher",
			},
		},
	}
}
