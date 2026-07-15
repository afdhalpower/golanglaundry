package models

import "time"

type ExpenseCategory struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Name      string `gorm:"size:100;not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (e *ExpenseCategory) TableName() string { return "expense_categories" }
