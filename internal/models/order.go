package models

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	OrderNumber      string         `gorm:"size:50;uniqueIndex;not null" json:"order_number"`
	CustomerID       uint           `gorm:"not null;index" json:"customer_id"`
	Customer         *Customer      `gorm:"foreignKey:CustomerID" json:"-"`
	UserID           uint           `gorm:"not null;index" json:"user_id"`
	User             *User          `gorm:"foreignKey:UserID" json:"-"`
	WeightKg         float64        `gorm:"type:decimal(10,2)" json:"weight_kg"`
	Discount         float64        `gorm:"type:decimal(12,2);default:0" json:"discount"`
	ExtraCost        float64        `gorm:"type:decimal(12,2);default:0" json:"extra_cost"`
	Total            float64        `gorm:"type:decimal(12,2)" json:"total"`
	EntryDate        time.Time      `json:"entry_date"`
	EstimatedDoneDate time.Time     `json:"estimated_done_date"`
	Status           string         `gorm:"size:20;default:menunggu;index" json:"status"`
	Notes            string         `gorm:"type:text" json:"notes"`
	Details          []OrderDetail  `gorm:"foreignKey:OrderID" json:"-"`
	Tracking         []OrderTracking `gorm:"foreignKey:OrderID;order:created_at ASC" json:"-"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

func (o *Order) TableName() string { return "orders" }
