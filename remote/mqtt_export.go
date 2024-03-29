package remote

import (
	"strconv"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/partyhall/easyws"
	"github.com/partyhall/partyhall/logs"
	"github.com/partyhall/partyhall/orm"
	"github.com/partyhall/partyhall/services"
)

func OnExportZip(client mqtt.Client, msg mqtt.Message) {
	eventIdStr := string(msg.Payload())
	eventId, err := strconv.Atoi(eventIdStr)
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
	BroadcastAdmin("EXPORT_STARTED", event)

	go func() {
		exportedEvent, err := (services.NewEventExporter(event)).Export()
		if err != nil {
			logs.Error(err)
			BroadcastAdmin(easyws.MSG_TYPE_ERR_MODAL, err.Error())
			return
		}

		BroadcastAdmin("EXPORT_COMPLETED", exportedEvent)
	}()
}
