package services

import (
	"time"

	"github.com/afdhalpower/golanglaundry/internal/repositories"
)

type ReportService struct {
	repo *repositories.ReportRepository
}

func NewReportService(repo *repositories.ReportRepository) *ReportService {
	return &ReportService{repo: repo}
}

type ReportData struct {
	StartDate    string
	EndDate      string
	TotalRevenue float64
	TotalExpense float64
	Profit       float64
	TotalOrders  int64
	RevenueDaily []repositories.RevenueReportItem
	OrderStats   []repositories.OrderReportItem
	TopCustomers []repositories.TopCustomerItem
	TopServices  []repositories.TopServiceItem
	Expenses     []interface{}
}

func (s *ReportService) GetRevenueReport(start, end time.Time) (float64, []repositories.RevenueReportItem, error) {
	total, err := s.repo.GetRevenueByDateRange(start, end)
	if err != nil {
		return 0, nil, err
	}
	daily, err := s.repo.GetRevenueDaily(start, end)
	if err != nil {
		return 0, nil, err
	}
	return total, daily, nil
}

func (s *ReportService) GetExpenseReport(start, end time.Time) (float64, error) {
	return s.repo.GetExpenseByDateRange(start, end)
}

func (s *ReportService) GetProfitReport(start, end time.Time) (float64, float64, float64, error) {
	revenue, err := s.repo.GetRevenueByDateRange(start, end)
	if err != nil {
		return 0, 0, 0, err
	}
	expense, err := s.repo.GetExpenseByDateRange(start, end)
	if err != nil {
		return 0, 0, 0, err
	}
	return revenue, expense, revenue - expense, nil
}

func (s *ReportService) GetOrderReport(start, end time.Time) (int64, []repositories.OrderReportItem, error) {
	count, err := s.repo.GetOrderTotalCount(start, end)
	if err != nil {
		return 0, nil, err
	}
	stats, err := s.repo.GetOrderStats(start, end)
	if err != nil {
		return 0, nil, err
	}
	return count, stats, nil
}

func (s *ReportService) GetTopCustomers() ([]repositories.TopCustomerItem, error) {
	return s.repo.GetTopCustomers(10)
}

func (s *ReportService) GetTopServices() ([]repositories.TopServiceItem, error) {
	return s.repo.GetTopServices(10)
}
