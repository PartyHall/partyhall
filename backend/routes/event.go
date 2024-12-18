package routes

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/partyhall/partyhall/api_errors"
	"github.com/partyhall/partyhall/config"
	"github.com/partyhall/partyhall/dal"
	"github.com/partyhall/partyhall/mercure_client"
	"github.com/partyhall/partyhall/models"
	routes_requests "github.com/partyhall/partyhall/routes/requests"
	"github.com/partyhall/partyhall/state"
	"github.com/partyhall/partyhall/utils"
)

func routeGetEvents(c *gin.Context) {
	offset := 0

	page := c.Query("page")

	var pageInt int = 1
	var err error
	if len(page) > 0 {
		pageInt, err = strconv.Atoi(page)
		if err != nil {
			c.Render(http.StatusBadRequest, api_errors.INVALID_PARAMETERS.WithExtra(map[string]any{
				"page": "The page should be an integer",
			}))

			return
		}

		offset = (pageInt - 1) * config.AMT_RESULTS_PER_PAGE
	}

	events, err := dal.EVENTS.GetCollection(config.AMT_RESULTS_PER_PAGE, offset)
	if err != nil {
		c.Render(http.StatusInternalServerError, api_errors.DATABASE_ERROR.WithExtra(map[string]any{
			"err": err.Error(),
		}))

		return
	}

	events.Page = pageInt

	c.JSON(http.StatusOK, events)
}

func routeGetEvent(c *gin.Context) {
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

func routeDeleteEvent(c *gin.Context) {
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

func routeCreateEvent(c *gin.Context) {
	var req routes_requests.CreateEvent
	if err := c.ShouldBindJSON(&req); err != nil {
		api_errors.RenderValidationErr(c, err)
		return
	}

	date, _ := time.Parse(time.RFC3339, req.Date)

	evt := models.Event{
		Name:     req.Name,
		Author:   req.Author,
		Date:     date,
		Location: req.Location,
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
		mercure_client.CLIENT.PublishEvent("/event", evt)
		dal.EVENTS.Set(&evt)
	}

	c.JSON(200, evt)
}

func routeUpdateEvent(c *gin.Context) {
	var req routes_requests.CreateEvent
	if err := c.ShouldBindJSON(&req); err != nil {
		api_errors.RenderValidationErr(c, err)
		return
	}

	date, _ := time.Parse(time.RFC3339, req.Date)

	evt := models.Event{
		Id:       int64(req.Id),
		Name:     req.Name,
		Author:   req.Author,
		Date:     date,
		Location: req.Location,
		NexusId:  req.NexusID,
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
