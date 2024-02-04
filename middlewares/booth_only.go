package middlewares

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/partyhall/partyhall/config"
	"github.com/partyhall/partyhall/utils"
)

func BoothOnlyMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if utils.IsRemote(c) {
			if !config.GET.DebugMode && !config.IsInDev() {
				return echo.NewHTTPError(http.StatusForbidden, "This can only be called from the appliance itself")
			}
		}

		return next(c)
	}
}
