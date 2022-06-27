package models

import (
	"time"

	"gorm.io/gorm"
)

type MRelationship struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"createdAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index" `
	TX        *gorm.DB       `json:"-" gorm:"-"`
	EventId   uint
}
