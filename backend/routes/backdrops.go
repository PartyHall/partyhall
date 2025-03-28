package routes

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/partyhall/partyhall/api_errors"
	"github.com/partyhall/partyhall/config"
	"github.com/partyhall/partyhall/dal"
)

func routeGetBackdropAlbums(c *gin.Context) {
	offset := 0

	search := c.Query("search")
	pageStr := c.Query("page")

	var page int = 1
	var err error
	if len(pageStr) > 0 {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			c.Render(http.StatusBadRequest, api_errors.INVALID_PARAMETERS.WithExtra(map[string]any{
				"page": "The page should be an integer",
			}))

			return
		}

		offset = (page - 1) * config.AMT_RESULTS_PER_PAGE
	}

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

func routeGetBackdropAlbum(c *gin.Context) {
	idStr := strings.TrimSpace(c.Params.ByName("albumId"))
	if len(idStr) == 0 {
		c.Render(http.StatusBadRequest, api_errors.BAD_REQUEST.WithExtra(map[string]any{
			"err": "Album ID should not be empty",
		}))

		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.Render(http.StatusBadRequest, api_errors.BAD_REQUEST.WithExtra(map[string]any{
			"err": "Album ID should be an integer",
		}))

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

func routeDownloadBackdrop(c *gin.Context) {
	albumIdStr := strings.TrimSpace(c.Params.ByName("albumId"))
	if len(albumIdStr) == 0 {
		c.Render(http.StatusBadRequest, api_errors.BAD_REQUEST.WithExtra(map[string]any{
			"err": "Album ID should not be empty",
		}))

		return
	}

	albumId, err := strconv.ParseInt(albumIdStr, 10, 64)
	if err != nil {
		c.Render(http.StatusBadRequest, api_errors.BAD_REQUEST.WithExtra(map[string]any{
			"err": "Album ID should be an integer",
		}))

		return
	}

	backdropIdStr := strings.TrimSpace(c.Params.ByName("backdropId"))
	if len(albumIdStr) == 0 {
		c.Render(http.StatusBadRequest, api_errors.BAD_REQUEST.WithExtra(map[string]any{
			"err": "Backdrop ID should not be empty",
		}))

		return
	}

	backdropId, err := strconv.ParseInt(backdropIdStr, 10, 64)
	if err != nil {
		c.Render(http.StatusBadRequest, api_errors.BAD_REQUEST.WithExtra(map[string]any{
			"err": "Backdrop ID should be an integer",
		}))

		return
	}

	backdrop, err := dal.BACKDROPS.Get(albumId, backdropId)
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
