package routes

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/partyhall/partyhall/api_errors"
	"github.com/partyhall/partyhall/config"
	"github.com/partyhall/partyhall/dal"
	"github.com/partyhall/partyhall/log"
	"github.com/partyhall/partyhall/mercure_client"
	"github.com/partyhall/partyhall/middlewares"
	"github.com/partyhall/partyhall/models"
	"github.com/partyhall/partyhall/nexus"
	"github.com/partyhall/partyhall/state"
)

type RoutesPhotobooth struct{}

func (h RoutesPhotobooth) Register(router *gin.RouterGroup) {
	// Has selected event & Authenticated
	router.POST(
		"take-picture",
		middlewares.HasEventLoaded(),
		middlewares.Authorized(),
		h.requestTakePicture,
	)

	// Has selected event & Appliance
	router.POST(
		"upload-picture",
		middlewares.HasEventLoaded(),
		middlewares.Authorized(models.ROLE_APPLIANCE),
		h.uploadPicture,
	)
}

func (h RoutesPhotobooth) requestTakePicture(c *gin.Context) {
	if state.STATE.CurrentEvent == nil {
		log.Error("Tried to take a picture from backend but no event selected")

		api_errors.ApiError(
			c,
			http.StatusBadRequest,
			"bad-request",
			"There is no event",
			"The appliance cannot take picture when there is no event. Please create one first",
		)
		return
	}

	mercure_client.CLIENT.SendTakePicture(false)
}

/**
 * @TODO: Validate input
 * Posting no image returns 200
 */
func (h RoutesPhotobooth) uploadPicture(c *gin.Context) {
	unattendedStr := c.PostForm("unattended")
	unattended, err := strconv.ParseBool(unattendedStr)
	if err != nil {
		c.Render(http.StatusBadRequest, api_errors.BAD_REQUEST.WithExtra(map[string]any{
			"err":  err.Error(),
			"line": "unattended parse bool",
		}))

		return
	}

	file, err := c.FormFile("picture")
	if err != nil {
		c.Render(http.StatusBadRequest, api_errors.BAD_REQUEST.WithExtra(map[string]any{
			"err":  err.Error(),
			"line": "get form picture file",
		}))

		return
	}

	/**
	 * Allow to keep a backup of the original image
	 * i.e. Using a backdrop, using special fx, etc..
	 **/
	alternateFile, err := c.FormFile("alternate_picture")
	hasAlternateFile := err != http.ErrMissingFile
	if err != nil && hasAlternateFile {
		c.Render(http.StatusBadRequest, api_errors.BAD_REQUEST.WithExtra(map[string]any{
			"err":  err.Error(),
			"line": "get form alternate picture file",
		}))

		return
	}

	basePath := filepath.Join(
		config.GET.EventPath,
		fmt.Sprintf("%v", state.STATE.CurrentEvent.Id), "photobooth",
	)

	if unattended {
		basePath = filepath.Join(basePath, "unattended")
	}

	err = os.MkdirAll(basePath, os.ModePerm)
	if err != nil {
		api_errors.ApiError(
			c,
			http.StatusInternalServerError,
			"generic-issue",
			"Failed to create folder",
			"This appliance is misconfigured",
		)

		return
	}

	imageName := fmt.Sprintf("%v.jpg", time.Now().Format("20060102-150405"))
	imagePath := filepath.Join(basePath, imageName)

	// Copy the image
	err = c.SaveUploadedFile(file, imagePath)
	if err != nil {
		api_errors.ApiErrorWithData(
			c,
			http.StatusInternalServerError,
			"generic-issue",
			"Failed to save image",
			"The image was not saved",
			map[string]any{
				"err": err,
			},
		)

		return
	}

	alternateImage := ""
	if hasAlternateFile {
		alternateImage = fmt.Sprintf("%v_alternate.jpg", time.Now().Format("20060102-150405"))
		alternateImagePath := filepath.Join(basePath, alternateImage)

		// Copy the image
		err = c.SaveUploadedFile(alternateFile, alternateImagePath)
		if err != nil {
			api_errors.ApiErrorWithData(
				c,
				http.StatusInternalServerError,
				"generic-issue",
				"Failed to save alternate image",
				"The alternate image was not saved",
				map[string]any{
					"err": err,
				},
			)

			return
		}
	}

	_, err = dal.EVENTS.InsertPicture(
		state.STATE.CurrentEvent.Id,
		imageName,
		unattended,
		alternateImage,
	)
	if err != nil {
		c.Render(http.StatusInternalServerError, api_errors.DATABASE_ERROR.WithExtra(map[string]any{
			"err":  err.Error(),
			"line": "insert picture in DB",
		}))

		return
	}

	go func() {
		event, err := dal.EVENTS.GetCurrent()
		if err != nil {
			log.Error("Failed to get new event after taking picture", "err", err)
		} else {
			mercure_client.CLIENT.SetCurrentEvent(event)
		}

		if err := nexus.INSTANCE.Sync(state.STATE.CurrentEvent); err != nil {
			log.Error("Failed to sync taken picture to Nexus", "err", err)
		}
	}()

	if !unattended {
		c.File(imagePath)

		return
	}

	c.Status(http.StatusCreated)
}
