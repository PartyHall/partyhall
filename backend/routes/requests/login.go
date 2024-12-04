package routes_requests

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginGuestRequest struct {
	Username string `json:"username"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}
