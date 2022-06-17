package models

import (
	"altsub/base"
	"errors"

	"gorm.io/gorm"
)

type MEvent struct {
	BaseModel
	Eventid string  `json:"eventid" gorm:"column:col_eventid;not null;comment:类似于zabbix的eventid，用于区分同一条告警的故障、恢复"`
	Data    JSON `json:"data" gorm:"column:col_data;not null;comment:告警事件原始数据"`
}

type MEvents struct {
	PQ  PageQuery
	TX  *gorm.DB
	All []*MEvent
}

func (*MEvent) TableName() string {
	return "tb_events"
}

func (evs *MEvents) Add() error {
	if evs.TX == nil {
		err := errors.New("nil db object")
		base.NewLog("error", err, "告警事件写库失败", "models:events.Add()")
		return err
	}
	err := evs.TX.Create(evs).Error
	base.NewLog("", err, "告警事件写库", "models:events.Add()")
	return err
}
