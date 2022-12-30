package socket

type Sockets []*Socket

var SOCKETS Sockets = []*Socket{}

func (s Sockets) BroadcastTo(to string, msgType string, data interface{}) {
	for _, socket := range s {
		if len(to) > 0 && socket.Type != to {
			continue
		}

		socket.Send(msgType, data)
	}
}

func (s Sockets) BroadcastBooth(msgType string, data interface{}) {
	s.BroadcastTo(SOCKET_TYPE_BOOTH, msgType, data)
}

func (s Sockets) BroadcastTakePicture() {
	for _, socket := range s {
		if socket.Type != SOCKET_TYPE_BOOTH {
			continue
		}

		socket.TakePicture()
	}
}

func (s Sockets) BroadcastAdmin(msgType string, data interface{}) {
	s.BroadcastTo(SOCKET_TYPE_ADMIN, msgType, data)
}

func (s Sockets) BroadcastState() {
	for _, socket := range s {
		socket.SendState()
	}
}
