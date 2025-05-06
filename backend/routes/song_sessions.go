package routes

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/partyhall/partyhall/api_errors"
	"github.com/partyhall/partyhall/dal"
	"github.com/partyhall/partyhall/log"
	"github.com/partyhall/partyhall/mercure_client"
	"github.com/partyhall/partyhall/middlewares"
	"github.com/partyhall/partyhall/models"
	"github.com/partyhall/partyhall/nexus"
	"github.com/partyhall/partyhall/services"
	"github.com/partyhall/partyhall/state"
	"github.com/partyhall/partyhall/utils"
)

type RoutesSongSession struct{}

func (h RoutesSongSession) Register(router *gin.RouterGroup) {
	// Has event & Authenticated
	router.POST(
		"",
		middlewares.HasEventLoaded(),
		middlewares.Authorized(),
		h.create,
	)

	// Has event & Authenticated
	router.POST(
		":songSessionId/start",
		middlewares.HasEventLoaded(),
		middlewares.Authorized(),
		h.start,
	)

	// Has event & Appliance
	router.POST(
		":songSessionId/ended",
		middlewares.HasEventLoaded(),
		middlewares.Authorized(models.ROLE_APPLIANCE),
		h.ended,
	)

	// Has event & Authenticated
	router.POST(
		":songSessionId/move/:direction",
		middlewares.HasEventLoaded(),
		middlewares.Authorized(),
		h.moveInQueue,
	)

	// Has event & Authenticated
	router.DELETE(
		":songSessionId",
		middlewares.HasEventLoaded(),
		middlewares.Authorized(),
		h.delete,
	)
}

