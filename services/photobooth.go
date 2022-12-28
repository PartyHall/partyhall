package services

import (
	"strconv"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/partyhall/partyhall/logs"
	"github.com/partyhall/partyhall/models"
	"github.com/partyhall/partyhall/orm"
)

type PartyHall struct {
	prv *Provider

	CurrentState    models.AppState
	IsTakingPicture bool
	CurrentMode     string
	DisplayDebug    bool
}

func (pb *PartyHall) OnSyncRequested(client mqtt.Client, msg mqtt.Message) {
	logs.Info("Sync requested")
}

func (pb *PartyHall) OnExportEvent(client mqtt.Client, msg mqtt.Message) {
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
	pb.prv.Sockets.BroadcastAdmin("EXPORT_STARTED", event)

	go func() {
		exportedEvent, err := (NewEventExporter(event)).Export()
		if err != nil {
			logs.Error(err)
			pb.prv.Sockets.BroadcastAdmin("ERR_MODAL", err.Error())
			return
		}

		pb.prv.Sockets.BroadcastAdmin("EXPORT_COMPLETED", exportedEvent)
	}()
}
