package models

import (
	"time"
)

type Expense struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	ExpenseCategoryID uint      `json:"expense_category_id"`
	Amount            float64   `gorm:"type:decimal(12,2);not null" json:"amount"`
	Description       string    `gorm:"type:text" json:"description"`
	Date              time.Time `json:"date"`
	CreatedBy         uint      `json:"created_by"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func (e *Expense) TableName() string { return "expenses" }
