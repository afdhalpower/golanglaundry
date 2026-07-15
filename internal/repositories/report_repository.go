package repositories

import (
	"time"

	"github.com/afdhalpower/golanglaundry/internal/models"
	"gorm.io/gorm"
)

type ReportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

type RevenueReportItem struct {
	Date   string  `json:"date"`
	Amount float64 `json:"amount"`
	Count  int64   `json:"count"`
}

type OrderReportItem struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
	Total  float64 `json:"total"`
}

type TopCustomerItem struct {
	ID       uint    `json:"id"`
	Name     string  `json:"name"`
	OrderCount int64 `json:"order_count"`
	Total    float64 `json:"total"`
}

type TopServiceItem struct {
	ID           uint    `json:"id"`
	Name         string  `json:"name"`
	UsageCount   int64   `json:"usage_count"`
	TotalRevenue float64 `json:"total_revenue"`
}

func (r *ReportRepository) GetRevenueByDateRange(start, end time.Time) (float64, error) {
	var total float64
	err := r.db.Model(&models.Payment{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("status = 'lunas' AND payment_date BETWEEN ? AND ?", start, end).
		Scan(&total).Error
	return total, err
}

func (r *ReportRepository) GetRevenueDaily(start, end time.Time) ([]RevenueReportItem, error) {
	var items []RevenueReportItem
	err := r.db.Model(&models.Payment{}).
		Select("DATE(payment_date) as date, SUM(amount) as amount, COUNT(*) as count").
		Where("status = 'lunas' AND payment_date BETWEEN ? AND ?", start, end).
		Group("DATE(payment_date)").
		Order("date ASC").
		Scan(&items).Error
	return items, err
}

func (r *ReportRepository) GetExpenseByDateRange(start, end time.Time) (float64, error) {
	var total float64
	err := r.db.Model(&models.Expense{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("date BETWEEN ? AND ?", start, end).
		Scan(&total).Error
	return total, err
}

func (r *ReportRepository) GetExpenseByCategory(start, end time.Time) ([]models.Expense, error) {
	var expenses []models.Expense
	err := r.db.Preload("ExpenseCategory").
		Where("date BETWEEN ? AND ?", start, end).
		Order("date DESC").
		Find(&expenses).Error
	return expenses, err
}

func (r *ReportRepository) GetOrderStats(start, end time.Time) ([]OrderReportItem, error) {
	var items []OrderReportItem
	err := r.db.Model(&models.Order{}).
		Select("status, COUNT(*) as count, COALESCE(SUM(total), 0) as total").
		Where("created_at BETWEEN ? AND ?", start, end).
		Group("status").
		Order("status").
		Scan(&items).Error
	return items, err
}

func (r *ReportRepository) GetOrderTotalCount(start, end time.Time) (int64, error) {
	var count int64
	err := r.db.Model(&models.Order{}).
		Where("created_at BETWEEN ? AND ?", start, end).
		Count(&count).Error
	return count, err
}

func (r *ReportRepository) GetTopCustomers(limit int) ([]TopCustomerItem, error) {
	var items []TopCustomerItem
	err := r.db.Model(&models.Order{}).
		Select("orders.customer_id as id, c.name, COUNT(*) as order_count, COALESCE(SUM(orders.total), 0) as total").
		Joins("LEFT JOIN customers c ON c.id = orders.customer_id").
		Group("orders.customer_id, c.name").
		Order("total DESC").
		Limit(limit).
		Scan(&items).Error
	return items, err
}

func (r *ReportRepository) GetTopServices(limit int) ([]TopServiceItem, error) {
	var items []TopServiceItem
	err := r.db.Model(&models.OrderDetail{}).
		Select("order_details.service_id as id, s.name, COUNT(*) as usage_count, COALESCE(SUM(order_details.price_per_kg), 0) as total_revenue").
		Joins("LEFT JOIN services s ON s.id = order_details.service_id").
		Group("order_details.service_id, s.name").
		Order("usage_count DESC").
		Limit(limit).
		Scan(&items).Error
	return items, err
}
