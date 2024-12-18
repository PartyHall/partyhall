package routes

import (
	"database/sql"
	"errors"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/partyhall/partyhall/api_errors"
	"github.com/partyhall/partyhall/dal"
	"github.com/partyhall/partyhall/log"
	"github.com/partyhall/partyhall/mercure_client"
	"github.com/partyhall/partyhall/nexus"
	"github.com/partyhall/partyhall/pipewire"
	routes_requests "github.com/partyhall/partyhall/routes/requests"
	"github.com/partyhall/partyhall/services"
	"github.com/partyhall/partyhall/state"
)

func routeSetMode(c *gin.Context) {
	mode := strings.ToLower(c.Params.ByName("mode"))
	if !slices.Contains(state.MODES, mode) {
		api_errors.ApiError(
			c,
			http.StatusBadRequest,
			"bad-request",
			"Selected mode not valid",
			"The selected mode is not available for this appliance",
		)

		return
	}

	if state.STATE.CurrentEvent == nil {
		state.STATE.CurrentMode = state.MODE_DISABLED

		mercure_client.CLIENT.PublishEvent("/mode", map[string]any{
			"mode": state.STATE.CurrentMode,
		})

		api_errors.ApiError(
			c,
			http.StatusBadRequest,
			"bad-request",
			"There is no event",
			"The appliance cannot be enabled when there is no event. Please create one first",
		)

		return
	}

	if state.STATE.CurrentMode != mode {
		state.STATE.CurrentMode = mode

		err := mercure_client.CLIENT.PublishEvent("/mode", map[string]any{
			"mode": mode,
		})
		if err != nil {
			c.Render(http.StatusInternalServerError, api_errors.MERCURE_PUBLISH_FAILURE)
			return
		}

	}

	c.Status(http.StatusNoContent)
}

func routeSetEvent(c *gin.Context) {
	eventStr := strings.ToLower(c.Params.ByName("event"))
	eventId, err := strconv.ParseInt(eventStr, 10, 64)
	if err != nil {
		api_errors.ApiError(
			c,
			http.StatusBadRequest,
			"bad-request",
			"Event ID not valid",
			"The event id could not be parsed.",
		)

		return
	}

	event, err := dal.EVENTS.Get(eventId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.Render(http.StatusInternalServerError, api_errors.DATABASE_ERROR.WithExtra(map[string]any{
				"err": err,
			}))
		}

		c.Render(http.StatusInternalServerError, api_errors.DATABASE_ERROR.WithExtra(map[string]any{
			"err": err,
		}))

		return
	}

	if state.STATE.CurrentEvent == nil || state.STATE.CurrentEvent.Id != event.Id {
		state.STATE.CurrentEvent = event

		err := mercure_client.CLIENT.PublishEvent("/event", event)
		if err != nil {
			c.Render(http.StatusInternalServerError, api_errors.MERCURE_PUBLISH_FAILURE)

			return
		}

		dal.EVENTS.Set(event)
	}

	c.JSON(http.StatusOK, event)
}

func routeSetDebug(c *gin.Context) {
	err := services.ShowDebug()
	if err != nil {
		c.Render(http.StatusInternalServerError, api_errors.MERCURE_PUBLISH_FAILURE)

		return
	}

	c.Status(http.StatusOK)
}

func routeGetAudioDevices(c *gin.Context) {
	devices, err := pipewire.GetDevices()
	if err != nil {
		api_errors.BAD_REQUEST.WithExtra(map[string]any{
			"err": err,
		})

		return
	}

	c.JSON(http.StatusOK, devices)
}

func routeSetAudioDevices(c *gin.Context) {
	var req routes_requests.AudioSetDevices
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

	mercure_client.CLIENT.PublishEvent(
		"/audio-devices",
		devices,
	)

	err = pipewire.LinkDevice(devices.DefaultSource, devices.DefaultSink)
	if err != nil {
		api_errors.BAD_REQUEST.WithExtra(map[string]any{
			"err": err,
		})

		return
	}

	c.JSON(200, devices)
}

func routeSetAudioDeviceVolume(c *gin.Context) {
	var req routes_requests.AudioSetDeviceVolume
	if err := c.ShouldBindJSON(&req); err != nil {
		api_errors.RenderValidationErr(c, err)

		return
	}

	deviceIdStr := strings.TrimSpace(c.Params.ByName("id"))
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

	var device *pipewire.Device = nil
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

	mercure_client.CLIENT.PublishEvent(
		"/audio-devices",
		devices,
	)

	c.JSON(200, devices)
}

func routeForceSync(c *gin.Context) {
	go func() {
		err := nexus.INSTANCE.Sync(state.STATE.CurrentEvent)
		if err != nil {
			log.Error("Failed to sync", "err", err)
		}
	}()

	c.Status(http.StatusOK)
}

func routeShutdown(c *gin.Context) {
	err := services.Shutdown()
	if err != nil {
		log.Error("Failed to shutdown", "err", err)
	}
}
