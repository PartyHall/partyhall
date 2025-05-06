package routes

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/partyhall/partyhall/api_errors"
	"github.com/partyhall/partyhall/mercure_client"
	"github.com/partyhall/partyhall/middlewares"
	"github.com/partyhall/partyhall/models"
	"github.com/partyhall/partyhall/nexus"
	"github.com/partyhall/partyhall/state"
	"github.com/partyhall/partyhall/utils"
)

type RoutesNexus struct{}

func (h RoutesNexus) Register(router *gin.RouterGroup) {
	// Has an event & Admin
	router.POST(
		"create_event/:eventId",
		middlewares.HasEventLoaded(),
		middlewares.Authorized(models.ROLE_ADMIN),
		h.createEventOnNexus,
	)

	// Admin
	router.POST(
		"sync",
		middlewares.Authorized(models.ROLE_ADMIN),
		h.sync,
	)
}

func (h RoutesNexus) createEventOnNexus(c *gin.Context) {
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

	eventId, parseFailed := utils.ParamAsIntOrError(c, "eventId")
	if parseFailed {
		return
	}

	err := nexus.INSTANCE.CreateEvent(eventId)
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

	// @WARN wtf is that ?
	// Do we really want that creating an event on nexus sets it as the current event?
	// I don't think so but I won't touch it yet
	// We probably just want to send the JSON of the given event
	mercure_client.CLIENT.SetCurrentEvent(state.STATE.CurrentEvent)
	c.JSON(http.StatusOK, state.STATE.CurrentEvent)
}

func (h RoutesNexus) sync(c *gin.Context) {
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

	// Do we really want to sync only the current event?
	// Should an admin not be able to force-sync another event?
	err := nexus.INSTANCE.Sync(state.STATE.CurrentEvent)
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
