package services

import (
	"github.com/afdhalpower/golanglaundry/internal/models"
	"github.com/afdhalpower/golanglaundry/internal/repositories"
)

type ServiceService struct {
	repo *repositories.ServiceRepository
}

func NewServiceService(repo *repositories.ServiceRepository) *ServiceService {
	return &ServiceService{repo: repo}
}

func (s *ServiceService) GetAll(page, limit int, search string) ([]models.Service, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	return s.repo.FindAll(page, limit, search)
}

func (s *ServiceService) GetByID(id uint) (*models.Service, error) {
	return s.repo.FindByID(id)
}

func (s *ServiceService) Create(service *models.Service) error {
	return s.repo.Create(service)
}

func (s *ServiceService) Update(service *models.Service) error {
	return s.repo.Update(service)
}

func (s *ServiceService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *ServiceService) GetActive() ([]models.Service, error) {
	return s.repo.FindActive()
}
