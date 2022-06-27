package event

import (
	"altsub/models"
	"altsub/modules/rule"
	"altsub/modules/schema"
)

func checkRules(srcName string, ev schema.SchemaedEvent) []models.MRule {
	return rule.Check(srcName, ev)
}
