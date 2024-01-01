package module_photobooth

import (
	"github.com/partyhall/easyws"
	"github.com/partyhall/partyhall/remote"
)

type TakePictureHandler struct{}

func (h TakePictureHandler) GetType() string {
	return "photobooth/TAKE_PICTURE"
}

func (h TakePictureHandler) Do(s *easyws.Socket, payload interface{}) {
	INSTANCE.Actions.TakePicture(s)
}

type RemoteTakePictureHandler struct{}

func (h RemoteTakePictureHandler) GetType() string {
	return "photobooth/REMOTE_TAKE_PICTURE"
}

func (h RemoteTakePictureHandler) Do(s *easyws.Socket, payload interface{}) {
	for _, s := range remote.EasyWS.Sockets {
		INSTANCE.Actions.TakePicture(s)
	}
}
