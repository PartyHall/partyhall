package routes_requests

type SetBackdropRequest struct {
	BackdropAlbum    *int64 `json:"backdrop_album"`
	SelectedBackdrop int    `json:"selected_backdrop"`
}
