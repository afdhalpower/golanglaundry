package models

import "time"

type OrderTracking struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	OrderID   uint      `gorm:"not null;index" json:"order_id"`
	Status    string    `gorm:"size:20;not null" json:"status"`
	Note      string    `gorm:"type:text" json:"note"`
	CreatedBy uint      `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}

func (ot *OrderTracking) TableName() string { return "order_tracking" }
