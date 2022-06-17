package models

import (
	"altsub/base"
	"errors"

	"gorm.io/gorm"
)

type MEvent struct {
	BaseModel
	Data JSON `json:"data" gorm:"column:col_data;not null;comment:告警事件原始数据"`
}

type MEvents struct {
	PQ  PageQuery
	TX  *gorm.DB
	All []*MEvent
}

func (*MEvent) TableName() string {
	return "tb_events"
}

func (ev *MEvent) Add() error {
	if ev.TX == nil {
		err := errors.New("nil db object")
		base.NewLog("error", err, "告警事件写库失败", "models:events.Add()")
		return err
	}
	err := ev.TX.Create(ev).Error
	base.NewLog("", err, "告警事件写库", "models:events.Add()")
	return err
}
