package nexus

import "fmt"

type ApiSong struct {
	Iri           string  `json:"@id"`
	Id            int64   `json:"id"`
	Title         string  `json:"title"`
	Artist        string  `json:"artist"`
	Format        string  `json:"format"`
	Quality       string  `json:"quality"`
	SpotifyID     *string `json:"spotifyId"`
	MusicbrainzID *string `json:"musicBrainzId"`
	NexusBuildId  *string `json:"nexusBuildId"`
	Ready         bool    `json:"ready"`
	Cover         bool    `json:"cover"`
	CoverUrl      string  `json:"coverUrl"`
	Vocals        bool    `json:"vocals"`
	Combined      bool    `json:"combined"`
}

func (song ApiSong) String() string {
	return fmt.Sprintf("%s - %s", song.Title, song.Artist)
}

type PhkSong struct {
	//#region @TODO: Remove those from the generated json file
	Context   string `json:"@context"`
	Iri       string `json:"@id"`
	HydraType string `json:"@type"`
	//#endregion

	NexusId       string `json:"nexusBuildId"`
	Title         string `json:"title"`
	Artist        string `json:"artist"`
	Format        string `json:"format"`
	Quality       string `json:"quality"`
	MusicBrainzId string `json:"musicBrainzId"`
	SpotifyId     string `json:"spotifyId"`
	Duration      int64  `json:"duration"`
	Hotspot       *int64 `json:"hotspot"`
}

type ApiBackdropAlbum struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Author  string `json:"author"`
	Version int64  `json:"version"`
}

type ApiBackdrop struct {
	Id    int64  `json:"id"`
	Title string `json:"title"`
	Url   string `json:"url"`
}
