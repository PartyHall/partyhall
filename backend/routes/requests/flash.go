package routes_requests

type SetFlashRequest struct {
	Powered    bool `json:"powered"`
	Brightness int  `json:"brightness"`
}
