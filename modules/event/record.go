package event

import (
	"altsub/base"
	"altsub/models"
	"errors"
)

// 原始事件存储
func StoreToDb(ev *models.MEvent) error {
	if ev.TX == nil {
		err := errors.New("nil db object")
		base.NewLog("error", err, "事件入库失败", "event:StoreToDb()")
		return err
	}
	err := ev.Add()
	base.NewLog("", err, "事件入库", "event:StoreToDb()")
	return err
}

// 记录事件发送（成功、失败）（事件id、发送时间、接受者id、成功与否）
