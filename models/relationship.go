package models

import (
	"time"

	"gorm.io/gorm"
)

type MRelationship struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index" `
	TX        *gorm.DB       `json:"-" gorm:"-"`
	EventId   uint
}
