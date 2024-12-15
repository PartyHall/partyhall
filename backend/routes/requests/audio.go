package routes_requests

type AudioSetDevices struct {
	SourceId int `json:"source_id" binding:"required"`
	SinkId   int `json:"sink_id" binding:"required"`
}

type AudioSetDeviceVolume struct {
	Volume int `json:"volume"`
}
