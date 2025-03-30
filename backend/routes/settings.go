package routes

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/partyhall/partyhall/api_errors"
	"github.com/partyhall/partyhall/config"
	"github.com/partyhall/partyhall/mercure_client"
	"github.com/partyhall/partyhall/middlewares"
	"github.com/partyhall/partyhall/models"
	"github.com/partyhall/partyhall/mqtt"
	"github.com/partyhall/partyhall/pipewire"
	"github.com/partyhall/partyhall/state"
)

type RoutesSettings struct{}

func (h RoutesSettings) Register(router *gin.RouterGroup) {
	// Non-onboarded | Admin
	router.PUT(
		"flash",
		middlewares.NotOnboardedOrRole("ADMIN"),
		h.setFlash,
	)

	// Non-onboarded | Admin
	router.PUT(
		"webcam",
		middlewares.NotOnboardedOrRole("ADMIN"),
		h.setWebcam,
	)

	// Non-onboarded | Admin
	router.PUT(
		"unattended",
		middlewares.NotOnboardedOrRole("ADMIN"),
		h.setUnattended,
	)

	//region Maybe to rework
	// Non-onboarded | Admin
	router.GET(
		"audio-devices",
		middlewares.NotOnboardedOrRole("ADMIN"),
		h.getAudioDevices,
	)
	// Non-onboarded | Admin
	router.POST(
		"audio-devices",
		middlewares.NotOnboardedOrRole("ADMIN"),
		h.setAudioDevices,
	)

	// Non-onboarded | Admin
	router.PUT(
		"audio-devices/:deviceId",
		middlewares.NotOnboardedOrRole("ADMIN"),
		h.setAudioDeviceVolume,
	)
	//endregion
}

func (h RoutesSettings) setWebcam(c *gin.Context) {
	var req struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		api_errors.RenderValidationErr(c, err)

		return
	}

	config.GET.UserSettings.Photobooth.Resolution.Width = req.Width
	config.GET.UserSettings.Photobooth.Resolution.Height = req.Height
	config.GET.UserSettings.Save()

	state.STATE.UserSettings = config.GET.UserSettings
	mercure_client.CLIENT.SendUserSettings()

	c.JSON(200, config.GET.UserSettings.Photobooth.Resolution)
}

func (h RoutesSettings) setFlash(c *gin.Context) {
	var req struct {
		Powered    bool `json:"powered"`
		Brightness int  `json:"brightness"`
	}

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

	if req.Brightness != config.GET.UserSettings.Photobooth.FlashBrightness {
		config.GET.UserSettings.Photobooth.FlashBrightness = req.Brightness
		config.GET.UserSettings.Save()

		state.STATE.UserSettings = config.GET.UserSettings
		mercure_client.CLIENT.SendUserSettings()
	}

	mqtt.SetFlash(req.Powered, req.Brightness)

	c.Status(200)
}

func (h RoutesSettings) setUnattended(c *gin.Context) {
	var req struct {
		Enabled  bool `json:"enabled"`
		Interval int  `json:"interval"`
	}

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

	dirty := false
	if req.Enabled != config.GET.UserSettings.Photobooth.Unattended.Enabled {
		dirty = true
		config.GET.UserSettings.Photobooth.Unattended.Enabled = req.Enabled
	}

	if req.Interval != config.GET.UserSettings.Photobooth.Unattended.Interval {
		dirty = true
		config.GET.UserSettings.Photobooth.Unattended.Interval = req.Interval
	}

	if dirty {
		config.GET.UserSettings.Save()

		state.STATE.UserSettings = config.GET.UserSettings
		mercure_client.CLIENT.SendUserSettings()
	}

	c.Status(200)
}

func (h RoutesSettings) getAudioDevices(c *gin.Context) {
	devices, err := pipewire.GetDevices()
	if err != nil {
		c.Render(http.StatusBadRequest, api_errors.BAD_REQUEST.WithExtra(map[string]any{
			"err": err.Error(),
		}))

		return
	}

	c.JSON(http.StatusOK, devices)
}

func (h RoutesSettings) setAudioDevices(c *gin.Context) {
	var req struct {
		SourceId int `json:"source_id" binding:"required"`
		SinkId   int `json:"sink_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		api_errors.RenderValidationErr(c, err)

		return
	}

	err := pipewire.SetDefaultDevices(req.SourceId, req.SinkId)
	if err != nil {
		api_errors.BAD_REQUEST.WithExtra(map[string]any{
			"err": err,
		})

		return
	}

	// Note that the links will not be updated in the response
	// but we don't care, the font do not have to use them
	devices, err := pipewire.GetDevices()
	if err != nil {
		api_errors.BAD_REQUEST.WithExtra(map[string]any{
			"err": err,
		})

		return
	}

	mercure_client.CLIENT.SendAudioDevices(devices)

	err = pipewire.LinkDevice(devices.DefaultSource, devices.DefaultSink)
	if err != nil {
		api_errors.BAD_REQUEST.WithExtra(map[string]any{
			"err": err,
		})

		return
	}

	c.JSON(200, devices)
}

func (h RoutesSettings) setAudioDeviceVolume(c *gin.Context) {
	var req struct {
		Volume int `json:"volume"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		api_errors.RenderValidationErr(c, err)

		return
	}

	deviceIdStr := strings.TrimSpace(c.Params.ByName("deviceId"))
	deviceId, err := strconv.Atoi(deviceIdStr)

	if len(deviceIdStr) == 0 || err != nil {
		api_errors.INVALID_PARAMETERS.WithExtra(map[string]any{
			"err": err,
		})

		return
	}

	devices, err := pipewire.GetDevices()
	if err != nil {
		api_errors.BAD_REQUEST.WithExtra(map[string]any{
			"err": err,
		})

		return
	}

	var device *models.PwDevice = nil
	if devices.DefaultSink != nil && deviceId == devices.DefaultSink.ID {
		device = devices.DefaultSink
	} else if deviceId == devices.KaraokeSink.ID {
		device = &devices.KaraokeSink
	} else {
		api_errors.BAD_REQUEST.WithExtra(map[string]any{
			"err": "The requested device is not set as default",
		})

		return
	}

	err = pipewire.SetVolume(device, float64(req.Volume)/100)
	if err != nil {
		api_errors.BAD_REQUEST.WithExtra(map[string]any{
			"err": err,
		})

		return
	}

	devices, err = pipewire.GetDevices()
	if err != nil {
		api_errors.BAD_REQUEST.WithExtra(map[string]any{
			"err": err,
		})

		return
	}

	mercure_client.CLIENT.SendAudioDevices(devices)

	c.JSON(200, devices)
}
