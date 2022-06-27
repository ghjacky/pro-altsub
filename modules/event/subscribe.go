package event

import (
	"altsub/models"
)

func checkSubscribe(rs []models.MRule) []models.MReceiver {
	var rcvs = []models.MReceiver{}
	for _, r := range rs {
		rcvs = append(rcvs, r.Receivers...)
	}
	return rcvs
}
