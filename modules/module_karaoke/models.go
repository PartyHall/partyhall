package module_karaoke

import (
	"fmt"
	"time"

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

func (s SongSession) AsExportMetadata() map[string]any {
	metadata := map[string]any{
		"uuid":     s.Song.Uuid,
		"artist":   s.Song.Artist,
		"title":    s.Song.Title,
		"hotspot":  s.Song.Hotspot,
		"sung_by":  s.SungBy,
		"added_at": time.Time(*s.AddedAt).Format("2006-01-02 15:04:05"),
	}

	if s.StartedAt != nil {
		metadata["sung_at"] = time.Time(*s.StartedAt).Format("2006-01-02 15:04:05")
	}

	if s.CancelledAt != nil {
		metadata["cancelled_at"] = time.Time(*s.CancelledAt).Format("2006-01-02 15:04:05")
	}

	return metadata
}

type SongImage struct {
	Id            int               `json:"id" db:"id"`
	SongSessionId int               `json:"song_session_id" db:"song_session_id"`
	CreatedAt     *models.Timestamp `json:"created_at" db:"created_at"`
}
