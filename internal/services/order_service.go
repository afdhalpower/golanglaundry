package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/afdhalpower/golanglaundry/internal/models"
	"github.com/afdhalpower/golanglaundry/internal/repositories"
)

var validTransitions = map[string][]string{
	"menunggu":     {"dicuci", "dibatalkan"},
	"dicuci":       {"dikeringkan", "dibatalkan"},
	"dikeringkan":  {"disetrika", "dibatalkan"},
	"disetrika":    {"siap_diambil", "dibatalkan"},
	"siap_diambil": {"sudah_diambil", "dibatalkan"},
	"sudah_diambil": {},
	"dibatalkan":   {},
}

type OrderService struct {
	orderRepo    *repositories.OrderRepository
	customerRepo *repositories.CustomerRepository
	serviceRepo  *repositories.ServiceRepository
}

func NewOrderService(
	orderRepo *repositories.OrderRepository,
	customerRepo *repositories.CustomerRepository,
	serviceRepo *repositories.ServiceRepository,
) *OrderService {
	return &OrderService{
		orderRepo:    orderRepo,
		customerRepo: customerRepo,
		serviceRepo:  serviceRepo,
	}
}

type CreateOrderRequest struct {
	CustomerID   uint
	UserID       uint
	ServiceID    uint
	WeightKg     float64
	PricePerKg   float64
	Discount     float64
	ExtraCost    float64
	EntryDate    string
	Notes        string
}

func (s *OrderService) Create(req CreateOrderRequest) (*models.Order, error) {
	// Validate customer exists
	customer, err := s.customerRepo.FindByID(req.CustomerID)
	if err != nil {
		return nil, errors.New("pelanggan tidak ditemukan")
	}
	_ = customer

	// Validate service exists
	service, err := s.serviceRepo.FindByID(req.ServiceID)
	if err != nil {
		return nil, errors.New("layanan tidak ditemukan")
	}

	entryDate, err := time.Parse("2006-01-02", req.EntryDate)
	if err != nil {
		entryDate = time.Now()
	}

	// Calculate total
	total := (service.PricePerKg * req.WeightKg) - req.Discount + req.ExtraCost
	if total < 0 {
		total = 0
	}

	// Generate order number
	orderNumber, err := s.orderRepo.GenerateOrderNumber(entryDate)
	if err != nil {
		return nil, fmt.Errorf("gagal generate nomor order: %w", err)
	}

	estimatedDone := entryDate.Add(time.Duration(service.EstimatedHours) * time.Hour)

	order := &models.Order{
		OrderNumber:      orderNumber,
		CustomerID:       req.CustomerID,
		UserID:           req.UserID,
		WeightKg:         req.WeightKg,
		Discount:         req.Discount,
		ExtraCost:        req.ExtraCost,
		Total:            total,
		EntryDate:        entryDate,
		EstimatedDoneDate: estimatedDone,
		Status:           "menunggu",
		Notes:            req.Notes,
		Details: []models.OrderDetail{
			{
				ServiceID:  req.ServiceID,
				PricePerKg: service.PricePerKg,
			},
		},
		Tracking: []models.OrderTracking{
			{
				Status:    "menunggu",
				Note:      "Pesanan dibuat",
				CreatedBy: req.UserID,
			},
		},
	}

	// Use transaction
	tx := s.orderRepo.BeginTx()
	if err := s.orderRepo.Create(tx, order); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("gagal menyimpan order: %w", err)
	}
	tx.Commit()

	return order, nil
}

func (s *OrderService) GetByID(id uint) (*models.Order, error) {
	return s.orderRepo.FindByID(id)
}

func (s *OrderService) GetAll(page, limit int, status, search string) ([]models.Order, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	return s.orderRepo.FindAll(page, limit, status, search)
}

func (s *OrderService) UpdateStatus(id uint, newStatus, note string, userID uint) error {
	order, err := s.orderRepo.FindByID(id)
	if err != nil {
		return errors.New("order tidak ditemukan")
	}

	validNext, ok := validTransitions[order.Status]
	if !ok {
		return errors.New("status tidak valid")
	}

	isValid := false
	for _, s := range validNext {
		if s == newStatus {
			isValid = true
			break
		}
	}
	if !isValid {
		return fmt.Errorf("tidak bisa mengubah status dari '%s' ke '%s'", order.Status, newStatus)
	}

	tx := s.orderRepo.BeginTx()

	if err := s.orderRepo.UpdateStatus(tx, id, newStatus); err != nil {
		tx.Rollback()
		return err
	}

	tracking := &models.OrderTracking{
		OrderID:   id,
		Status:    newStatus,
		Note:      note,
		CreatedBy: userID,
	}
	if err := s.orderRepo.AddTracking(tx, tracking); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (s *OrderService) Delete(id uint) error {
	return s.orderRepo.Delete(id)
}

func (s *OrderService) GetStatusList() []string {
	return []string{"menunggu", "dicuci", "dikeringkan", "disetrika", "siap_diambil", "sudah_diambil", "dibatalkan"}
}

func (s *OrderService) GetValidNextStatuses(currentStatus string) []string {
	return validTransitions[currentStatus]
}

func (s *OrderService) GetStatusCounts() (map[string]int64, error) {
	return s.orderRepo.GetStatusCounts()
}