func (h RoutesSongSession) create(c *gin.Context) {
	var req struct {
		SongId      string `json:"song_id" binding:"required"`
		DisplayName string `json:"display_name" binding:"required"`
		DirectPlay  bool   `json:"direct_play"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		api_errors.RenderValidationErr(c, err)
		return
	}

	req.SongId = strings.ToLower(req.SongId)

	currentSong := state.STATE.Karaoke.Current
	if currentSong != nil && strings.ToLower(currentSong.NexusId) == req.SongId {
		c.Status(409)

		return
	}

	for _, session := range state.STATE.KaraokeQueue {
		if strings.ToLower(session.NexusId) == req.SongId {
			c.Status(409)

			return
		}
	}

	song, err := dal.SONGS.Get(req.SongId)
	if err != nil {
		log.Error("Failed to find song to add in queue", "err", err)

		c.Render(http.StatusBadRequest, api_errors.BAD_REQUEST.WithExtra(map[string]any{
			"err": "Song does not exists",
		}))

		return
	}

	ctxClaims, _ := c.Get("TokenClaims")
	claims := ctxClaims.(*models.JwtCustomClaims)

	log.Info(
		"Song requested",
		"username", claims.Username,
		"displayName", req.DisplayName,
		"direct", req.DirectPlay,
		"song", song.NexusId,
		"title", fmt.Sprintf("%v by %v", song.Title, song.Artist),
	)

	session := models.SongSession{
		EventId:  state.STATE.CurrentEvent.Id,
		Title:    song.Title,
		Artist:   song.Artist,
		SungBy:   req.DisplayName,
		SungById: fmt.Sprintf("%v", claims.Subject),
		Song:     song,
		AddedAt: models.JsonnableNullTime{
			Time:  time.Now(),
			Valid: true,
		},
		StartedAt:   models.JsonnableNullTime{},
		EndedAt:     models.JsonnableNullTime{},
		CancelledAt: models.JsonnableNullTime{},
	}

	directPlay := state.STATE.Karaoke.Current == nil || req.DirectPlay

	err = dal.SONGS.CreateSession(&session)
	if err != nil {
		log.Error("Failed to add in queue", "err", err)

		c.Render(http.StatusInternalServerError, api_errors.DATABASE_ERROR.WithExtra(map[string]any{
			"err": err,
		}))

		return
	}

	if directPlay {
		if err := services.KARAOKE.StartSong(&session); err != nil {
			log.Error("Failed to direct play the song: ", "err", err)
			c.Render(http.StatusInternalServerError, api_errors.DATABASE_ERROR.WithExtra(map[string]any{
				"err": err,
			}))

			return
		}
	} else {
		state.STATE.KaraokeQueue = append(state.STATE.KaraokeQueue, &session)

		mercure_client.CLIENT.SendKaraokeQueue()
	}

	c.JSON(200, session)
}

func (h RoutesSongSession) start(c *gin.Context) {
	sessionId, parseFailed := utils.ParamAsIntOrError(c, "songSessionId")
	if parseFailed {
		return
	}

	var idx = -1
	var session *models.SongSession = nil
	for i, s := range state.STATE.KaraokeQueue {
		if s.Id == sessionId {
			session = s
			idx = i
			break
		}
	}

	if session == nil {
		c.Render(http.StatusBadRequest, api_errors.BAD_REQUEST.WithExtra(map[string]any{
			"err": "Session ID is invalid (not in queue)",
		}))

		return
	}

	ctxClaims, _ := c.Get("TokenClaims")
	claims := ctxClaims.(*models.JwtCustomClaims)
	log.Info(
		"Directplay from queue",
		"username", claims.Username,
		"session", session.Id,
		"song", session.Song.NexusId,
		"title", fmt.Sprintf("%v by %v", session.Song.Title, session.Song.Artist),
	)

	state.STATE.KaraokeQueue = append(state.STATE.KaraokeQueue[:idx], state.STATE.KaraokeQueue[idx+1:]...)

	mercure_client.CLIENT.SendKaraokeQueue()

	if err := services.KARAOKE.StartSong(session); err != nil {
		log.Error("Failed to direct play the song: ", "err", err)
		c.Render(http.StatusInternalServerError, api_errors.DATABASE_ERROR.WithExtra(map[string]any{
			"err": err,
		}))

		return
	}

	c.JSON(200, session)
}

func (h RoutesSongSession) ended(c *gin.Context) {
	sessionId, parseFailed := utils.ParamAsIntOrError(c, "songSessionId")
	if parseFailed {
		return
	}

	current := state.STATE.Karaoke.Current

	if current == nil || current.Id != sessionId {
		c.Render(http.StatusBadRequest, api_errors.BAD_REQUEST.WithExtra(map[string]any{
			"err": "Session ID is invalid",
		}))

		return
	}

	current.EndedAt = models.JsonnableNullTime{
		Time:  time.Now(),
		Valid: true,
	}

	log.Info(
		"Song ended",
		"song", current.Song.NexusId,
		"title", fmt.Sprintf("%v by %v", current.Song.Title, current.Song.Artist),
	)

	err := dal.SONGS.UpdateSession(current)
	if err != nil {
		c.Render(http.StatusInternalServerError, api_errors.DATABASE_ERROR.WithExtra(map[string]any{
			"err": err,
		}))

		return
	}

	err = nexus.INSTANCE.Sync(state.STATE.CurrentEvent)
	if err != nil {
		log.Error("Failed to sync", "err", err)
	}

	state.STATE.Karaoke.Current = nil

	services.KARAOKE.StartNextSong()
}

func (h RoutesSongSession) moveInQueue(c *gin.Context) {
	sessionId, parseFailed := utils.ParamAsIntOrError(c, "songSessionId")
	if parseFailed {
		return
	}

	direction := strings.ToLower(c.Param("direction"))

	if direction != "up" && direction != "down" {
		c.Render(http.StatusBadRequest, api_errors.BAD_REQUEST.WithExtra(map[string]any{
			"err": "The direction should be either up or down",
		}))

		return
	}

	sessionIdx := -1

	for idx, session := range state.STATE.KaraokeQueue {
		if session.Id == sessionId {
			sessionIdx = idx
			break
		}
	}

	if sessionIdx < 0 {
		c.Render(http.StatusBadRequest, api_errors.BAD_REQUEST.WithExtra(map[string]any{
			"err": "The session was not found in the queue",
		}))

		return
	}

	if len(state.STATE.KaraokeQueue) < 2 {
		c.Render(http.StatusBadRequest, api_errors.BAD_REQUEST.WithExtra(map[string]any{
			"err": "Not enough sessions in the queue to perform movement",
		}))

		return
	}

	newPos := sessionIdx - 1
	if direction == "up" {
		if sessionIdx == 0 {
			c.Render(http.StatusBadRequest, api_errors.BAD_REQUEST.WithExtra(map[string]any{
				"err": "Cannot move the first session up",
			}))

			return
		}

	} else if direction == "down" {
		if sessionIdx == len(state.STATE.KaraokeQueue)-1 {
			c.Render(http.StatusBadRequest, api_errors.BAD_REQUEST.WithExtra(map[string]any{
				"err": "Cannot move the last session down",
			}))

			return
		}

		newPos = sessionIdx + 1
	}

	ctxClaims, _ := c.Get("TokenClaims")
	claims := ctxClaims.(*models.JwtCustomClaims)
	log.Info(
		"Session moved in queue",
		"username", claims.Username,
		"session", state.STATE.KaraokeQueue[sessionIdx].Id,
		"song", state.STATE.KaraokeQueue[sessionIdx].Song.NexusId,
		"title", fmt.Sprintf("%v by %v", state.STATE.KaraokeQueue[sessionIdx].Song.Title, state.STATE.KaraokeQueue[sessionIdx].Song.Artist),
		"previousPos", sessionIdx,
		"newPos", newPos,
	)

	state.STATE.KaraokeQueue[sessionIdx], state.STATE.KaraokeQueue[newPos] = state.STATE.KaraokeQueue[newPos], state.STATE.KaraokeQueue[sessionIdx]

	mercure_client.CLIENT.SendKaraokeQueue()

	c.JSON(200, state.STATE.KaraokeQueue)
}

func (h RoutesSongSession) delete(c *gin.Context) {
	sessionId := c.Param("songSessionId")

	session, err := dal.SONGS.GetSession(sessionId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.Status(404)
			return
		}

		c.Render(http.StatusInternalServerError, api_errors.DATABASE_ERROR.WithExtra(map[string]any{
			"err": err.Error(),
		}))

		return
	}

	if session.EndedAt.Valid || session.CancelledAt.Valid {
		c.Render(http.StatusBadRequest, api_errors.BAD_REQUEST.WithExtra(map[string]any{
			"err": "Session is already over",
		}))

		return
	}

	session.CancelledAt = models.JsonnableNullTime{Valid: true, Time: time.Now()}

	err = dal.SONGS.UpdateSession(session)
	if err != nil {
		c.Render(http.StatusInternalServerError, api_errors.DATABASE_ERROR.WithExtra(map[string]any{
			"err": err.Error(),
		}))

		return
	}

	title := ""
	nexusId := ""
	if session != nil && session.Song != nil {
		title = fmt.Sprintf("%v by %v", session.Song.Title, session.Song.Artist)
		nexusId = session.Song.NexusId
	} else {
		title = fmt.Sprintf("%v by %v", session.Title, session.Artist)
		nexusId = "Song unloaded"
	}

	ctxClaims, _ := c.Get("TokenClaims")
	claims := ctxClaims.(*models.JwtCustomClaims)
	log.Info(
		"Removing from queue",
		"username", claims.Username,
		"session", session.Id,
		"song", nexusId,
		"title", title,
	)

	if state.STATE.Karaoke.Current != nil && state.STATE.Karaoke.Current.Id == session.Id {
		services.KARAOKE.StartNextSong()
	} else {
		newSessions := []*models.SongSession{}
		for _, currSession := range state.STATE.KaraokeQueue {
			if currSession.Id != session.Id {
				newSessions = append(newSessions, currSession)
			}
		}

		state.STATE.KaraokeQueue = newSessions

		mercure_client.CLIENT.SendKaraokeQueue()
	}

	c.JSON(200, session)
}
