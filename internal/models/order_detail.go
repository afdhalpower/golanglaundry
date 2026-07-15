package models

type OrderDetail struct {
	ID        uint     `gorm:"primaryKey" json:"id"`
	OrderID   uint     `gorm:"not null;index" json:"order_id"`
	ServiceID uint     `gorm:"not null" json:"service_id"`
	Service   *Service `gorm:"foreignKey:ServiceID" json:"-"`
	PricePerKg float64 `gorm:"type:decimal(12,2)" json:"price_per_kg"`
}

func (od *OrderDetail) TableName() string { return "order_details" }
