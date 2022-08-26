package models

import (
	"altsub/base"
	"errors"
	"time"

	"gorm.io/gorm"
)

type MMaintenance struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"created_at"`
	// DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index" `
	TX          *gorm.DB `json:"-" gorm:"-"`
	Name        string   `json:"name" gorm:"column:col_name;type:varchar(64);not null;uniqueIndex;comment:维护项名称，唯一"`
	Description string   `json:"description" gorm:"column:col_description;type:text;comment:维护项描述信息"`
	StartAt     int64    `json:"start_at" gorm:"column:col_start_at;not null;comment:维护项生效时段开始秒级时间戳"`
	EndAt       int64    `json:"end_at" gorm:"column:col_end_at;not null;comment:维护项生效时段结束秒级时间戳"`
	Rules       []MRule  `json:"rules" gorm:"many2many:tb_maintenances_rules"`
}

type MMaintenances struct {
	PQ  PageQuery
	TX  *gorm.DB
	All []*MMaintenance
}

func (*MMaintenance) TableName() string {
	return "tb_maintenances"
}

func (m *MMaintenance) Add() error {
	if m.TX == nil {
		err := errors.New("nil db object")
		base.NewLog("error", err, "新增维护项失败", "models:maintenance.Add()")
		return err
	}
	if len(m.Name) <= 0 {
		err := errors.New("empty maintenance name")
		base.NewLog("error", err, "新增维护项失败", "models;maintenance.Add()")
		return err
	}
	if err := m.TX.Create(m).Error; err != nil {
		base.NewLog("error", err, "新增维护项", "models:maintenance.Add()")
		return err
	}
	return nil
}
