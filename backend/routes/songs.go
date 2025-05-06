package routes

import (
	"database/sql"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/partyhall/partyhall/api_errors"
	"github.com/partyhall/partyhall/config"
	"github.com/partyhall/partyhall/dal"
	"github.com/partyhall/partyhall/middlewares"
	"github.com/partyhall/partyhall/utils"
)

type RoutesSong struct{}

func (h RoutesSong) Register(router *gin.RouterGroup) {
	// Authenticated
	router.GET(
		"",
		middlewares.Authorized(),
		h.getCollection,
	)

	// Authenticated
	router.GET(
		":songId",
		middlewares.Authorized(),
		h.get,
	)

	router.GET(
		":songId/cover",
		h.getCover,
	)

	// Appliance
	router.GET(
		":songId/file/:songFilename",
		// middlewares.Authorized(models.ROLE_APPLIANCE), // too lazy to get the files manually to authenticate
		h.getFile,
	)
}

func (h RoutesSong) getCollection(c *gin.Context) {
	page, offset, err := utils.ParsePageOffset(c)
	if err != nil {
		return
	}

	search := c.Query("search")
	formatsStr := c.Query("formats")
	hasVocals, err := utils.ParseTrilean(c, "has_vocals")

	if err != nil {
		return
	}

	formats := []string{}
	if len(formatsStr) > 0 {
		formats = strings.Split(strings.ToLower(formatsStr), ",")
	}

	songs, err := dal.SONGS.GetCollection(
		search,
		formats,
		hasVocals,
		config.AMT_RESULTS_PER_PAGE,
		offset,
	)
	if err != nil {
		c.Render(http.StatusInternalServerError, api_errors.DATABASE_ERROR.WithExtra(map[string]any{
			"err": err.Error(),
		}))

		return
	}

	songs.Page = page

	c.JSON(http.StatusOK, songs)
}

func (h RoutesSong) get(c *gin.Context) {
	id := strings.TrimSpace(c.Params.ByName("songId"))
	if len(id) == 0 {
		c.Render(http.StatusBadRequest, api_errors.BAD_REQUEST.WithExtra(map[string]any{
			"err": "Song ID should not be empty",
		}))

		return
	}

	evt, err := dal.SONGS.Get(id)
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

	c.JSON(http.StatusOK, evt)
}

func (h RoutesSong) getCover(c *gin.Context) {
	songID := c.Param("songId")

	songPath := filepath.Join(config.GET.RootPath, "karaoke", songID, "cover.jpg")
	if _, err := os.Stat(songPath); os.IsNotExist(err) {
		c.Render(http.StatusNotFound, api_errors.BAD_REQUEST.WithExtra(map[string]any{
			"err": "Cover file does not exist",
		}))

		return
	}

	c.File(songPath)
}

func (h RoutesSong) getFile(c *gin.Context) {
	songID := c.Param("songId")
	file := c.Param("songFilename")

	filePath := filepath.Join(config.GET.RootPath, "karaoke", songID, file)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.Render(http.StatusNotFound, api_errors.BAD_REQUEST.WithExtra(map[string]any{
			"err": "The specified file does not exist",
		}))

		return
	}

	c.File(filePath)
}
