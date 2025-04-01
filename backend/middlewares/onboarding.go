package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/partyhall/partyhall/api_errors"
	"github.com/partyhall/partyhall/config"
)

func Onboarded(shouldBeOnboarded bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		isOnboarded := config.GET.UserSettings.Onboarded

		if shouldBeOnboarded && !isOnboarded {
			c.Render(http.StatusBadRequest, api_errors.SHOULD_BE_ONBOARDED)
			c.Abort()
		}

		if !shouldBeOnboarded && isOnboarded {
			c.Render(http.StatusBadRequest, api_errors.ALREADY_ONBOARDED)
			c.Abort()
		}

		c.Next()
	}
}

func NotOnboardedOrRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		isOnboarded := config.GET.UserSettings.Onboarded

		if !isOnboarded {
			c.Next()

			return
		}

		authorizedMiddleware := Authorized(roles...)

		authorizedMiddleware(c)
	}
}
