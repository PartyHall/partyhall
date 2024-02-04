package routes

import (
	"github.com/labstack/echo/v4"
)

func registerAdminRoutes(g *echo.Group) {
	// r.Use(AuthenticatedMiddleware)

	g.POST("/event", eventPost)
	g.GET("/event/:id", eventGet)
	g.PUT("/event/:id", eventPut)

	g.GET("/event/:id/export", eventExportsGet)
	g.GET("/event/:id/export/download", eventExportsDownload)
	// g.POST("/event/:id/export", eventExportsPost) // Start an export
}
