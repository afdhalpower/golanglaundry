package models

import (
	"time"

	"gorm.io/gorm"
)

type Service struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	Name           string         `gorm:"size:255;not null" json:"name" validate:"required,min=2,max=255"`
	PricePerKg     float64        `gorm:"type:decimal(12,2);not null" json:"price_per_kg" validate:"required,gt=0"`
	EstimatedHours int            `gorm:"not null;default:24" json:"estimated_hours" validate:"required,gt=0"`
	Description    string         `gorm:"type:text" json:"description"`
	IsActive       bool           `gorm:"default:true" json:"is_active"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (s *Service) TableName() string { return "services" }
