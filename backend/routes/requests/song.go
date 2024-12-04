package routes_requests

type SongEnqueue struct {
	SongId      string `json:"song_id" binding:"required"`
	DisplayName string `json:"display_name" binding:"required"`
	DirectPlay  bool   `json:"direct_play"`
}

type SetTimecode struct {
	Timecode int `json:"timecode"`
}
