package routes

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/partyhall/partyhall/api_errors"
	"github.com/partyhall/partyhall/config"
	"github.com/partyhall/partyhall/dal"
	"github.com/partyhall/partyhall/middlewares"
	"github.com/partyhall/partyhall/utils"
)

type RoutesBackdrops struct{}

func (h RoutesBackdrops) Register(router *gin.RouterGroup) {
	// Onboarded only (Too lazy to implement the authenticated images)
	router.GET(
		":backdropId/download",
		middlewares.Onboarded(true),
		h.download,
	)
}

func (h RoutesBackdrops) download(c *gin.Context) {
	backdropId, parseFailed := utils.ParamAsIntOrError(c, "backdropId")
	if parseFailed {
		return
	}

	backdrop, err := dal.BACKDROPS.Get(backdropId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.Status(http.StatusNotFound)
			return
		}

		c.Render(http.StatusInternalServerError, api_errors.DATABASE_ERROR.WithExtra(map[string]any{
			"err": err,
		}))

		return
	}

	c.File(filepath.Join(
		config.GET.RootPath,
		"backdrops",
		fmt.Sprintf("%v", backdrop.AlbumId),
		backdrop.Filename,
	))
}
