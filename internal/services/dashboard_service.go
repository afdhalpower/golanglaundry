package services

import (
	"time"

	"github.com/afdhalpower/golanglaundry/internal/repositories"
)

type DashboardStats struct {
	TotalOrdersToday     int64
	OrdersInProgress     int64
	OrdersCompleted      int64
	RevenueToday         float64
	RevenueThisMonth     float64
	TotalCustomers       int64
	TotalExpensesThisMonth float64
}

type DashboardService struct {
	dashboardRepo *repositories.DashboardRepository
}

func NewDashboardService(dashboardRepo *repositories.DashboardRepository) *DashboardService {
	return &DashboardService{dashboardRepo: dashboardRepo}
}

func (s *DashboardService) GetStats() (*DashboardStats, error) {
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	ordersToday, err := s.dashboardRepo.CountOrdersSince(todayStart)
	if err != nil {
		return nil, err
	}

	ordersInProgress, err := s.dashboardRepo.CountOrdersByStatus([]string{"dicuci", "dikeringkan", "disetrika", "siap_diambil"})
	if err != nil {
		return nil, err
	}

	ordersCompleted, err := s.dashboardRepo.CountOrdersByStatus([]string{"sudah_diambil"})
	if err != nil {
		return nil, err
	}

	revenueToday, err := s.dashboardRepo.SumPaymentsSince(todayStart)
	if err != nil {
		return nil, err
	}

	revenueMonth, err := s.dashboardRepo.SumPaymentsSince(monthStart)
	if err != nil {
		return nil, err
	}

	totalCustomers, err := s.dashboardRepo.CountAllCustomers()
	if err != nil {
		return nil, err
	}

	expensesMonth, err := s.dashboardRepo.SumExpensesSince(monthStart)
	if err != nil {
		return nil, err
	}

	return &DashboardStats{
		TotalOrdersToday:       ordersToday,
		OrdersInProgress:       ordersInProgress,
		OrdersCompleted:        ordersCompleted,
		RevenueToday:           revenueToday,
		RevenueThisMonth:       revenueMonth,
		TotalCustomers:         totalCustomers,
		TotalExpensesThisMonth: expensesMonth,
	}, nil
}
