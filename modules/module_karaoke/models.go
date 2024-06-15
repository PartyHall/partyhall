package module_karaoke

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/guregu/null/v5"
	"github.com/partyhall/partyhall/models"
)

type PhkSong struct {
	Id        int         `json:"-" db:"id"`
	Uuid      uuid.UUID   `json:"uuid" db:"uuid"`
	SpotifyID null.String `json:"spotify_id" db:"spotify_id"`
	Artist    string      `json:"artist" db:"artist"`
	Title     string      `json:"title" db:"title"`
	Hotspot   string      `json:"hotspot" db:"hotspot"`
	Format    string      `json:"format" db:"format"`

	HasCover  bool `json:"has_cover" db:"has_cover"`
	HasVocals bool `json:"has_vocals" db:"has_vocals"`
	HasFull   bool `json:"has_full" db:"has_full"`

	Filename string `json:"-" db:"filename"`
}

func (s PhkSong) String() string {
	return fmt.Sprintf("[%v] %v - %v (%v)", s.Uuid, s.Artist, s.Title, s.Format)
}

type SongSession struct {
	Id   int     `json:"id" db:"id"`
	Song PhkSong `json:"song" db:"song_id"`

	EventID int `json:"-" db:"event_id"`

	SungBy      string            `json:"sung_by" db:"sung_by"`
	AddedAt     *models.Timestamp `json:"added_at" db:"added_at"`
	StartedAt   *models.Timestamp `json:"started_at" db:"started_at"`
	EndedAt     *models.Timestamp `json:"ended_at" db:"ended_at"`
	CancelledAt *models.Timestamp `json:"cancelled_at" db:"cancelled_at"`
}

type SongImage struct {
	Id            int               `json:"id" db:"id"`
	SongSessionId int               `json:"song_session_id" db:"song_session_id"`
	CreatedAt     *models.Timestamp `json:"created_at" db:"created_at"`
}
