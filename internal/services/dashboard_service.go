package services

import (
	"fmt"
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

type ChartData struct {
	Revenue     []float64
	OrderCounts []int64
	Labels      []string
}

func (s *DashboardService) GetChartData() (*ChartData, error) {
	const days = 6 // last 7 days (today and 6 days back)

	revenueRows, err := s.dashboardRepo.GetDailyRevenue(days)
	if err != nil {
		return nil, err
	}

	orderRows, err := s.dashboardRepo.GetDailyOrderCounts(days)
	if err != nil {
		return nil, err
	}

	// Build lookup maps from DB results
	revenueMap := make(map[string]float64, len(revenueRows))
	for _, r := range revenueRows {
		revenueMap[r.Date.Format("2006-01-02")] = r.Total
	}

	orderMap := make(map[string]int64, len(orderRows))
	for _, r := range orderRows {
		orderMap[r.Date.Format("2006-01-02")] = r.Count
	}

	// Indonesian day labels (dynamic, aligned with data order)
	labelMap := map[string]string{
		"Monday": "Sen", "Tuesday": "Sel", "Wednesday": "Rab",
		"Thursday": "Kam", "Friday": "Jum", "Saturday": "Sab", "Sunday": "Min",
	}

	labels := make([]string, 7)
	revenue := make([]float64, 7)
	orderCounts := make([]int64, 7)
	now := time.Now()

	for i := 6; i >= 0; i-- {
		date := now.AddDate(0, 0, -i)
		dateStr := date.Format("2006-01-02")
		idx := 6 - i // 0 = 6 days ago, 6 = today
		revenue[idx] = revenueMap[dateStr]
		orderCounts[idx] = orderMap[dateStr]
		labels[idx] = labelMap[date.Weekday().String()]
	}

	return &ChartData{
		Revenue:     revenue,
		OrderCounts: orderCounts,
		Labels:      labels,
	}, nil
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

// GetRecentTracking returns recent tracking activity with customer names.
func (s *DashboardService) GetRecentTracking(limit int) ([]repositories.RecentTrackingItem, error) {
	return s.dashboardRepo.GetRecentTracking(limit)
}

// CountOverdueOrders returns how many orders are past their estimated done date.
func (s *DashboardService) CountOverdueOrders() (int64, error) {
	return s.dashboardRepo.CountOverdueOrders()
}

// GetOverdueOrders returns the list of overdue orders.
func (s *DashboardService) GetOverdueOrders() ([]repositories.OverdueOrderItem, error) {
	return s.dashboardRepo.GetOverdueOrders()
}

// TimeAgo converts a time to a human-readable relative time string in Indonesian.
func TimeAgo(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "baru saja"
	case d < 2*time.Minute:
		return "1 menit lalu"
	case d < time.Hour:
		return fmt.Sprintf("%d menit lalu", int(d.Minutes()))
	case d < 2*time.Hour:
		return "1 jam lalu"
	case d < 24*time.Hour:
		return fmt.Sprintf("%d jam lalu", int(d.Hours()))
	case d < 48*time.Hour:
		return "kemarin"
	default:
		return fmt.Sprintf("%d hari lalu", int(d.Hours()/24))
	}
}

// StatusIcon returns an emoji icon based on the order status.
func StatusIcon(status string) string {
	icons := map[string]string{
		"menunggu":      "🕐",
		"dicuci":        "🫧",
		"dikeringkan":   "💨",
		"disetrika":     "🔥",
		"siap_diambil":  "✅",
		"sudah_diambil": "📦",
		"dibatalkan":    "❌",
	}
	if icon, ok := icons[status]; ok {
		return icon
	}
	return "📋"
}

// StatusLabel returns a human-readable Indonesian status label.
func StatusLabel(status string) string {
	labels := map[string]string{
		"menunggu":      "Menunggu",
		"dicuci":        "Dicuci",
		"dikeringkan":   "Dikeringkan",
		"disetrika":     "Disetrika",
		"siap_diambil":  "Siap Diambil",
		"sudah_diambil": "Sudah Diambil",
		"dibatalkan":    "Dibatalkan",
	}
	if label, ok := labels[status]; ok {
		return label
	}
	return status
}
