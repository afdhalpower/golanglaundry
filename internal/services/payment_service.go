package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/afdhalpower/golanglaundry/internal/models"
	"github.com/afdhalpower/golanglaundry/internal/repositories"
)

type PaymentService struct {
	paymentRepo *repositories.PaymentRepository
	orderRepo   *repositories.OrderRepository
}

func NewPaymentService(paymentRepo *repositories.PaymentRepository, orderRepo *repositories.OrderRepository) *PaymentService {
	return &PaymentService{paymentRepo: paymentRepo, orderRepo: orderRepo}
}

func (s *PaymentService) GetAll(page, limit int, status, search string) ([]models.Payment, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	return s.paymentRepo.FindAll(page, limit, status, search)
}

func (s *PaymentService) GetByID(id uint) (*models.Payment, error) {
	return s.paymentRepo.FindByID(id)
}

func (s *PaymentService) GetByOrderID(orderID uint) (*models.Payment, error) {
	return s.paymentRepo.FindByOrderID(orderID)
}

func (s *PaymentService) CreateOrUpdate(orderID uint, amount float64, method, note string, userID uint) (*models.Payment, error) {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, errors.New("order tidak ditemukan")
	}

	if amount <= 0 {
		amount = order.Total
	}

	payment := &models.Payment{
		OrderID:     orderID,
		Amount:      amount,
		Method:      method,
		Status:      "lunas",
		PaymentDate: time.Now(),
		Note:        note,
		CreatedBy:   userID,
	}

	existing, err := s.paymentRepo.FindByOrderID(orderID)
	if err == nil {
		// Update existing
		existing.Amount = amount
		existing.Method = method
		existing.Status = "lunas"
		existing.PaymentDate = time.Now()
		existing.Note = note
		return existing, nil
	}

	if err := s.paymentRepo.Create(payment); err != nil {
		return nil, fmt.Errorf("gagal menyimpan pembayaran: %w", err)
	}

	return payment, nil
}

func (s *PaymentService) GetRevenueToday() (float64, error) {
	return s.paymentRepo.GetRevenueToday()
}

func (s *PaymentService) GetRevenueMonth() (float64, error) {
	return s.paymentRepo.GetRevenueMonth()
}
