package models

import "gorm.io/gorm"

type MSchema struct {
	BaseModel
	SourceID       uint    `json:"source_id" gorm:"column:col_source_id;not null;comment:schema相关联的 source id"`
	Data           JSON    `json:"data" gorm:"column:col_data;not null;comment:schema具体内容"`
	Source         MSource `json:"source" gorm:"foreignKey:SourceID;references:ID"`
}

type MSchemas struct {
	PQ  PageQuery
	TX  *gorm.DB
	All []*MSource
}

func (*MSchema) TableName() string {
	return "tb_schemas"
}
