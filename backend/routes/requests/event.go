package routes_requests

import "github.com/partyhall/partyhall/models"

type CreateEvent struct {
	Id       int                        `json:"id"`
	Name     string                     `json:"name" binding:"required"`
	Author   string                     `json:"author"`
	Date     string                     `json:"date" binding:"iso8601,required"`
	Location string                     `json:"location"`
	NexusID  models.JsonnableNullstring `json:"nexus_id"`
}
