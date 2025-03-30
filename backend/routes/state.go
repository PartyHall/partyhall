package routes

import (
	"database/sql"
	"errors"
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/partyhall/partyhall/api_errors"
	"github.com/partyhall/partyhall/dal"
	"github.com/partyhall/partyhall/mercure_client"
	"github.com/partyhall/partyhall/middlewares"
	"github.com/partyhall/partyhall/mqtt"
	"github.com/partyhall/partyhall/state"
)

type RoutesState struct{}

func (h RoutesState) Register(router *gin.RouterGroup) {
	// Unauthenticated
	router.GET("", h.getState)

	// Unauthenticated
	router.POST("debug", h.showDebug)

	// Onboarded & Admin
	router.PUT(
		"event",
		middlewares.Onboarded(true),
		middlewares.Authorized("ADMIN"),
		h.setEvent,
	)

	// Onboarded & Admin
	router.PUT(
		"mode",
		middlewares.Onboarded(true),
		middlewares.Authorized("ADMIN"),
		h.setMode,
	)

	// Onboarded & Authenticated
	router.PUT(
		"backdrops",
		middlewares.Onboarded(true),
		middlewares.Authorized(),
		h.setBackdrops,
	)

	// Onboarded & Authenticated
	router.PUT(
		"flash",
		middlewares.Onboarded(true),
		middlewares.Authorized(),
		h.setFlash,
	)
}

func (h RoutesState) getState(c *gin.Context) {
	c.JSON(200, state.STATE)
}

func (h RoutesState) showDebug(c *gin.Context) {
	err := mercure_client.CLIENT.ShowDebug()
	if err != nil {
		c.Render(http.StatusInternalServerError, api_errors.MERCURE_PUBLISH_FAILURE)

		return
	}

	c.Status(http.StatusOK)
}

func (h RoutesState) setEvent(c *gin.Context) {
	var req struct {
		EventId int64 `json:"event"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		api_errors.RenderValidationErr(c, err)
		return
	}

	event, err := dal.EVENTS.Get(req.EventId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.Render(http.StatusNotFound, api_errors.NOT_FOUND)
		}

		c.Render(http.StatusInternalServerError, api_errors.DATABASE_ERROR.WithExtra(map[string]any{
			"err": err,
		}))

		return
	}

	if state.STATE.CurrentEvent == nil || state.STATE.CurrentEvent.Id != event.Id {
		state.STATE.CurrentEvent = event

		if err := mercure_client.CLIENT.SetCurrentEvent(event); err != nil {
			c.Render(http.StatusInternalServerError, api_errors.MERCURE_PUBLISH_FAILURE)

			return
		}

		dal.EVENTS.Set(event)
	}

	c.JSON(http.StatusOK, event)
}

func (h RoutesState) setMode(c *gin.Context) {
	var req struct {
		Mode string `json:"mode" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		api_errors.RenderValidationErr(c, err)
		return
	}

	mode := strings.ToLower(req.Mode)

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

		mercure_client.CLIENT.SetMode(state.STATE.CurrentMode)

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

		if err := mercure_client.CLIENT.SetMode(mode); err != nil {
			c.Render(http.StatusInternalServerError, api_errors.MERCURE_PUBLISH_FAILURE)
			return
		}

	}

	c.Status(http.StatusNoContent)
}

func (h RoutesState) setBackdrops(c *gin.Context) {
	var req struct {
		BackdropAlbum    *int64 `json:"backdrop_album"`
		SelectedBackdrop int    `json:"selected_backdrop"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		api_errors.RenderValidationErr(c, err)

		return
	}

	if req.BackdropAlbum != nil {
		alb, err := dal.BACKDROPS.GetAlbum(*req.BackdropAlbum)
		if err != nil {
			api_errors.RenderValidationErr(c, err)

			return
		}

		state.STATE.SelectedBackdrop = req.SelectedBackdrop

		if state.STATE.BackdropAlbum == nil || alb.Id != state.STATE.BackdropAlbum.Id {
			state.STATE.SelectedBackdrop = 0
		}

		state.STATE.BackdropAlbum = &alb
	} else {
		state.STATE.BackdropAlbum = nil
		state.STATE.SelectedBackdrop = 0
	}

	mercure_client.CLIENT.SendBackdropState()

	c.JSON(200, map[string]any{
		"backdrop_album":    state.STATE.BackdropAlbum,
		"selected_backdrop": state.STATE.SelectedBackdrop,
	})
}

func (h RoutesState) setFlash(c *gin.Context) {
	var req struct {
		Powered bool `json:"powered"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		api_errors.RenderValidationErr(c, err)

		return
	}

	mqtt.SetFlash(req.Powered, state.STATE.UserSettings.Photobooth.FlashBrightness)

	c.Status(200)
}
