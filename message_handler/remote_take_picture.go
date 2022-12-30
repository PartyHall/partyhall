package message_handler

import (
	"github.com/partyhall/partyhall/socket"
)

type RemoteTakePictureHandler struct{}

func (h RemoteTakePictureHandler) GetType() string {
	return "REMOTE_TAKE_PICTURE"
}

func (h RemoteTakePictureHandler) Do(s *socket.Socket, payload interface{}) {
	for _, sock := range socket.SOCKETS {
		sock.TakePicture()
	}
}
