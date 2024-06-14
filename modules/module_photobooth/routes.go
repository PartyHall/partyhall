package module_photobooth

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/partyhall/partyhall/logs"
	"github.com/partyhall/partyhall/orm"
	"github.com/partyhall/partyhall/remote"
	"github.com/partyhall/partyhall/services"
)

func takePictureRoute(c echo.Context) error {
	event := c.FormValue("event")
	unattended := c.FormValue("unattended")
	image, err := c.FormFile("image")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Unable to save picture: Getting picture => "+err.Error())
	}

	if len(event) == 0 || len(unattended) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to save picture: bad request")
	}

	isUnattended, err := strconv.ParseBool(unattended)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to parse unattended var: "+err.Error())
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

	if isUnattended {
		basePath = filepath.Join(basePath, "unattended")
	} else {
		basePath = filepath.Join(basePath, "pictures")
	}

	err = os.MkdirAll(basePath, os.ModePerm)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create path: "+err.Error())
	}

	/**
	 * Insert the image
	 **/
	imageName := fmt.Sprintf("%v.jpg", time.Now().Format("20060102-150405"))

	img, err := orm.GET.Events.InsertImage(services.GET.CurrentState.CurrentEvent, isUnattended)
	if err != nil {
		logs.Error("Failed to insert image: ", err)
		logs.Error("Defaulting name to current timestamp")
	} else {
		imageName = fmt.Sprintf("%v.jpg", img.Id)
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

	if !isUnattended {
		return c.File(basePath)
	}

	return c.NoContent(http.StatusCreated)

}
