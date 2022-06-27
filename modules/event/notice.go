package event

import (
	"altsub/base"
	"altsub/models"
	"altsub/modules/notification"
	"altsub/modules/schema"
	"fmt"
)

func Notice(ev schema.SchemaedEvent, rcvs []models.MReceiver) {
	for _, rcv := range rcvs {
		nt := notification.NewNotification(rcv.Type)
		nt.SetEvent(ev)
		nt.ParseAuth(rcv.Auth)
		if err := nt.Notice(nt.RenderMsg()); err != nil {
			base.NewLog("error", err, fmt.Sprintf("告警发送失败 (rcv: %s)", string(rcv.Auth)), "event:Notice()")
			continue
		}
	}
}
