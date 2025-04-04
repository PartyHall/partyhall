package routes

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/partyhall/partyhall/api_errors"
	"github.com/partyhall/partyhall/config"
	"github.com/partyhall/partyhall/mercure_client"
	"github.com/partyhall/partyhall/middlewares"
	"github.com/partyhall/partyhall/models"
	"github.com/partyhall/partyhall/mqtt"
	"github.com/partyhall/partyhall/nexus"
	"github.com/partyhall/partyhall/os_mgmt"
	"github.com/partyhall/partyhall/state"
)

var spotifydNameRegex = regexp.MustCompile(`^[\p{L}0-9]{1,64}$`) // Letters + numbers between 1 & 64 characters (NO SPACE ! See spotifyd/spotifyd#146)

type RoutesSettings struct{}

func (h RoutesSettings) Register(router *gin.RouterGroup) {
	// Non-onboarded | Admin
	router.PUT(
		"flash",
		middlewares.NotOnboardedOrRole("ADMIN"),
		h.setFlash,
	)

	// Non-onboarded | Admin
	// In post as it also sets the mode
	router.POST(
		"button-mappings",
		middlewares.NotOnboardedOrRole("ADMIN"),
		h.getButtonMappings,
	)

	// Non-onboarded | Admin
	router.PUT(
		"button-mappings",
		middlewares.NotOnboardedOrRole("ADMIN"),
		h.setButtonMappings,
	)

	// Non-onboarded | Admin
	router.GET(
		"button-mappings/actions",
		middlewares.NotOnboardedOrRole("ADMIN"),
		h.getButtonMappingsActions,
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

	router.GET(
		"ap",
		middlewares.NotOnboardedOrRole("ADMIN"),
		h.getAp,
	)

	router.PUT(
		"ap",
		middlewares.NotOnboardedOrRole("ADMIN"),
		h.setAp,
	)

	router.GET(
		"spotify",
		middlewares.NotOnboardedOrRole("ADMIN"),
		h.getSpotify,
	)

	router.PUT(
		"spotify",
		middlewares.NotOnboardedOrRole("ADMIN"),
		h.setSpotify,
	)

	// Non-onboarded | Admin
	router.GET(
		"nexus",
		middlewares.NotOnboardedOrRole("ADMIN"),
		h.getNexus,
	)

	router.PUT(
		"nexus",
		middlewares.NotOnboardedOrRole("ADMIN"),
		h.setNexus,
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

func (h RoutesSettings) getButtonMappings(c *gin.Context) {
	state.STATE.SetMode(state.MODE_BTN_SETUP)
	mercure_client.CLIENT.SetMode(state.STATE.CurrentMode)

	c.JSON(200, state.STATE.UserSettings.ButtonMappings)
}

func (h RoutesSettings) getButtonMappingsActions(c *gin.Context) {
	c.JSON(200, mqtt.BUTTON_ACTIONS)
}

func (h RoutesSettings) setButtonMappings(c *gin.Context) {
	var mappings map[int]string

	if err := c.ShouldBindJSON(&mappings); err != nil {
		api_errors.RenderValidationErr(c, err)

		return
	}

	config.GET.UserSettings.ButtonMappings = mappings
	config.GET.UserSettings.Save()

	state.STATE.UserSettings = config.GET.UserSettings
	mercure_client.CLIENT.SendUserSettings()

	c.JSON(200, mappings)
}

func (h RoutesSettings) getNexus(c *gin.Context) {
	us := config.GET.UserSettings

	c.JSON(200, map[string]any{
		"nexus_url":   us.NexusURL,
		"hardware_id": us.HardwareID,
		"bypass_ssl":  us.NexusIgnoreSSL,
	})
}

func (h RoutesSettings) setNexus(c *gin.Context) {
	var req struct {
		BaseUrl    string `json:"base_url"`
		HardwareId string `json:"hardware_id"`
		ApiKey     string `json:"api_key"`
		BypassSsl  bool   `json:"bypass_ssl"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		api_errors.RenderValidationErr(c, err)

		return
	}

	baseUrl := strings.TrimSuffix(req.BaseUrl, "/")

	response := map[string]any{
		"nexus_url":   baseUrl,
		"hardware_id": req.HardwareId,
		"bypass_ssl":  req.BypassSsl,
	}

	var errMessage = ""
	var successful = true

	if len(baseUrl) > 0 {
		errMessage, successful = nexus.ValidateCredentials(
			baseUrl,
			req.HardwareId,
			req.ApiKey,
			req.BypassSsl,
		)
	}

	response["error"] = nil
	if successful {
		config.GET.UserSettings.NexusURL = baseUrl
		config.GET.UserSettings.HardwareID = req.HardwareId
		config.GET.UserSettings.ApiKey = req.ApiKey
		config.GET.UserSettings.NexusIgnoreSSL = req.BypassSsl
		config.GET.UserSettings.Save()
		state.STATE.UserSettings = config.GET.UserSettings
	} else {
		response["error"] = errMessage
	}

	c.JSON(200, response)
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

func (h RoutesSettings) getAp(c *gin.Context) {
	c.JSON(200, config.GET.UserSettings.WirelessAp)
}

func (h RoutesSettings) setAp(c *gin.Context) {
	var req struct {
		Enabled  bool   `json:"enabled"`
		Ssid     string `json:"ssid" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		api_errors.RenderValidationErr(c, err)

		return
	}

	// @TODO:
	// Validate the hostname & password
	// Build the hostapd.conf

	if req.Enabled {
		// systemctl enable hostapd
		// systemctl restart hostapd
	} else {
		// systemctl disable --now hostapd
	}

	config.GET.UserSettings.WirelessAp.Enabled = req.Enabled
	config.GET.UserSettings.WirelessAp.Ssid = req.Ssid
	config.GET.UserSettings.WirelessAp.Password = req.Password
	config.GET.UserSettings.Save()

	state.STATE.UserSettings = config.GET.UserSettings

	// @TODO: Ensure that the config is working (No error while running the command, idk where to find the logs though)
	c.JSON(200, req)
}

func (h RoutesSettings) getSpotify(c *gin.Context) {
	c.JSON(200, config.GET.UserSettings.Spotify)
}

func (h RoutesSettings) setSpotify(c *gin.Context) {
	var req struct {
		Enabled bool   `json:"enabled"`
		Name    string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		api_errors.RenderValidationErr(c, err)

		return
	}

	if req.Enabled && !(*spotifydNameRegex).MatchString(req.Name) {
		c.Render(http.StatusBadRequest, api_errors.BAD_REQUEST.WithExtra(map[string]any{
			"err": "Invalid device name, should only contains letters and numbers (No space!)",
		}))

		return
	}

	// Building & saving the Spotifyd config + restarting the service
	err := os_mgmt.SetSpotifySettings(req.Enabled, req.Name)
	if err != nil {
		c.Render(http.StatusInternalServerError, api_errors.BAD_REQUEST.WithExtra(map[string]any{
			"err": err.Error(),
		}))

		return
	}

	// Saving the stuff in PH config (so that we don't have to parse the config)
	config.GET.UserSettings.Spotify.Enabled = req.Enabled
	config.GET.UserSettings.Spotify.Name = req.Name
	config.GET.UserSettings.Save()

	state.STATE.UserSettings = config.GET.UserSettings

	c.JSON(200, req)
}

// #region
func (h RoutesSettings) getAudioDevices(c *gin.Context) {
	devices, err := os_mgmt.GetDevices()
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

	err := os_mgmt.SetDefaultDevices(req.SourceId, req.SinkId)
	if err != nil {
		api_errors.BAD_REQUEST.WithExtra(map[string]any{
			"err": err,
		})

		return
	}

	// Note that the links will not be updated in the response
	// but we don't care, the font do not have to use them
	devices, err := os_mgmt.GetDevices()
	if err != nil {
		api_errors.BAD_REQUEST.WithExtra(map[string]any{
			"err": err,
		})

		return
	}

	mercure_client.CLIENT.SendAudioDevices(devices)

	err = os_mgmt.LinkDevice(devices.DefaultSource, devices.DefaultSink)
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

	devices, err := os_mgmt.GetDevices()
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

	err = os_mgmt.SetVolume(device, float64(req.Volume)/100)
	if err != nil {
		api_errors.BAD_REQUEST.WithExtra(map[string]any{
			"err": err,
		})

		return
	}

	devices, err = os_mgmt.GetDevices()
	if err != nil {
		api_errors.BAD_REQUEST.WithExtra(map[string]any{
			"err": err,
		})

		return
	}

	mercure_client.CLIENT.SendAudioDevices(devices)

	c.JSON(200, devices)
}

//#endregion
