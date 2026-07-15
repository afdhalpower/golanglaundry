package models

import (
	"time"

	"gorm.io/gorm"
)

type Expense struct {
	ID                uint            `gorm:"primaryKey" json:"id"`
	ExpenseCategoryID uint            `gorm:"not null;index" json:"expense_category_id"`
	ExpenseCategory   *ExpenseCategory `gorm:"foreignKey:ExpenseCategoryID" json:"-"`
	Amount            float64         `gorm:"type:decimal(12,2);not null" json:"amount"`
	Description       string          `gorm:"type:text" json:"description"`
	Date              time.Time       `json:"date"`
	CreatedBy         uint            `json:"created_by"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
	DeletedAt         gorm.DeletedAt  `gorm:"index" json:"-"`
}

func (e *Expense) TableName() string { return "expenses" }
