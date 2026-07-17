package services

import (
	"errors"
	"time"

	"github.com/afdhalpower/golanglaundry/internal/models"
	"github.com/afdhalpower/golanglaundry/internal/repositories"
)

type ExpenseCategoryService struct {
	repo *repositories.ExpenseCategoryRepository
}

func NewExpenseCategoryService(repo *repositories.ExpenseCategoryRepository) *ExpenseCategoryService {
	return &ExpenseCategoryService{repo: repo}
}

func (s *ExpenseCategoryService) GetAll() ([]models.ExpenseCategory, error) {
	return s.repo.FindAll()
}

func (s *ExpenseCategoryService) Create(name string) (*models.ExpenseCategory, error) {
	if name == "" {
		return nil, errors.New("nama kategori harus diisi")
	}
	return s.repo.Create(name)
}

type ExpenseService struct {
	repo *repositories.ExpenseRepository
}

func NewExpenseService(repo *repositories.ExpenseRepository) *ExpenseService {
	return &ExpenseService{repo: repo}
}

func (s *ExpenseService) GetAll(page, limit int, categoryID, search string) ([]models.Expense, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	return s.repo.FindAll(page, limit, categoryID, search)
}

func (s *ExpenseService) GetByID(id uint) (*models.Expense, error) {
	return s.repo.FindByID(id)
}

func (s *ExpenseService) Create(expense *models.Expense) error {
	return s.repo.Create(expense)
}

func (s *ExpenseService) Update(expense *models.Expense) error {
	return s.repo.Update(expense)
}

func (s *ExpenseService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *ExpenseService) GetTotalThisMonth() (float64, error) {
	// Calculate via repo with date range
	return 0, nil
}

func getStartOfMonth() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
}
