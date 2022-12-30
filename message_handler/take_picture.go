package message_handler

import "github.com/partyhall/partyhall/socket"

type TakePictureHandler struct{}

func (h TakePictureHandler) GetType() string {
	return "TAKE_PICTURE"
}

func (h TakePictureHandler) Do(s *socket.Socket, payload interface{}) {
	s.TakePicture()
}
