package routes

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/partyhall/partyhall/dto"
	"github.com/partyhall/partyhall/logs"
	"github.com/partyhall/partyhall/orm"
	"github.com/partyhall/partyhall/remote"
	"github.com/partyhall/partyhall/services"
	"github.com/partyhall/partyhall/utils"
)

func eventPost(c echo.Context) error {
	event := new(dto.EventPost)
	if err := c.Bind(event); err != nil {
		fmt.Println(err)
		return c.String(http.StatusBadRequest, "bad request")
	}

	if err := c.Validate(event); err != nil {
		return err
	}

	dbEvent, err := orm.GET.Events.CreateEvent(event)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to create event: "+err.Error())
	}

	//#region Setting it as the event if its the only one
	events, err := orm.GET.Events.GetEvents()
	if err == nil {
		if len(events) == 1 {
			services.GET.CurrentState.CurrentEvent = &events[0].Id
			services.GET.CurrentState.CurrentEventObj = &events[0]

			orm.GET.AppState.SetState(services.GET.CurrentState)
		}
	}
	//#endregion

	remote.BroadcastState()

	return c.JSON(http.StatusCreated, dto.GetEvent(dbEvent))
}

func eventPut(c echo.Context) error {
	eventId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	event := new(dto.EventPut)
	if err := c.Bind(event); err != nil {
		fmt.Println(err)
		return c.String(http.StatusBadRequest, "bad request")
	}

	if err := c.Validate(event); err != nil {
		return err
	}

	dbEvent, err := orm.GET.Events.Update(eventId, event)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to update event: "+err.Error())
	}

	return c.JSON(http.StatusOK, dto.GetEvent(dbEvent))
}

func eventGet(c echo.Context) error {
	eventId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	dbEvent, err := orm.GET.Events.GetEvent(eventId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to update event: "+err.Error())
	}

	return c.JSON(http.StatusOK, dto.GetEvent(dbEvent))
}

func eventExportsGet(c echo.Context) error {
	eventId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	event, err := orm.GET.Events.GetEvent(eventId)
	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	exportedEvents, err := orm.GET.Events.GetExportedEvents(event, 5)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, exportedEvents)
}

func eventExportsDownload(c echo.Context) error {
	eventId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	exportedEvent, err := orm.GET.Events.GetExportedEvent(eventId)
	if err != nil {
		logs.Error(err)
		return c.NoContent(http.StatusNotFound)
	}

	path := utils.GetPath("images", fmt.Sprintf("%v", exportedEvent.EventId), "exports", exportedEvent.Filename)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return c.NoContent(http.StatusNotFound)
	}

	return c.File(path)
}
