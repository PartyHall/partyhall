package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/partyhall/partyhall/api_errors"
	"github.com/partyhall/partyhall/dal"
	"github.com/partyhall/partyhall/log"
	"github.com/partyhall/partyhall/middlewares"
	"github.com/partyhall/partyhall/models"
	"github.com/partyhall/partyhall/services"
	"github.com/partyhall/partyhall/state"
)

type RoutesAdmin struct{}

func (h RoutesAdmin) Register(router *gin.RouterGroup) {
	router.POST(
		"/create-admin",
		h.createAdmin,
	)

	router.GET(
		"logs",
		middlewares.Authorized(models.ROLE_ADMIN),
		h.getLogs,
	)

	router.POST(
		"shutdown",
		middlewares.Authorized(models.ROLE_ADMIN),
		h.shutdown,
	)
}

/**
 * This will be an infinite scroll at some point
 * This means that offset should be based on ID
 * Not arbitrary values
 **/
func (h RoutesAdmin) getLogs(c *gin.Context) {
	count, err := log.CountMessages()
	if err != nil {
		c.Render(http.StatusInternalServerError, api_errors.DATABASE_ERROR.WithExtra(map[string]any{
			"err": err,
		}))

		return
	}

	logs, err := log.GetMessages(100, 0)
	if err != nil {
		c.Render(http.StatusInternalServerError, api_errors.DATABASE_ERROR.WithExtra(map[string]any{
			"err": err,
		}))

		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"results":        logs,
		"total_count":    count,
		"per_page_count": 100,
	})
}

func (h RoutesAdmin) shutdown(c *gin.Context) {
	err := services.Shutdown()
	if err != nil {
		log.Error("Failed to shutdown", "err", err)
	}
}

func (h RoutesAdmin) createAdmin(c *gin.Context) {
	// Ultra-important check
	// Under NO CIRCUMSTANCES should we allow the creation of an admin if one already exists
	// as this endpoint is not protected by any authentication
	// It should ONLY be used during the onboarding process
	if state.STATE.AdminCreated {
		c.Status(http.StatusConflict)

		return
	}

	var req struct {
		Username string `json:"username" binding:"required"`
		Name     string `json:"name"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		api_errors.RenderValidationErr(c, err)
		return
	}

	hashed, err := services.GetArgon().Hash(req.Password)
	if err != nil {
		c.Render(http.StatusInternalServerError, api_errors.DATABASE_ERROR.WithExtra(map[string]any{"err": err.Error()}))
		return
	}

	userDTO := models.User{
		Username: req.Username,
		Name:     req.Name,
		Password: hashed,
		Roles: models.Roles([]string{
			models.ROLE_USER,
			models.ROLE_ADMIN,
		}),
	}

	created, err := dal.USERS.Create(userDTO)
	if err != nil {
		c.Render(http.StatusInternalServerError, api_errors.DATABASE_ERROR.WithExtra(map[string]any{"err": err.Error()}))
		return
	}

	// No need to send through mercure
	// as we expect only one connection to the admin interface
	// when not onboarded

	// EDIT: Or maybe we should, so that we can put a splashscreen
	// on the appliance until an admin is created
	// such as "Your appliance is nearly ready, please go to the admin interface to finish the setup"
	state.STATE.AdminCreated = true

	c.JSON(http.StatusOK, created)
}
