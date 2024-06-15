package module_karaoke

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/partyhall/partyhall/logs"
	"github.com/partyhall/partyhall/remote"
)

// @TODO: Factorize this with the photobooth module one maybe ?
// At least keep it in sync
func takePictureRoute(c echo.Context) error {
	if INSTANCE.CurrentSong == nil {
		errMsg := "Tried to take a karaoke unattended picture while not singing a song"
		logs.Error(errMsg)
		return echo.NewHTTPError(http.StatusBadRequest, errMsg)
	}

	image, err := c.FormFile("image")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Unable to save picture: Getting picture => "+err.Error())
	}

	src, err := image.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to open image: "+err.Error())
	}
	defer src.Close()

	basePath, err := getModuleEventDir()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get event module dir: "+err.Error())
	}

	basePath = filepath.Join(basePath, fmt.Sprintf("%v", INSTANCE.CurrentSong.Id))

	err = os.MkdirAll(basePath, os.ModePerm)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create path: "+err.Error())
	}

	/**
	 * Insert the image
	 **/
	imageName := fmt.Sprintf("%v.jpg", time.Now().Format("20060102-150405"))

	_, err = ormSaveUnattendedPicture(*INSTANCE.CurrentSong)
	if err != nil {
		logs.Error("Failed to insert image: ", err)
		logs.Error("Defaulting name to current timestamp")

		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save the image: "+err.Error())
	}

	basePath = filepath.Join(basePath, imageName)

	f, err := os.Create(basePath)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create image file: "+err.Error())
	}

	if _, err := io.Copy(f, src); err != nil {
		f.Close()
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save image file: "+err.Error())
	}

	if err = f.Sync(); err != nil {
		logs.Error("Failed to sync the data ! be careful")
	}

	f.Close()

	// Broadcasting the state so that the current event is refreshed on the admin panel
	remote.BroadcastState()

	return c.NoContent(http.StatusCreated)

}
