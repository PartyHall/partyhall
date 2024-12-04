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

func routeForceSync(c *gin.Context) {
	go func() {
		if state.STATE.CurrentEvent != nil {
			err := nexus.INSTANCE.SyncPictures(state.STATE.CurrentEvent)
			if err != nil {
				log.Error("Failed to sync pictures", "err", err)
			}
		}

		err := nexus.INSTANCE.SyncSongs()
		if err != nil {
			log.Error("Failed to sync songs", "err", err)
		}

		err = nexus.INSTANCE.SyncSessions()
		if err != nil {
			log.Error("Failed to sync songs sessions", "err", err)
		}
	}()

	c.Status(http.StatusOK)
}
