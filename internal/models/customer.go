package models

import (
	"time"

	"gorm.io/gorm"
)

type Customer struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:255;not null" json:"name" validate:"required,min=2,max=255"`
	Phone     string         `gorm:"size:20" json:"phone"`
	Whatsapp  string         `gorm:"size:20" json:"whatsapp"`
	Address   string         `gorm:"type:text" json:"address"`
	Notes     string         `gorm:"type:text" json:"notes"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (c *Customer) TableName() string { return "customers" }
