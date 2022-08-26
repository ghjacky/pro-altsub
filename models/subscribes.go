package models

import (
	"altsub/base"
	"errors"
	"time"

	"gorm.io/gorm"
)

const (
	SubscribeTypeSub = iota + 1
	SubscribeTypeAss
)

type MSubscribe struct {
	CreatedAt   time.Time `json:"created_at" form:"-"`
	Name        string    `json:"name" form:"name" gorm:"column:col_name;comment:订阅或指派名称"`
	Description string    `json:"description" form:"description" gorm:"column:col_description;type:text;comment:订阅或指派描述"`
	StartAt     int64     `json:"start_at" form:"start_at" gorm:"column:col_start_at;not null;comment:订阅或指派生效时段开始秒级时间戳"`
	EndAt       int64     `json:"end_at" form:"end_at" gorm:"column:col_end_at;not null;comment:订阅或指派生效时段结束秒级时间戳"`
	ReceiverID  uint      `json:"receiver_id" form:"-" gorm:"column:col_receiver_id;primaryKey"`
	RuleID      uint      `json:"rule_id" form:"-" gorm:"column:col_rule_id;primaryKey"`
	Type        int       `json:"type" form:"type" gorm:"column:col_type;comment:订阅、指派"`
	Receiver    MReceiver `json:"receiver" form:"-" gorm:"-"`
	Rule        MRule     `json:"rule" form:"-" gorm:"-"`
	Source      MSource   `json:"source" form:"-" gorm:"-"`
}

type MSubscribes struct {
	PQ  PageQuery
	TX  *gorm.DB
	All []*MSubscribe
}

func (s *MSubscribe) BeforeSave(tx *gorm.DB) error {
	if tx.Statement.Context.Value("subscribe") == nil {
		return nil
	}
	sub, ok := tx.Statement.Context.Value("subscribe").(MSubscribe)
	if !ok {
		err := errors.New("wrong struct value in tx.Context")
		base.NewLog("error", err, "新增订阅或指派", "models;subscribe.BeforeCreate()")
		return err
	}
	now := time.Now().Local()
	s.CreatedAt = now
	s.Name = sub.Name
	s.Description = sub.Description
	s.StartAt = sub.StartAt
	s.EndAt = sub.EndAt
	s.Type = sub.Type
	return nil
}

func (*MSubscribe) TableName() string {
	return "tb_subscribes"
}

func (ss *MSubscribes) Fetch() error {
	if ss.TX == nil {
		err := errors.New("nil db objec")
		base.NewLog("error", err, "获取订阅/指派关系失败", "models:subscribes.fetch()")
		return err
	}
	if err := ss.PQ.Query(ss.TX, &ss.All, &MSubscribe{}).Error; err != nil {
		base.NewLog("error", err, "获取订阅/指派关系失败", "models:subscribes.fetch()")
		return err
	}
	for _, s := range ss.All {
		s.Receiver.ID = s.ReceiverID
		s.Receiver.TX = ss.TX
		s.Receiver.Get()
		s.Rule.ID = s.RuleID
		s.Rule.TX = ss.TX
		s.Rule.Get()
		s.Source.ID = s.Rule.SourceID
		s.Source.TX = ss.TX
		s.Source.Get()
	}
	return nil
}

func (s *MSubscribe) DeleteByReceiver(tx *gorm.DB) error {
	return tx.Where("col_receiver_id = ?", s.ReceiverID).Unscoped().Delete(&MSubscribe{}).Error
}

func (s *MSubscribe) DeleteByRule(tx *gorm.DB) error {
	return tx.Where("col_rule_id = ?", s.RuleID).Unscoped().Delete(&MSubscribe{}).Error
}
