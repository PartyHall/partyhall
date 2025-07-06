package models

import (
	"time"
)

type Event struct {
	Id       int64     `db:"id" json:"id"`
	Name     string    `db:"name" json:"name"`
	Author   string    `db:"author" json:"author"`
	Date     time.Time `db:"date" json:"date"`
	Location string    `db:"location" json:"location"`

	NexusId             JsonnableNullstring `db:"nexus_id" json:"nexus_id"`
	UserRegistrationUrl JsonnableNullstring `db:"registration_url" json:"registration_url"`

	AmtImageHandtaken  int `db:"amt_images_handtaken" json:"amt_images_handtaken"`
	AmtImageUnattended int `db:"amt_images_unattended" json:"amt_images_unattended"`
}

func (e Event) AsJson() map[string]any {
	return map[string]any{
		"id":                    e.Id,
		"name":                  e.Name,
		"author":                e.Author,
		"date":                  e.Date.Format(time.RFC3339),
		"location":              e.Location,
		"amt_images_handtaken":  e.AmtImageHandtaken,
		"amt_images_unattended": e.AmtImageUnattended,
	}
}

type Picture struct {
	Id         string              `db:"id" json:"id"`
	TakenAt    time.Time           `db:"taken_at" json:"taken_at"`
	Unattended bool                `db:"unattended" json:"unattended"`
	Filename   string              `db:"filename" json:"-"`
	EventId    int64               `db:"event_id" json:"-"`
	NexusId    JsonnableNullstring `db:"nexus_id" json:"nexus_id"`
}

func (p Picture) AsJson() map[string]any {
	return map[string]any{
		"id":       p.Id,
		"filename": p.Filename,
		"taken_at": p.TakenAt.Format(time.RFC3339),
	}
}
