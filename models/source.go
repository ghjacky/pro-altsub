package models

import (
	"errors"

	"gorm.io/gorm"
)

type MSource struct {
	BaseModel
	Name        string `json:"name" gorm:"column:col_name;type:varchar(32);not null;uniqueIndex;comment:告警源名称"`
	Description string `json:"description" gorm:"column:col_description;type:text;null;default:NULL;comment:告警源描述"`
}

type MSources struct {
	PQ  PageQuery
	TX  *gorm.DB
	All []*MSource
}

func (*MSource) TableName() string {
	return "tb_sources"
}

func (ss *MSources) Fetch() error {
	if ss.TX == nil {
		return errors.New("nil db object")
	}
	return ss.PQ.Query(ss.TX, &ss.All).Error
}

func (s *MSource) Add() error {
	if s.DB == nil {
		return errors.New("nil db object")
	}
	if len(s.Name) <= 0 {
		return errors.New("empty source name")
	}
	return s.DB.Create(s).Error
}
