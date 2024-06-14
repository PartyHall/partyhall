package routes

import (
	"crypto/sha512"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/partyhall/partyhall/config"
	"github.com/partyhall/partyhall/dto"
	"github.com/partyhall/partyhall/logs"
	"github.com/partyhall/partyhall/middlewares"
	"github.com/partyhall/partyhall/models"
	"github.com/partyhall/partyhall/modules"
	"github.com/partyhall/partyhall/orm"
	"github.com/partyhall/partyhall/remote"
	"github.com/partyhall/partyhall/services"
	"github.com/partyhall/partyhall/utils"
)

func Register(g *echo.Group) {
	g.GET("/settings", settings)
	g.POST("/login", login)
	g.POST("/login-guest", loginGuest)
	g.POST("/refresh", refresh)

	g.GET("/socket/:type", remote.EasyWS.Route, services.GET.EchoWsJwtMiddleware)

	registerAdminRoutes(g.Group("/admin", services.GET.EchoJwtMiddleware, middlewares.RequireAdmin))
	modules.RegisterRoutes(g.Group("/modules"))
}

func settings(c echo.Context) error {
	return c.JSON(http.StatusOK, services.BuildFrontendSettings())
}

func login(c echo.Context) error {
	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var loginRequest LoginRequest
	if err := c.Bind(&loginRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON format"})
	}

	dbUser, err := orm.GET.Users.FindByUsername(loginRequest.Username)
	if err != nil {
		// @TODO Check if DB error
		return c.NoContent(http.StatusNotFound)
	}

	match, _ := services.GetArgon().VerifyPassword(loginRequest.Password, dbUser.Password)
	if !match {
		return c.NoContent(http.StatusNotFound)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, utils.GetClaimsFromUser(dbUser))
	tokenString, err := token.SignedString(services.GET.EchoJWTPrivateKey)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})
	}

	newRT, err := services.GetArgon().GenerateRandomBytes(128)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate refresh token"})
	}

	hasher := sha512.New()
	hasher.Write(newRT)
	newRefreshToken := fmt.Sprintf("%x", hasher.Sum(nil))

	_, err = orm.GET.Users.CreateRefreshToken(dbUser.Id, newRefreshToken)
	if err != nil {
		logs.Error(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save refresh token"})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token":         tokenString,
		"refresh_token": newRefreshToken,
	})
}

func loginGuest(c echo.Context) error {
	if !config.GET.GuestsAllowed {
		return c.NoContent(http.StatusForbidden)
	}

	type LoginRequest struct {
		Username string `json:"username"`
	}

	var loginRequest LoginRequest
	if err := c.Bind(&loginRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON format"})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, utils.GetClaimsFromUser(&models.User{
		Id:       0,
		Name:     loginRequest.Username,
		Username: loginRequest.Username,
		Roles:    []string{"GUEST"},
	}))
	tokenString, err := token.SignedString(services.GET.EchoJWTPrivateKey)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token":         tokenString,
		"refresh_token": "",
	})
}

func refresh(c echo.Context) error {
	var rt dto.UserRefresh
	if err := c.Bind(&rt); err != nil {
		return err
	}

	if len(strings.TrimSpace(rt.RefreshToken)) == 0 {
		return c.NoContent(http.StatusBadRequest)
	}

	dbUser, err := orm.GET.Users.FindByRefreshToken(rt.RefreshToken)
	err2 := orm.GET.Users.DeleteRefreshToken(rt.RefreshToken)
	newRT, err3 := services.GetArgon().GenerateRandomBytes(128)
	if err != nil || err2 != nil || err3 != nil {
		logs.Error(err, err2, err3)
		return c.NoContent(http.StatusUnauthorized)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, utils.GetClaimsFromUser(dbUser))
	tokenString, err := token.SignedString(services.GET.EchoJWTPrivateKey)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})
	}

	hasher := sha512.New()
	hasher.Write(newRT)
	newRefreshToken := fmt.Sprintf("%x", hasher.Sum(nil))

	_, err = orm.GET.Users.CreateRefreshToken(dbUser.Id, newRefreshToken)
	if err != nil {
		logs.Error(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save refresh token"})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token":         tokenString,
		"refresh_token": newRefreshToken,
	})
}
