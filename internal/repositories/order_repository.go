package repositories

import (
	"fmt"
	"time"

	"github.com/afdhalpower/golanglaundry/internal/models"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) GenerateOrderNumber(date time.Time) (string, error) {
	dateStr := date.Format("20060102")
	prefix := fmt.Sprintf("INV/%s/", dateStr)

	var lastOrder models.Order
	err := r.db.Unscoped().Where("order_number LIKE ?", prefix+"%").
		Order("order_number DESC").First(&lastOrder).Error

	var seq int
	if err != nil {
		seq = 1
	} else {
		fmt.Sscanf(lastOrder.OrderNumber, "INV/"+dateStr+"/%d", &seq)
		seq++
	}

	return fmt.Sprintf("INV/%s/%04d", dateStr, seq), nil
}

func (r *OrderRepository) FindAll(page, limit int, status, search string) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	query := r.db.Model(&models.Order{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if search != "" {
		query = query.Where("order_number ILIKE ?", "%"+search+"%")
	}

	query.Count(&total)

	offset := (page - 1) * limit
	err := query.
		Preload("Customer").
		Preload("Details").
		Preload("Details.Service").
		Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&orders).Error

	return orders, total, err
}

func (r *OrderRepository) FindByID(id uint) (*models.Order, error) {
	var order models.Order
	err := r.db.
		Preload("Customer").
		Preload("User").
		Preload("Details").
		Preload("Details.Service").
		Preload("Tracking").
		First(&order, id).Error
	return &order, err
}

func (r *OrderRepository) Create(tx *gorm.DB, order *models.Order) error {
	return tx.Create(order).Error
}

func (r *OrderRepository) Update(order *models.Order) error {
	return r.db.Save(order).Error
}

func (r *OrderRepository) Delete(id uint) error {
	return r.db.Delete(&models.Order{}, id).Error
}

func (r *OrderRepository) AddTracking(tx *gorm.DB, tracking *models.OrderTracking) error {
	return tx.Create(tracking).Error
}

func (r *OrderRepository) UpdateStatus(tx *gorm.DB, id uint, status string) error {
	return tx.Model(&models.Order{}).Where("id = ?", id).Update("status", status).Error
}

func (r *OrderRepository) BeginTx() *gorm.DB {
	return r.db.Begin()
}

func (r *OrderRepository) GetStatusCounts() (map[string]int64, error) {
	type StatusCount struct {
		Status string
		Count  int64
	}
	var results []StatusCount
	err := r.db.Model(&models.Order{}).
		Select("status, count(*) as count").
		Group("status").
		Find(&results).Error
	if err != nil {
		return nil, err
	}

	counts := make(map[string]int64)
	for _, r := range results {
		counts[r.Status] = r.Count
	}
	return counts, nil
}
