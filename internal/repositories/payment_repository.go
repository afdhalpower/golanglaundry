package repositories

import (
	"time"

	"github.com/afdhalpower/golanglaundry/internal/models"
	"gorm.io/gorm"
)

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) FindAll(page, limit int, status, search string) ([]models.Payment, int64, error) {
	var payments []models.Payment
	var total int64
	query := r.db.Model(&models.Payment{})
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if search != "" {
		query = query.Joins("JOIN orders ON orders.id = payments.order_id").
			Where("orders.order_number ILIKE ?", "%"+search+"%")
	}
	query.Count(&total)
	offset := (page - 1) * limit
	err := query.Preload("Order").Preload("Order.Customer").Offset(offset).Limit(limit).Order("created_at DESC").Find(&payments).Error
	return payments, total, err
}

func (r *PaymentRepository) FindByID(id uint) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.Preload("Order").Preload("Order.Customer").First(&payment, id).Error
	return &payment, err
}

func (r *PaymentRepository) FindByOrderID(orderID uint) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.Where("order_id = ?", orderID).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, err
}

func (r *PaymentRepository) Create(payment *models.Payment) error {
	return r.db.Create(payment).Error
}

func (r *PaymentRepository) GetRevenueToday() (float64, error) {
	var total float64
	today := time.Now().Truncate(24 * time.Hour)
	err := r.db.Model(&models.Payment{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("status = 'lunas' AND payment_date >= ?", today).
		Scan(&total).Error
	return total, err
}

func (r *PaymentRepository) GetRevenueMonth() (float64, error) {
	var total float64
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	err := r.db.Model(&models.Payment{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("status = 'lunas' AND payment_date >= ?", startOfMonth).
		Scan(&total).Error
	return total, err
}
