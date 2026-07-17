package repositories

import (
	"time"

	"github.com/afdhalpower/golanglaundry/internal/models"
	"gorm.io/gorm"
)

// RecentTrackingItem holds the result of an order_tracking JOIN with orders & customers.
type RecentTrackingItem struct {
	OrderNumber  string
	CustomerName string
	Status       string
	Note         string
	CreatedBy    uint
	CreatedAt    time.Time
}

// OverdueOrderItem holds minimal overdue order info for the alert card.
type OverdueOrderItem struct {
	ID               uint
	OrderNumber      string
	CustomerName     string
	EstimatedDoneDate time.Time
}

type DashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) *DashboardRepository {
	return &DashboardRepository{db: db}
}

func (r *DashboardRepository) CountOrdersSince(since time.Time) (int64, error) {
	var count int64
	err := r.db.Model(&models.Order{}).Where("created_at >= ?", since).Count(&count).Error
	return count, err
}

func (r *DashboardRepository) CountOrdersByStatus(statuses []string) (int64, error) {
	var count int64
	err := r.db.Model(&models.Order{}).Where("status IN ?", statuses).Count(&count).Error
	return count, err
}

func (r *DashboardRepository) SumPaymentsSince(since time.Time) (float64, error) {
	var result struct {
		Total float64
	}
	err := r.db.Model(&models.Payment{}).
		Select("COALESCE(SUM(amount), 0) as total").
		Where("created_at >= ? AND status = 'lunas'", since).
		Scan(&result).Error
	return result.Total, err
}

type DailyRevenue struct {
	Date  time.Time
	Total float64
}

type DailyOrderCount struct {
	Date  time.Time
	Count int64
}

func (r *DashboardRepository) GetDailyRevenue(days int) ([]DailyRevenue, error) {
	var results []DailyRevenue
	since := time.Now().AddDate(0, 0, -days)
	err := r.db.Model(&models.Payment{}).
		Select("DATE(created_at) as date, COALESCE(SUM(amount), 0) as total").
		Where("created_at >= ? AND status = 'lunas'", since).
		Group("DATE(created_at)").
		Order("DATE(created_at) ASC").
		Scan(&results).Error
	return results, err
}

func (r *DashboardRepository) GetDailyOrderCounts(days int) ([]DailyOrderCount, error) {
	var results []DailyOrderCount
	since := time.Now().AddDate(0, 0, -days)
	err := r.db.Model(&models.Order{}).
		Select("DATE(created_at) as date, COUNT(*) as count").
		Where("created_at >= ?", since).
		Group("DATE(created_at)").
		Order("DATE(created_at) ASC").
		Scan(&results).Error
	return results, err
}

func (r *DashboardRepository) CountAllCustomers() (int64, error) {
	var count int64
	err := r.db.Model(&models.Customer{}).Count(&count).Error
	return count, err
}

// GetRecentTracking returns the most recent tracking entries with order & customer info.
func (r *DashboardRepository) GetRecentTracking(limit int) ([]RecentTrackingItem, error) {
	var items []RecentTrackingItem
	err := r.db.Table("order_tracking").
		Select(`orders.order_number, customers.name as customer_name, order_tracking.status,
		        COALESCE(order_tracking.note, '') as note, order_tracking.created_by, order_tracking.created_at`).
		Joins("JOIN orders ON orders.id = order_tracking.order_id").
		Joins("JOIN customers ON customers.id = orders.customer_id").
		Order("order_tracking.created_at DESC").
		Limit(limit).
		Scan(&items).Error
	return items, err
}

// CountOverdueOrders returns the count of orders past their estimated done date and not yet completed/cancelled.
func (r *DashboardRepository) CountOverdueOrders() (int64, error) {
	var count int64
	now := time.Now()
	err := r.db.Model(&models.Order{}).
		Where("estimated_done_date < ? AND status NOT IN ?", now, []string{"sudah_diambil", "dibatalkan"}).
		Count(&count).Error
	return count, err
}

// GetOverdueOrders returns the list of overdue orders with customer name.
func (r *DashboardRepository) GetOverdueOrders() ([]OverdueOrderItem, error) {
	var items []OverdueOrderItem
	now := time.Now()
	err := r.db.Table("orders").
		Select("orders.id, orders.order_number, customers.name as customer_name, orders.estimated_done_date").
		Joins("JOIN customers ON customers.id = orders.customer_id").
		Where("orders.estimated_done_date < ? AND orders.status NOT IN ?", now, []string{"sudah_diambil", "dibatalkan"}).
		Order("orders.estimated_done_date ASC").
		Scan(&items).Error
	return items, err
}

func (r *DashboardRepository) SumExpensesSince(since time.Time) (float64, error) {
	var result struct {
		Total float64
	}
	err := r.db.Model(&models.Expense{}).
		Select("COALESCE(SUM(amount), 0) as total").
		Where("date >= ?", since).
		Scan(&result).Error
	return result.Total, err
}
