package models

import (
	"time"
)

type Payment struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	OrderID     uint      `gorm:"uniqueIndex;not null" json:"order_id"`
	Amount      float64   `gorm:"type:decimal(12,2);not null" json:"amount"`
	Method      string    `gorm:"size:20;not null;default:tunai" json:"method"`
	Status      string    `gorm:"size:20;not null;default:belum_lunas" json:"status"`
	PaymentDate time.Time `json:"payment_date"`
	Note        string    `gorm:"type:text" json:"note"`
	CreatedBy   uint      `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (p *Payment) TableName() string { return "payments" }
