package models

import (
	"fmt"
)

type Song struct {
	NexusId   string              `db:"nexus_id" json:"nexus_id"`
	Title     string              `db:"title" json:"title"`
	Artist    string              `db:"artist" json:"artist"`
	Format    string              `db:"format" json:"format"`
	Duration  int64               `db:"duration" json:"duration"`
	SpotifyId JsonnableNullstring `db:"spotify_id" json:"spotify_id"`
	Hotspot   JsonnableNullInt64  `db:"hotspot" json:"hotspot"`

	HasCover    bool `db:"has_cover" json:"has_cover"`
	HasVocals   bool `db:"has_vocals" json:"has_vocals"`
	HasCombined bool `db:"has_combined" json:"has_combined"`
}

type SongSession struct {
	Id             int64              `db:"id" json:"id"`
	EventId        int64              `db:"event_id" json:"event_id"`
	SessionNexusId JsonnableNullInt64 `db:"session_nexus_id" json:"session_nexus_id"`
	NexusId        string             `db:"nexus_id" json:"nexus_id"`
	Title          string             `db:"title" json:"title"`
	Artist         string             `db:"artist" json:"artist"`
	SungBy         string             `db:"sung_by" json:"sung_by"`
	SungById       string             `db:"-" json:"sung_by_id"`
	Song           *Song              `db:"-" json:"song"`
	AddedAt        JsonnableNullTime  `db:"added_at" json:"added_at"`
	StartedAt      JsonnableNullTime  `db:"started_at" json:"started_at"`
	EndedAt        JsonnableNullTime  `db:"ended_at" json:"ended_at"`
	CancelledAt    JsonnableNullTime  `db:"cancelled_at" json:"cancelled_at"`
}

func (session SongSession) String() string {
	return fmt.Sprintf(
		"SongSession[id=%v, event_id=%v, nexus_id=%v, title=%v, artist=%v, sung_by=%v, sung_by_id=%v]",
		session.Id,
		session.EventId,
		session.NexusId,
		session.Title,
		session.Artist,
		session.SungBy,
		session.SungById,
	)
}
