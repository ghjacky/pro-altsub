package event

import (
	"altsub/models"
	"altsub/modules/maintenance"
)

func checkMaintenance(rs []models.MRule) bool {
	return maintenance.Check(rs)
}
