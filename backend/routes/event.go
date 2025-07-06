package routes

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/partyhall/partyhall/api_errors"
	"github.com/partyhall/partyhall/config"
	"github.com/partyhall/partyhall/dal"
	"github.com/partyhall/partyhall/mercure_client"
	"github.com/partyhall/partyhall/middlewares"
	"github.com/partyhall/partyhall/models"
	"github.com/partyhall/partyhall/state"
	"github.com/partyhall/partyhall/utils"
)

type upsertEventRequest struct {
	Id              int                        `json:"id"`
	Name            string                     `json:"name" binding:"required"`
	Author          string                     `json:"author"`
	Date            string                     `json:"date" binding:"iso8601,required"`
	Location        string                     `json:"location"`
	NexusID         models.JsonnableNullstring `json:"nexus_id"`
	RegistrationUrl models.JsonnableNullstring `json:"registration_url"`
}

type RoutesEvent struct{}

func (h RoutesEvent) Register(router *gin.RouterGroup) {
	// Not onboarded or Admin
	router.GET(
		"",
		middlewares.Authorized(models.ROLE_ADMIN),
		h.getCollection,
	)

	// Not onboarded or Admin
	router.GET(
		":eventId",
		middlewares.Authorized(models.ROLE_ADMIN),
		h.get,
	)

	// Anyone
	router.GET(
		":eventId/registration-qr",
		h.getQr,
	)

	// Not onboarded or Admin
	router.POST(
		"",
		middlewares.Authorized(models.ROLE_ADMIN),
		h.create,
	)

	// Not onboarded or Admin
	router.PUT(
		":eventId",
		middlewares.Authorized(models.ROLE_ADMIN),
		h.update,
	)

	router.DELETE(
		":eventId",
		middlewares.Authorized(models.ROLE_ADMIN),
		h.delete,
	)
}

func (h RoutesEvent) getCollection(c *gin.Context) {
	page, offset, err := utils.ParsePageOffset(c)
	if err != nil {
		return
	}

	events, err := dal.EVENTS.GetCollection(config.AMT_RESULTS_PER_PAGE, offset)
	if err != nil {
		c.Render(http.StatusInternalServerError, api_errors.DATABASE_ERROR.WithExtra(map[string]any{
			"err": err.Error(),
		}))

		return
	}

	events.Page = page

	c.JSON(http.StatusOK, events)
}

func (h RoutesEvent) get(c *gin.Context) {
	id, parseErr := utils.ParamAsIntOrError(c, "eventId")
	if parseErr {
		return
	}

	evt, err := dal.EVENTS.Get(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.Status(http.StatusNotFound)
			return
		}

		c.Render(http.StatusInternalServerError, api_errors.DATABASE_ERROR.WithExtra(map[string]any{
			"err": err,
		}))

		return
	}

	c.JSON(http.StatusOK, evt)
}

func (h RoutesEvent) getQr(c *gin.Context) {
	id, parseErr := utils.ParamAsIntOrError(c, "eventId")
	if parseErr {
		return
	}

	evt, err := dal.EVENTS.Get(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.Status(http.StatusNotFound)
			return
		}

		c.Render(http.StatusInternalServerError, api_errors.DATABASE_ERROR.WithExtra(map[string]any{
			"err": err,
		}))

		return
	}

	if !evt.UserRegistrationUrl.Valid || len(evt.UserRegistrationUrl.String) == 0 {
		c.Render(http.StatusBadRequest, api_errors.INVALID_PARAMETERS.WithExtra(map[string]any{
			"eventId": "This event does not have a registration URL set-up",
		}))

		return
	}

	utils.GenerateQrCodeWithoutLogo(evt.UserRegistrationUrl.String, c)
}

func (h RoutesEvent) create(c *gin.Context) {
	var req upsertEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api_errors.RenderValidationErr(c, err)
		return
	}

	date, _ := time.Parse(time.RFC3339, req.Date)

	evt := models.Event{
		Name:                req.Name,
		Author:              req.Author,
		Date:                date,
		Location:            req.Location,
		UserRegistrationUrl: req.RegistrationUrl,
	}

	err := dal.EVENTS.Create(&evt)
	if err != nil {
		c.Render(http.StatusInternalServerError, api_errors.DATABASE_ERROR.WithExtra(map[string]any{
			"err": err.Error(),
		}))

		return
	}

	/** If no event was set-up, then this event is the current one */
	if state.STATE.CurrentEvent == nil {
		state.STATE.CurrentEvent = &evt

		mercure_client.CLIENT.SetCurrentEvent(state.STATE.CurrentEvent)
		dal.EVENTS.Set(state.STATE.CurrentEvent)
	}

	c.JSON(200, evt)
}

func (h RoutesEvent) update(c *gin.Context) {
	var req upsertEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api_errors.RenderValidationErr(c, err)
		return
	}

	date, _ := time.Parse(time.RFC3339, req.Date)

	/** @TODO: No, the id should come from the queryParams, not the body **/
	evt := models.Event{
		Id:                  int64(req.Id),
		Name:                req.Name,
		Author:              req.Author,
		Date:                date,
		Location:            req.Location,
		NexusId:             req.NexusID,
		UserRegistrationUrl: req.RegistrationUrl,
	}

	err := dal.EVENTS.Update(&evt)
	if err != nil {
		c.Render(http.StatusInternalServerError, api_errors.DATABASE_ERROR.WithExtra(map[string]any{
			"err": err.Error(),
		}))

		return
	}

	if state.STATE.CurrentEvent != nil && state.STATE.CurrentEvent.Id == evt.Id {
		state.STATE.CurrentEvent = &evt
	}

	c.JSON(200, evt)
}

func (h RoutesEvent) delete(c *gin.Context) {
	id, parseErr := utils.ParamAsIntOrError(c, "eventId")
	if parseErr {
		return
	}

	// Should never be nil as we check for a current event in the routes
	// But as a safeguard
	if state.STATE.CurrentEvent == nil || state.STATE.CurrentEvent.Id == id {
		c.Render(http.StatusBadRequest, api_errors.INVALID_PARAMETERS.WithExtra(map[string]any{
			"eventId": "You cannot delete the current event",
		}))

		return
	}

	err := dal.EVENTS.Delete(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.Status(http.StatusNotFound)
			return
		}

		c.Render(http.StatusInternalServerError, api_errors.DATABASE_ERROR.WithExtra(map[string]any{
			"err": err,
		}))

		return
	}

	c.Status(http.StatusNoContent)
}
