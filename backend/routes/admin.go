package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/partyhall/partyhall/api_errors"
	"github.com/partyhall/partyhall/log"
	"github.com/partyhall/partyhall/middlewares"
	"github.com/partyhall/partyhall/services"
)

type RoutesAdmin struct{}

func (h RoutesAdmin) Register(router *gin.RouterGroup) {
	// Onboarded & Admin
	router.GET(
		"logs",
		middlewares.Onboarded(true),
		middlewares.Authorized("ADMIN"),
		h.getLogs,
	)

	// Onboarded & Admin
	router.POST(
		"shutdown",
		middlewares.Onboarded(true),
		middlewares.Authorized("ADMIN"),
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
