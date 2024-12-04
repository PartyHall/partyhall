package routes

import (
	"crypto/sha512"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/partyhall/partyhall/api_errors"
	"github.com/partyhall/partyhall/config"
	"github.com/partyhall/partyhall/dal"
	"github.com/partyhall/partyhall/log"
	"github.com/partyhall/partyhall/models"
	routes_requests "github.com/partyhall/partyhall/routes/requests"
	"github.com/partyhall/partyhall/services"
	"github.com/partyhall/partyhall/utils"
)

func routeLogin(c *gin.Context) {
	var loginRequest routes_requests.LoginRequest
	if err := c.Bind(&loginRequest); err != nil {
		api_errors.ApiError(
			c,
			http.StatusBadRequest,
			"invalid-json",
			"Invalid JSON provided",
			"The JSON does not match the expected content or is invalid",
		)

		return
	}

	dbUser, err := dal.USERS.FindByUsername(loginRequest.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			api_errors.ApiError(
				c,
				http.StatusBadRequest,
				"bad-login",
				"Bad login",
				"No user matching the username/password combo has been found.",
			)

			return
		}

		api_errors.ApiError(
			c,
			http.StatusInternalServerError,
			"",
			"Failed to query the database",
			err.Error(),
		)

		return
	}

	match, _ := services.GetArgon().VerifyPassword(loginRequest.Password, dbUser.Password)
	if !match {
		api_errors.ApiError(
			c,
			http.StatusBadRequest,
			"bad-login",
			"Bad login",
			"No user matching the username/password combo has been found.",
		)

		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, utils.GetClaimsFromUser(dbUser, false))
	tokenString, err := token.SignedString(config.GET.Mercure.SubscriberKey)
	if err != nil {
		api_errors.ApiError(
			c,
			http.StatusInternalServerError,
			"auth-issue",
			"Failed to generate token",
			err.Error(),
		)

		return
	}

	newRT, err := services.GetArgon().GenerateRandomBytes(128)
	if err != nil {
		api_errors.ApiError(
			c,
			http.StatusInternalServerError,
			"auth-issue",
			"Failed to generate refresh token",
			err.Error(),
		)

		return
	}

	hasher := sha512.New()
	hasher.Write(newRT)
	newRefreshToken := fmt.Sprintf("%x", hasher.Sum(nil))

	_, err = dal.USERS.CreateRefreshToken(dbUser.Id, newRefreshToken)
	if err != nil {
		log.Error("Auth failure", "err", err)
		api_errors.ApiError(
			c,
			http.StatusInternalServerError,
			"db-issue",
			"Failed to save refresh token",
			err.Error(),
		)

		return
	}

	c.JSON(http.StatusOK, map[string]string{
		"token":         tokenString,
		"refresh_token": newRefreshToken,
	})
}

func routeLoginGuest(c *gin.Context) {
	if !config.GET.GuestsAllowed {
		c.Status(http.StatusForbidden)
		return
	}

	var loginRequest routes_requests.LoginGuestRequest
	if err := c.Bind(&loginRequest); err != nil {
		api_errors.ApiError(
			c,
			http.StatusBadRequest,
			"invalid-json",
			"Invalid JSON provided",
			"The JSON does not match the expected content or is invalid",
		)

		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, utils.GetClaimsFromUser(&models.User{
		Id:       0,
		Name:     loginRequest.Username,
		Username: loginRequest.Username,
		Roles:    []string{"GUEST"},
	}, true))

	tokenString, err := token.SignedString(config.GET.Mercure.SubscriberKey)
	if err != nil {
		api_errors.ApiError(
			c,
			http.StatusInternalServerError,
			"auth-issue",
			"Failed to generate token",
			err.Error(),
		)

		return
	}

	c.JSON(http.StatusOK, map[string]string{
		"token":         tokenString,
		"refresh_token": "no-token",
	})
}

func routeLoginRefresh(c *gin.Context) {
	var rt routes_requests.RefreshRequest
	if err := c.Bind(&rt); err != nil {
		api_errors.ApiErrorWithData(
			c,
			http.StatusBadRequest,
			"invalid-json",
			"Invalid JSON provided",
			"The JSON does not match the expected content or is invalid",
			map[string]any{
				"err": err,
			},
		)

		return
	}

	if len(strings.TrimSpace(rt.RefreshToken)) == 0 {
		api_errors.ApiError(
			c,
			http.StatusBadRequest,
			"invalid-refresh",
			"Invalid refresh token",
			"The given refresh token is invalid",
		)

		return
	}

	dbUser, err := dal.USERS.FindByRefreshToken(rt.RefreshToken)
	err2 := dal.USERS.DeleteRefreshToken(rt.RefreshToken)
	newRT, err3 := services.GetArgon().GenerateRandomBytes(128)
	if err != nil || err2 != nil || (err3 != nil && !errors.Is(err3, sql.ErrNoRows)) {
		log.Error("Failed to refresh", "err", err, "err2", err2, "err3", err3)

		c.Status(http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, utils.GetClaimsFromUser(dbUser, true))
	tokenString, err := token.SignedString(config.GET.Mercure.SubscriberKey)
	if err != nil {
		api_errors.ApiError(
			c,
			http.StatusInternalServerError,
			"auth-issue",
			"Failed to generate token",
			err.Error(),
		)

		return
	}

	hasher := sha512.New()
	hasher.Write(newRT)
	newRefreshToken := fmt.Sprintf("%x", hasher.Sum(nil))

	_, err = dal.USERS.CreateRefreshToken(dbUser.Id, newRefreshToken)
	if err != nil {
		log.Error("Failed to refresh", "err", err)
		api_errors.ApiError(
			c,
			http.StatusInternalServerError,
			"db-issue",
			"Failed to save refresh token",
			err.Error(),
		)

		return
	}

	c.JSON(http.StatusOK, map[string]string{
		"token":         tokenString,
		"refresh_token": newRefreshToken,
	})
}
