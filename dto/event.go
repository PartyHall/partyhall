package dto

import "github.com/partyhall/partyhall/models"

type EventPost struct {
	Name     string           `json:"name" validate:"required"`
	Date     models.Timestamp `json:"date" validate:"required"`
	Author   string           `json:"author" validate:"required"`
	Location *string          `json:"location"`
}

type EventPut struct {
	Name     string           `json:"name,omitempty"`
	Date     models.Timestamp `json:"date,omitempty"`
	Author   string           `json:"author,omitempty"`
	Location *string          `json:"location,omitempty"`
}

type EventGet struct {
	Id       int              `json:"id"`
	Name     string           `json:"name"`
	Date     models.Timestamp `json:"date"`
	Author   string           `json:"author"`
	Location *string          `json:"location"`
}

func GetEvent(dbEvent *models.Event) EventGet {
	return EventGet{
		Id:       dbEvent.Id,
		Name:     dbEvent.Name,
		Date:     dbEvent.Date,
		Author:   dbEvent.Author,
		Location: dbEvent.Location,
	}
}
