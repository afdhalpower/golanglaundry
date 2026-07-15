package models

import (
	"time"

	"gorm.io/gorm"
)

type Inventory struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Category    string         `gorm:"size:50;not null" json:"category"`
	Stock       int            `gorm:"not null;default:0" json:"stock"`
	MinStock    int            `gorm:"not null;default:0" json:"min_stock"`
	Unit        string         `gorm:"size:20;not null;default:pcs" json:"unit"`
	Description string         `gorm:"type:text" json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (i *Inventory) TableName() string { return "inventories" }
