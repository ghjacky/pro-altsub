package models

import (
	"altsub/base"
	"errors"
	"time"

	"gorm.io/gorm"
)

type MSchema struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index" `
	TX        *gorm.DB       `json:"-" gorm:"-"`
	Data      JSON           `json:"data" gorm:"column:col_data;not null;comment:schema具体内容"`
	EvField   string         `json:"ev_field" gorm:"column:col_ev_field;not null;default:.;comment:event数据从哪个字段中获取，默认：'.'（代表上报上来的整个数据即为event数据本身）"`
	EvType    string         `json:"ev_type" gorm:"column:col_ev_type;not null;default:map;comment:指定获取event数据字段的类型，一般为map或者array"`
	SourceID  uint           `json:"source_id" gorm:"column:col_source_id;not null;uniqueIndex;comment:schema相关联的 source id"`
	Source    MSource        `json:"source" gorm:"foreignKey:SourceID;references:ID"`
}

type MSchemas struct {
	PQ  PageQuery
	TX  *gorm.DB
	All []*MSchema
}

func (*MSchema) TableName() string {
	return "tb_schemas"
}

func (s *MSchema) Get() error {
	if s.TX == nil {
		err := errors.New("nil db object")
		base.NewLog("error", err, "获取schema失败", "models:schema.Get()")
		return err
	}
	if err := s.TX.First(s).Error; err != nil {
		base.NewLog("error", err, "获取schema", "models:schema.Get()")
		return err
	}
	return nil
}

func (ss *MSchemas) Fetch() error {
	if ss.TX == nil {
		err := errors.New("nil db object")
		base.NewLog("error", err, "拉取schema失败", "models:schema.Add()")
		return err
	}
	if err := ss.PQ.Query(ss.TX, &ss.All, &MSchema{}).Error; err != nil {
		base.NewLog("error", err, "拉取schema", "models:schema.Add()")
		return err
	}
	return nil
}

func (s *MSchema) Add() error {
	if s.TX == nil {
		err := errors.New("nil db object")
		base.NewLog("error", err, "新增schema失败", "models:schema.Add()")
		return err
	}
	if len(s.Data) <= 0 {
		err := errors.New("empty schema data")
		base.NewLog("error", err, "新增schema失败", "models:schema.Add()")
		return err
	}
	if s.Source.ID == 0 {
		if len(s.Source.Name) == 0 {
			err := errors.New("wrong ralated source")
			base.NewLog("error", err, "新增schema失败", "models:schema.Add()")
			return err
		} else {
			if err := s.Source.GetByName(); err != nil {
				base.NewLog("error", err, "新增schema失败", "models:schema.Add()")
				return err
			}
		}
	}
	if err := s.TX.Create(s).Error; err != nil {
		base.NewLog("error", err, "新增schema", "models:schema.Add()")
		return err
	}
	return nil
}

func (s *MSchema) Update() error {
	if s.TX == nil {
		err := errors.New("nil db object")
		base.NewLog("error", err, "更新schema失败", "models:schema.Update()")
		return err
	}
	if len(s.Data) <= 0 {
		err := errors.New("empty schema data")
		base.NewLog("error", err, "更新schema失败", "models:schema.Update()")
		return err
	}
	if s.Source.ID == 0 {
		if len(s.Source.Name) == 0 {
			err := errors.New("wrong ralated source")
			base.NewLog("error", err, "更新schema失败", "models:schema.Update()")
			return err
		} else {
			if err := s.Source.GetByName(); err != nil {
				base.NewLog("error", err, "更新schema失败", "models:schema.Update()")
				return err
			}
		}
	}
	if err := s.TX.Save(s).Error; err != nil {
		base.NewLog("error", err, "更新schema", "models:schema.Update()")
		return err
	}
	return nil
}

func (s *MSchema) GetBySourceID() error {
	if s.TX == nil {
		err := errors.New("nil db object")
		base.NewLog("error", err, "根据source_id获取schema失败", "models:schema.GetBySourceID()")
		return err
	}
	if s.SourceID == 0 {
		err := errors.New("zero source id")
		base.NewLog("error", err, "根据source_id获取schema失败", "models:schema.GetBySourceID()")
		return err
	}
	if err := s.TX.Where("col_source_id = ?", s.SourceID).First(s).Error; err != nil {
		base.NewLog("error", err, "根据source_id获取schema", "models:schema.GetBySourceID()")
		return err
	}
	return nil
}
