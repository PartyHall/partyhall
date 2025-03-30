package routes

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/partyhall/partyhall/api_errors"
	"github.com/partyhall/partyhall/config"
	"github.com/partyhall/partyhall/dal"
	"github.com/partyhall/partyhall/middlewares"
	"github.com/partyhall/partyhall/utils"
)

type RoutesBackdropAlbums struct{}

func (h RoutesBackdropAlbums) Register(router *gin.RouterGroup) {
	// Onboarded & Authenticated
	router.GET(
		"",
		middlewares.Onboarded(true),
		middlewares.Authorized(),
		h.getCollection,
	)

	// Onboarded & Authenticated
	router.GET(
		":backdropAlbumId",
		middlewares.Onboarded(true),
		middlewares.Authorized(),
		h.getCollection,
	)
}

func (h RoutesBackdropAlbums) getCollection(c *gin.Context) {
	page, offset, err := utils.ParsePageOffset(c)
	if err != nil {
		return
	}

	search := c.Query("search")

	albums, err := dal.BACKDROPS.GetAlbumCollection(
		search,
		config.AMT_RESULTS_PER_PAGE,
		offset,
	)
	if err != nil {
		c.Render(http.StatusInternalServerError, api_errors.DATABASE_ERROR.WithExtra(map[string]any{
			"err": err.Error(),
		}))

		return
	}

	albums.Page = page

	c.JSON(http.StatusOK, albums)
}

func (h RoutesBackdropAlbums) get(c *gin.Context) {
	id, parseFailed := utils.ParamAsIntOrError(c, "backdropAlbumId")
	if parseFailed {
		return
	}

	album, err := dal.BACKDROPS.GetAlbum(id)
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

	c.JSON(http.StatusOK, album)
}
