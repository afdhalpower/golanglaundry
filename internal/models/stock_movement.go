package models

import (
	"time"

	"gorm.io/gorm"
)

type StockMovement struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	InventoryID   uint           `gorm:"not null;index" json:"inventory_id"`
	Inventory     *Inventory     `gorm:"foreignKey:InventoryID" json:"inventory,omitempty"`
	Type          string         `gorm:"size:20;not null" json:"type"` // "in" or "out"
	Quantity      int            `gorm:"not null" json:"quantity"`
	PreviousStock int            `gorm:"not null" json:"previous_stock"`
	NewStock      int            `gorm:"not null" json:"new_stock"`
	Note          string         `gorm:"type:text" json:"note"`
	CreatedBy     uint           `gorm:"not null" json:"created_by"`
	User          *User          `gorm:"foreignKey:CreatedBy" json:"user,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (s *StockMovement) TableName() string {
	return "stock_movements"
}
