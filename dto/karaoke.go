package dto

import "github.com/partyhall/partyhall/models"

type SongDto struct {
	*models.Song
	SungBy string `json:"sung_by"`
}
