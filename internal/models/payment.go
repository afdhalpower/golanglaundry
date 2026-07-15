package models

import (
	"time"

	"gorm.io/gorm"
)

type Payment struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	OrderID     uint           `gorm:"uniqueIndex;not null" json:"order_id"`
	Order       *Order         `gorm:"foreignKey:OrderID" json:"-"`
	Amount      float64        `gorm:"type:decimal(12,2);not null" json:"amount"`
	Method      string         `gorm:"size:20;not null;default:tunai" json:"method"`
	Status      string         `gorm:"size:20;not null;default:lunas" json:"status"`
	PaymentDate time.Time      `json:"payment_date"`
	Note        string         `gorm:"type:text" json:"note"`
	CreatedBy   uint           `json:"created_by"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (p *Payment) TableName() string { return "payments" }
