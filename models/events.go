package models

import (
	"altsub/base"
	"errors"
	"time"

	"gorm.io/gorm"
)

type MEvent struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"createdAt"`
	TX        *gorm.DB  `json:"-" gorm:"-"`
	Data      JSON      `json:"data" gorm:"column:col_data;not null;comment:告警事件原始数据"`
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
	if err := ev.TX.Create(ev).Error; err != nil {
		base.NewLog("error", err, "告警事件写库", "models:events.Add()")
		return err
	}
	return nil
}

type MSchemaedEvent struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"createdAt"`
	TX        *gorm.DB  `json:"-" gorm:"-"`
	Data      JSON      `json:"data" gorm:"column:col_data;not null;comment:解析过后的告警事件数据"` // SchemaedEvent
	Rules     []MRule   `json:"rules" gorm:"many2many:tb_events_rules"`
}

type MSchemaedEvents struct {
	PQ  PageQuery
	TX  *gorm.DB
	All []*MSchemaedEvent
}

func (*MSchemaedEvent) TableName() string {
	return "tb_schemaedevents"
}

func (se *MSchemaedEvent) Add() error {
	if se.TX == nil {
		err := errors.New("nil db object")
		base.NewLog("error", err, "新增解析后告警数据失败", "models:schemaedEvent.Add()")
		return err
	}
	if err := se.TX.Create(se).Error; err != nil {
		base.NewLog("error", err, "新增解析后告警数据", "models:schemaedEvent.Add()")
		return err
	}
	return nil
}

func (se *MSchemaedEvent) AppendRule(r MRule) error {
	if se.TX == nil {
		err := errors.New("nil db object")
		base.NewLog("error", err, "解析后告警数据新增关联规则失败", "models:schemaedEvent.AppendRule()")
		return err
	}
	if err := se.TX.Model(se).Association("Rules").Append(&r); err != nil {
		base.NewLog("error", err, "解析后告警数据新增关联规则", "models:schemaedEvent.AppendRule()")
		return err
	}
	return nil
}

func (ses *MSchemaedEvents) FetchAfter(t time.Time) error {
	if ses.TX == nil {
		err := errors.New("nil db object")
		base.NewLog("error", err, "获取指定时间段解析后告警数据失败", "models:schemaedEvent.AppendRule()")
		return err
	}
	tstring := t.Format("2006-01-02 15:04:05")
	if err := ses.TX.Where("createdAt > ?", tstring).Find(&ses.All).Error; err != nil {
		base.NewLog("error", err, "获取指定时间段解析后告警数据", "models:schemaedEvent.AppendRule()")
		return err
	}
	return nil
}
