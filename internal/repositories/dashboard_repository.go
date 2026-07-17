package repositories

import (
	"time"

	"github.com/afdhalpower/golanglaundry/internal/models"
	"gorm.io/gorm"
)

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
