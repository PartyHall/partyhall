package module_karaoke

type DtoSongCreate struct {
	Title     string `form:"title" validate:"required"`
	Artist    string `form:"artist" validate:"required"`
	Format    string `form:"format" validate:"required,oneof=CDG WEBM MP4"`
	CoverType string `form:"cover_type" validate:"required,oneof=LINK UPLOADED NO_COVER"`
	// Song      *multipart.FileHeader `form:"song" validate:"required"`
	// CDG       *multipart.FileHeader `form:"cdg"`

	// Cover    *multipart.FileHeader `form:"cover"`
	CoverUrl *string `form:"cover_url"`
}

type DtoSongGet struct {
	Id       int    `json:"id"`
	Title    string `json:"title"`
	Artist   string `json:"artist"`
	Format   string `json:"format"`
	Filename string `json:"filename"`
}

func songGet(dbSong *PhkSong) DtoSongGet {
	return DtoSongGet{
		Id:     dbSong.Id,
		Title:  dbSong.Title,
		Artist: dbSong.Artist,
		Format: dbSong.Format,
	}
}
