package models

// Too lazy to implement hydra in go
// So i'll do my thing for some metadata

type ResponseMetadata struct {
	LastPage int `json:"last_page"`
	Total    int `json:"total"`
}

type ContextualizedResponse struct {
	Results interface{}      `json:"results"`
	Meta    ResponseMetadata `json:"meta"`
}
