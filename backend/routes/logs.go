package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/partyhall/partyhall/api_errors"
	"github.com/partyhall/partyhall/log"
)

func routeGetLogs(c *gin.Context) {
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
