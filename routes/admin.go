package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func registerAdminRoutes(g *echo.Group) {
	// r.Use(AuthenticatedMiddleware)

	g.POST("/login", isPasswordValid)

	g.POST("/event", eventPost)
	g.GET("/event/:id", eventGet)
	g.PUT("/event/:id", eventPut)

	g.GET("/event/:id/export", eventExportsGet)
	g.GET("/event/:id/export/download", eventExportsDownload)
	// g.POST("/event/:id/export", eventExportsPost) // Start an export

	// r.HandleFunc("/set_mode/{mode}", setMode).Methods(http.MethodPost) // Unused, but @TODO since we want to remove the most from WS
}

// @TODO: Oops
func isPasswordValid(c echo.Context) error {
	return c.String(http.StatusOK, "yes")
}
