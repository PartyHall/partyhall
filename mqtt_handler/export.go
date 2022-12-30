package mqtt_handler

import (
	"strconv"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/partyhall/partyhall/logs"
	"github.com/partyhall/partyhall/orm"
	"github.com/partyhall/partyhall/services"
	"github.com/partyhall/partyhall/socket"
)

type ExportHandler struct{}

func (h ExportHandler) GetTopic() string {
	return "export"
}

func (h ExportHandler) Do(client mqtt.Client, msg mqtt.Message) {
	eventIdStr := string(msg.Payload())
	eventId, err := strconv.ParseInt(eventIdStr, 10, 64)
	if err != nil {
		logs.Error("Failed to export event: bad eventid => ", eventIdStr)
		logs.Error(err)

		return
	}

	event, err := orm.GET.Events.GetEvent(eventId)
	if err != nil {
		logs.Error("Failed to export event:", err)
		return
	}

	logs.Info("Export requested")
	socket.SOCKETS.BroadcastAdmin("EXPORT_STARTED", event)

	go func() {
		exportedEvent, err := (services.NewEventExporter(event)).Export()
		if err != nil {
			logs.Error(err)
			socket.SOCKETS.BroadcastAdmin("ERR_MODAL", err.Error())
			return
		}

		socket.SOCKETS.BroadcastAdmin("EXPORT_COMPLETED", exportedEvent)
	}()
}
