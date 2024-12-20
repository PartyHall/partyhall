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
	"github.com/partyhall/partyhall/mqtt"
	"github.com/partyhall/partyhall/nexus"
	routes_requests "github.com/partyhall/partyhall/routes/requests"
	"github.com/partyhall/partyhall/state"
)

func routeTakePicture(c *gin.Context) {
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

	mercure_client.CLIENT.PublishEvent("/take-picture", map[string]any{
		"unattended": false,
	})
}

/**
 * @TODO: Validate input
 * Posting no image returns 200
 */
func routeUploadPicture(c *gin.Context) {
	if state.STATE.CurrentEvent == nil {
		log.Error("Tried to upload a picture but no event selected")

		api_errors.ApiError(
			c,
			http.StatusBadRequest,
			"bad-request",
			"There is no event",
			"The appliance cannot take picture when there is no event. Please create one first",
		)
		return
	}

	unattendedStr := c.PostForm("unattended")
	unattended, err := strconv.ParseBool(unattendedStr)
	if err != nil {
		c.Render(http.StatusBadRequest, api_errors.BAD_REQUEST.WithExtra(map[string]any{
			"err": err.Error(),
		}))

		return
	}

	file, err := c.FormFile("picture")
	if err != nil {
		c.Render(http.StatusBadRequest, api_errors.BAD_REQUEST.WithExtra(map[string]any{
			"err": err.Error(),
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

	_, err = dal.EVENTS.InsertPicture(
		state.STATE.CurrentEvent.Id,
		imageName,
		unattended,
	)
	if err != nil {
		c.Render(http.StatusInternalServerError, api_errors.DATABASE_ERROR.WithExtra(map[string]any{
			"err": err.Error(),
		}))

		return
	}

	go func() {
		event, err := dal.EVENTS.GetCurrent()
		if err != nil {
			log.Error("Failed to get new event after taking picture", "err", err)
		} else {
			mercure_client.CLIENT.PublishEvent("/event", event)
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

func routeSetFlash(c *gin.Context) {
	var req routes_requests.SetFlashRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api_errors.RenderValidationErr(c, err)

		return
	}

	if mqtt.EasyMqtt == nil {
		api_errors.BAD_REQUEST.WithExtra(map[string]any{
			"err": "No MQTT server for some reason",
		}).Render(c.Writer)

		return
	}

	mqtt.SetFlash(req.Powered, req.Brightness)

	c.Status(200)
}
