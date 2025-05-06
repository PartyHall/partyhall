package routes

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/partyhall/partyhall/api_errors"
	"github.com/partyhall/partyhall/log"
	"github.com/partyhall/partyhall/mercure_client"
	"github.com/partyhall/partyhall/middlewares"
	"github.com/partyhall/partyhall/models"
	"github.com/partyhall/partyhall/state"
	"github.com/partyhall/partyhall/utils"
)

/**
 * Those routes are separated instead of on the same endpoint
 * so that the timecode that FLOODS the request does not
 * overlap and cause race errors with the admin requests
 **/

type RoutesKaraoke struct{}

func (h RoutesKaraoke) Register(router *gin.RouterGroup) {
	// hasEvent & Appliance
	router.PUT(
		"timecode",
		middlewares.HasEventLoaded(),
		middlewares.Authorized(models.ROLE_APPLIANCE),
		h.setTimecode,
	)

	router.PUT(
		"playing_status",
		middlewares.HasEventLoaded(),
		middlewares.Authorized(),
		h.setPlayingStatus,
	)

	router.PUT(
		"volume",
		middlewares.HasEventLoaded(),
		middlewares.Authorized(),
		h.setVolume,
	)
}

func (h RoutesKaraoke) setTimecode(c *gin.Context) {
	var req struct {
		Timecode int `json:"timecode"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		api_errors.RenderValidationErr(c, err)
		return
	}

	if state.STATE.Karaoke.Timecode != req.Timecode {
		state.STATE.Karaoke.Timecode = req.Timecode
		mercure_client.CLIENT.SendKaraokeTimecode()
	}

	c.Status(200)
}

func (h RoutesKaraoke) setPlayingStatus(c *gin.Context) {
	var req struct {
		Status bool `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		api_errors.RenderValidationErr(c, err)
		return
	}

	state.STATE.Karaoke.IsPlaying = req.Status
	tx := "paused"
	if req.Status {
		tx = "resumed"
	}

	ctxClaims, _ := c.Get("TokenClaims")
	claims := ctxClaims.(*models.JwtCustomClaims)
	log.Info(
		"Song "+tx,
		"username", claims.Username,
		"session", state.STATE.Karaoke.Current.Id,
		"song", state.STATE.Karaoke.Current.Song.NexusId,
		"title", fmt.Sprintf("%v by %v", state.STATE.Karaoke.Current.Song.Title, state.STATE.Karaoke.Current.Song.Artist),
	)

	mercure_client.CLIENT.SendKaraokeState()
}

func (h RoutesKaraoke) setVolume(c *gin.Context) {
	var req struct {
		Type   string `json:"type" binding:"required"`
		Volume int    `json:"volume"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		api_errors.RenderValidationErr(c, err)
		return
	}

	volumeType := strings.ToLower(req.Type)

	if volumeType != "instrumental" && volumeType != "vocals" && volumeType != "combined" {
		c.Render(http.StatusBadRequest, api_errors.BAD_REQUEST.WithExtra(map[string]any{
			"err": "Volume type should be either instrumental, vocals or combined",
		}))

		return
	}

	vol := utils.ClampInt(req.Volume, 0, 100)

	if volumeType == "instrumental" {
		state.STATE.Karaoke.Volume = vol
	} else if volumeType == "vocals" {
		state.STATE.Karaoke.VolumeVocals = vol
	}

	mercure_client.CLIENT.SendKaraokeState()

	c.JSON(200, state.STATE.Karaoke)
}
