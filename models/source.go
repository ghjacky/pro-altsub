package models

import (
	"altsub/base"
	"errors"
	"fmt"

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
		err := errors.New("nil db object")
		base.NewLog("error", err, "拉取告警源失败", "models:source.Add()")
		return err
	}
	err := ss.PQ.Query(ss.TX, &ss.All).Error
	base.NewLog("", err, "拉取告警源失败", "models:source.Add()")
	return err
}

func (s *MSource) Add() error {
	if s.TX == nil {
		err := errors.New("nil db object")
		base.NewLog("error", err, "新增告警源失败", "models:source.Add()")
		return err
	}
	if len(s.Name) <= 0 {
		err := errors.New("empty source name")
		base.NewLog("error", err, "新增告警源失败", "models:source.Add()")
		return err
	}
	err := s.TX.Create(s).Error
	base.NewLog("", err, fmt.Sprintf("新增告警源：%s", s.Name), "models:source.Add()")
	return err
}

func (s *MSource) GetByName() error {
	if s.TX == nil {
		err := errors.New("nil db object")
		base.NewLog("error", err, "根据名称获取告警源失败", "models:source.GetByName()")
		return err
	}
	if len(s.Name) <= 0 {
		err := errors.New("empty source name")
		base.NewLog("error", err, "根据名称获取告警源失败", "models:source.GetByName()")
		return err
	}
	err := s.TX.Where("col_name = ?", s.Name).First(s).Error
	base.NewLog("", err, fmt.Sprintf("根据名称（%s）获取告警源", s.Name), "models:source.GetByName()")
	return err
}
