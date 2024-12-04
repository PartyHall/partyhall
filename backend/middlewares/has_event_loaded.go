package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/partyhall/partyhall/api_errors"
	"github.com/partyhall/partyhall/state"
)

func HasEventLoaded() gin.HandlerFunc {
	return func(c *gin.Context) {
		if state.STATE.CurrentEvent == nil {
			api_errors.ApiError(
				c,
				http.StatusBadRequest,
				"no-event",
				"No event selected",
				"Before using PartyHall you should ensure that an event is created and selected",
			)
			c.Abort()

			return
		}

		c.Next()
	}
}
