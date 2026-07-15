package services

import (
	"github.com/afdhalpower/golanglaundry/internal/models"
	"github.com/afdhalpower/golanglaundry/internal/repositories"
)

type CustomerService struct {
	repo *repositories.CustomerRepository
}

func NewCustomerService(repo *repositories.CustomerRepository) *CustomerService {
	return &CustomerService{repo: repo}
}

func (s *CustomerService) GetAll(page, limit int, search string) ([]models.Customer, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	return s.repo.FindAll(page, limit, search)
}

func (s *CustomerService) GetByID(id uint) (*models.Customer, error) {
	return s.repo.FindByID(id)
}

func (s *CustomerService) Create(customer *models.Customer) error {
	return s.repo.Create(customer)
}

func (s *CustomerService) Update(customer *models.Customer) error {
	return s.repo.Update(customer)
}

func (s *CustomerService) Delete(id uint) error {
	return s.repo.Delete(id)
}
