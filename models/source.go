package models

import (
	"altsub/base"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type MSource struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	CreatedAt   time.Time      `json:"createdAt"`
	DeletedAt   gorm.DeletedAt `json:"deletedAt" gorm:"index" `
	TX          *gorm.DB       `json:"-" gorm:"-"`
	Name        string         `json:"name" gorm:"column:col_name;type:varchar(32);not null;uniqueIndex;comment:告警源名称"`
	Type        string         `json:"type" gorm:"column:col_type;type:varchar(32);not null;comment:告警源类型"`
	Description string         `json:"description" gorm:"column:col_description;type:text;null;default:NULL;comment:告警源描述"`
	Rules       []MRule        `json:"rules" gorm:"foreignKey:SourceID"`
}

type MSources struct {
	PQ  PageQuery
	TX  *gorm.DB
	All []*MSource
}

func (*MSource) TableName() string {
	return "tb_sources"
}

func (ss *MSources) Fetch(preloads ...string) error {
	if ss.TX == nil {
		err := errors.New("nil db object")
		base.NewLog("error", err, "拉取告警源失败", "models:source.Fetch()")
		return err
	}
	for _, p := range preloads {
		ss.TX = ss.TX.Preload(p)
	}
	if err := ss.PQ.Query(ss.TX, &ss.All).Error; err != nil {
		base.NewLog("error", err, "拉取告警源", "models:source.Fetch()")
		return err
	}
	return nil
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
	if err := s.TX.Create(s).Error; err != nil {
		base.NewLog("error", err, fmt.Sprintf("新增告警源：%s", s.Name), "models:source.Add()")
		return err
	}
	return nil
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
	if err := s.TX.Where("col_name = ?", s.Name).First(s).Error; err != nil {
		base.NewLog("error", err, fmt.Sprintf("根据名称（%s）获取告警源", s.Name), "models:source.GetByName()")
		return err
	}
	return nil
}

func (s *MSource) Get() error {
	if s.TX == nil {
		err := errors.New("nil db object")
		base.NewLog("error", err, "获取数据源失败", "models:source.Get()")
		return err
	}
	if err := s.TX.First(s).Error; err != nil {
		base.NewLog("error", err, "获取数据源失败", "models:source.Get()")
		return err
	}
	return nil
}

func (ss *MSources) FetchTypes() error {
	if ss.TX == nil {
		err := errors.New("nil db object")
		base.NewLog("error", err, "获取数据源类型失败", "models:source.FetchTypes()")
		return err
	}
	if err := ss.TX.Distinct("col_type").Select("col_type").Find(&ss.All).Error; err != nil {
		base.NewLog("error", err, "获取数据源类型失败", "models:source.FetchTypes()")
		return err
	}
	return nil
}
