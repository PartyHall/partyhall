package routes

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/partyhall/partyhall/api_errors"
	"github.com/partyhall/partyhall/mercure_client"
	"github.com/partyhall/partyhall/nexus"
	"github.com/partyhall/partyhall/state"
)

func routeCreateOnNexus(c *gin.Context) {
	if !nexus.INSTANCE.IsSetup {
		api_errors.ApiError(
			c,
			http.StatusBadRequest,
			api_errors.NEXUS_ERR_TYPE,
			"PartyNexus is not setup",
			"You are trying to synchronise with PartyNexus but it is not setup properly",
		)

		return
	}

	eventIdStr := c.Params.ByName("id")

	eventId, err := strconv.ParseInt(eventIdStr, 10, 64)
	if err != nil {
		c.Render(http.StatusBadRequest, api_errors.BAD_REQUEST.WithExtra(map[string]any{
			"err": "Bad event id",
		}))

		return
	}

	err = nexus.INSTANCE.CreateEvent(eventId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			api_errors.ApiError(
				c,
				http.StatusNotFound,
				api_errors.NEXUS_ERR_TYPE,
				"This event does not exist",
				fmt.Sprintf("No event with id %v", eventId),
			)
		}

		api_errors.ApiError(
			c,
			http.StatusInternalServerError,
			api_errors.NEXUS_ERR_TYPE,
			"Failed to create event on PartyNexus",
			err.Error(),
		)

		return
	}

	mercure_client.CLIENT.PublishEvent("/event", state.STATE.CurrentEvent)
	c.JSON(http.StatusOK, state.STATE.CurrentEvent)
}

func routeSync(c *gin.Context) {
	if !nexus.INSTANCE.IsSetup {
		api_errors.ApiError(
			c,
			http.StatusBadRequest,
			api_errors.NEXUS_ERR_TYPE,
			"PartyNexus is not setup",
			"You are trying to synchronise with PartyNexus but it is not setup properly",
		)

		return
	}

	err := nexus.INSTANCE.SyncPictures(state.STATE.CurrentEvent)
	if err != nil {
		api_errors.ApiError(
			c,
			http.StatusInternalServerError,
			api_errors.NEXUS_ERR_TYPE,
			"Failed to synchronise with PartyNexus",
			"An error occured while synchronizing with PartyNexus: "+err.Error(),
		)

		return
	}

	c.Status(http.StatusOK)
}
