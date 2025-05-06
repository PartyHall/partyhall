package routes

import (
	"github.com/gin-gonic/gin"
)

/**
	Summary of routes and permissions

	Definitions:
		- authenticated: Logged in anonymously (temp user) or as a real user
		- appliance: Only callable from localhost (the appliance itself)
		- onboarded: The admin has completed the initial setup process
		- hasEvent: An event was created and is currently selected

		Note: The appliance frontend has the USER, ADMIN, APPLIANCE roles

	@TODO
	Nouvelle règle:
	    - Onboarded no longer exists as we have onboarding for multiple stuff
		- We should re-check to be sure that EVENT IS LOADED before most of the routes

	#region API:
	/
		login          onboarded & unauthenticated
		guest-login    onboarded & unauthenticated
		refresh        onboarded & unauthenticated

	/state
		GET /                 unauthenticated
		POST /debug           unauthenticated
		PUT /event            onboarded & admin
		PUT /mode             onboarded & admin
		PUT /backdrops        onboarded & authenticated
		GET /flash            unauthenticated => Get the state of the flash (brightness + powered)
		PUT /flash            onboarded & authenticated // Only sets whether it is on or off, used by appliance AND backend

	/events
		GET /          non-onboarded | Admin
		GET /:id       non-onboarded | Admin
		POST /         non-onboarded | Admin
		PATCH /:id     non-onboarded | Admin
		DELETE /:id    onboarded & Admin

	/photobooth
		POST /take-picture    onboarded & hasEvent & authenticated
		POST /upload-picture  onboarded & hasEvent & Appliance

	/backdrops_albums
		GET /                       onboarded & authenticated
		GET /:backdropAlbumId       onboarded & authenticated

	/backdrops
		GET /:backdropId/download   onboarded // No authenticated check because I'm too lazy to handle the auth in img tags

	/songs
		GET /                            onboarded & Authenticated
		GET /:songId                     onboarded & Authenticated
		GET /:songId/cover               onboarded // No authenticated check because I'm too lazy to handle the auth in img tags
		GET /:songId/file/:songFilename  onboarded & Appliance

	/song_sessions
		POST   /                               onboarded & hasEvent & Authenticated
		POST   /:songSessionId/start           onboarded & hasEvent & Authenticated
		POST   /:songSessionId/ended           onboarded & hasEvent & Appliance (Or Authenticated ? Need to check how I did this)
		POST   /:songSessionId/move/:direction onboarded & hasEvent & Authenticated
		DELETE /:songSessionId                 onboarded & hasEvent & Authenticated

	/karaoke
		PUT /timecode       onboarded & hasEvent & appliance
		PUT /playing_status onboarded & hasEvent & authenticated
		PUT /volume         onboarded & hasEvent & authenticated

	/nexus
		POST /sync          onboarded & admin
		POST /events/:id    onboarded & hasEvent & admin

	/admin
		POST /create-admin           Only if no admin was created
		GET /logs                    onboarded & admin
		POST /shutdown               onboarded & admin

	/settings non-onboarded | admin
		GET /ap
		PUT /ap
		GET /webcam // Not sure if useful
		PUT /webcam
		POST /flash                  => Set the default brightness value
		GET /photobooth
		PUT /photobooth
		GET /audio-devices
		POST /audio-devices
		PUT /audio-devices/:deviceId => Set the volume for a device
		GET /spotify
		PUT /spotify
		GET /nexus
		PUT /nexus
		GET /physical-buttons
		PUT /physical-buttons
		POST create-admin           non-onboarded only
		POST conclude-onboarding    non-onboarded only
	#endregion

	#region Mercure
	Topic             | Frontends   | Permissions
	------------------|-------------|-------------
	/time             | admin front | everyone
	/event            | admin front | everyone
	/flash            | admin       | admin
	/mode             | admin front | Pourquoi pas l'admin?
	/sync-progress    | admin       | everyone (?)
	/snackbar         | admin front | everyone
	/backdrop-state   | admin front | everyone
	/ip-addresses     |       front | appliance
	/debug            |       front | appliance
	/take-picture     |       front | appliance
	/karaoke-queue    |       front | everyone (Pourquoi l'admin ne l'a pas ?)
	/audio-devices    | admin       | admin
	/karaoke          | admin       | everyone
	/karaoke-queue    | admin       | everyone
	/karaoke-timecode | admin       | everyone
	/logs             | admin       | admin
	#endregion
**/

func RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/login", routeLogin)
	router.POST("/guest-login", routeLoginGuest)
	router.POST("/refresh", routeLoginRefresh)

	// Special endpoints
	(RoutesState{}).Register(router.Group("/state"))
	(RoutesPhotobooth{}).Register(router.Group("/photobooth"))
	(RoutesKaraoke{}).Register(router.Group("/karaoke"))
	(RoutesNexus{}).Register(router.Group("/nexus"))
	(RoutesSettings{}).Register(router.Group("/settings"))
	(RoutesAdmin{}).Register(router.Group("/admin"))

	// Standard CRUD++ endpoint
	(RoutesEvent{}).Register(router.Group("/events"))
	(RoutesBackdropAlbums{}).Register(router.Group("/backdrop_albums"))
	(RoutesBackdrops{}).Register(router.Group("/backdrops"))
	(RoutesSong{}).Register(router.Group("/songs"))
	(RoutesSongSession{}).Register(router.Group("/song_sessions"))
	(RoutesUser{}).Register(router.Group("/users"))
}
