package models

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	OrderNumber string       `gorm:"size:50;uniqueIndex;not null" json:"order_number"`
	CustomerID uint          `gorm:"not null" json:"customer_id"`
	UserID     uint          `gorm:"not null" json:"user_id"`
	WeightKg   float64       `gorm:"type:decimal(10,2)" json:"weight_kg"`
	PricePerKg float64       `gorm:"type:decimal(12,2)" json:"price_per_kg"`
	Discount   float64       `gorm:"type:decimal(12,2);default:0" json:"discount"`
	ExtraCost  float64       `gorm:"type:decimal(12,2);default:0" json:"extra_cost"`
	Total      float64       `gorm:"type:decimal(12,2)" json:"total"`
	EntryDate  time.Time     `json:"entry_date"`
	Status     string        `gorm:"size:20;default:menunggu" json:"status"`
	Notes      string        `gorm:"type:text" json:"notes"`
	CreatedAt  time.Time     `json:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (o *Order) TableName() string { return "orders" }
