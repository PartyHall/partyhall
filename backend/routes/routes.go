package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/partyhall/partyhall/middlewares"
	"github.com/partyhall/partyhall/models"
	"github.com/partyhall/partyhall/state"
)

const GET_COLLECTION_LIMIT = 50

func RegisterWebappRoutes(router *gin.RouterGroup) {
	router.POST("/login", routeLogin)
	router.POST("/guest-login", routeLoginGuest)
	router.POST("/refresh", routeLoginRefresh)

	router.GET("/status", routeStatus)

	r := router.Group("/webapp")
	//#region CRUD Event
	events := r.Group("/events", middlewares.Authorized(models.ROLE_ADMIN))
	events.POST("", routeCreateEvent)
	events.GET("", routeGetEvents)
	events.GET("/:eventId", routeGetEvent)
	events.PUT("/:eventId", routeUpdateEvent)
	events.DELETE("/events/:eventId", routeDeleteEvent)
	//#endregion

	//#region Karaoke related stuff
	songs := r.Group("/songs")
	songs.GET("", routeGetSongs)
	songs.GET("/:songId", routeGetSong)
	songs.GET("/:songId/data/*filepath", routeServeSong)

	session := r.Group("/session", middlewares.Authorized(models.ROLE_GUEST, models.ROLE_USER, models.ROLE_APPLIANCE))
	session.POST("/", routeQueueAdd)
	session.POST("/:sessionId/start", routeSessionDirectPlay)
	session.POST("/:sessionId/ended", routeSetEnded)
	session.POST("/:sessionId/move/:direction", routeMoveInQueue)
	session.DELETE("/:sessionId", routeQueueRemove)

	karaoke := r.Group("/karaoke", middlewares.Authorized(models.ROLE_GUEST, models.ROLE_USER, models.ROLE_APPLIANCE))
	karaoke.POST("/timecode", routeSetTimecode)
	karaoke.POST("/playing-status/:status", routeSetPlayingStatus)
	karaoke.POST("/set-volume/:type/:volume", routeSetVolume)
	//#endregion

	//#region Backdrop related stuff
	backdrops := r.Group("/backdrops")
	backdrops.GET("", routeGetBackdropAlbums)
	backdrops.GET(":albumId", routeGetBackdropAlbum)
	backdrops.GET(":albumId/image/:backdropId/download", routeDownloadBackdrop)
	//#endregion

	r.GET("/logs", middlewares.Authorized(models.ROLE_ADMIN), routeGetLogs)

	r.POST("/picture", middlewares.HasEventLoaded(), middlewares.Authorized(), routeTakePicture)
	r.POST("/flash", middlewares.HasEventLoaded(), routeSetFlash)

	// Admin
	settings := r.Group("/settings")
	settings.POST("/shutdown", middlewares.Authorized("ADMIN"), routeShutdown)
	settings.POST("/mode/:mode", middlewares.Authorized("ADMIN"), routeSetMode)
	settings.POST("/event/:event", middlewares.Authorized("ADMIN"), routeSetEvent)
	settings.POST("/debug", routeSetDebug)
	settings.POST("/force-sync", middlewares.Authorized("ADMIN"), routeForceSync)
	settings.POST("/backdrops", routeSetBackdrops)

	settings.GET("/audio-devices", routeGetAudioDevices)
	settings.POST("/audio-devices", middlewares.Authorized("ADMIN"), routeSetAudioDevices)
	settings.POST("/audio-devices/:id/volume", routeSetAudioDeviceVolume)

	nexusRoutes := r.Group("/nexus", middlewares.HasEventLoaded(), middlewares.Authorized("ADMIN"))
	nexusRoutes.POST("/sync", routeSync)
	nexusRoutes.POST("/events/:id", routeCreateOnNexus)
}

func RegisterApplianceRoutes(router *gin.RouterGroup) {
	r := router.Group("/appliance", middlewares.Authorized("APPLIANCE"))
	r.POST("/picture", middlewares.HasEventLoaded(), routeUploadPicture)
}

func routeStatus(c *gin.Context) {
	c.JSON(200, state.STATE)
}
