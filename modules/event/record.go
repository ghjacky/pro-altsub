package event

import (
	"altsub/base"
	"altsub/models"
	"altsub/modules/schema"
	"encoding/json"
)

// 原始事件存储
func StoreRawToDb(ev *models.MEvent) error {
	ev.TX = base.DB()
	if err := ev.Add(); err != nil {
		base.NewLog("error", err, "原始事件入库", "event:StoreRawToDb()")
		return err
	}
	return nil
}

// 解析后数据入库
func StoreParsedEvToDb(parsedEv schema.SchemaedEvent, rs []models.MRule) error {
	var ev = models.MSchemaedEvent{TX: base.DB().Begin()}
	ev.Data, _ = json.Marshal(parsedEv)
	ev.Rules = rs
	if err := ev.Add(); err != nil {
		base.NewLog("error", err, "解析后数据入库失败", "event:StoreParsedEvToDb()")
		ev.TX.Rollback()
		return err
	}
	ev.TX.Commit()
	return nil
}

// 记录事件发送（成功、失败）（事件id、发送时间、接受者id、成功与否）
