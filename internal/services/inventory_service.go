package services

import (
	"github.com/afdhalpower/golanglaundry/internal/models"
	"github.com/afdhalpower/golanglaundry/internal/repositories"
)

type InventoryService struct {
	repo *repositories.InventoryRepository
}

func NewInventoryService(repo *repositories.InventoryRepository) *InventoryService {
	return &InventoryService{repo: repo}
}

func (s *InventoryService) GetAll(page, limit int, search string) ([]models.Inventory, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	return s.repo.FindAll(page, limit, search)
}

func (s *InventoryService) GetByID(id uint) (*models.Inventory, error) {
	return s.repo.FindByID(id)
}

func (s *InventoryService) Create(item *models.Inventory) error {
	return s.repo.Create(item)
}

func (s *InventoryService) Update(item *models.Inventory) error {
	return s.repo.Update(item)
}

func (s *InventoryService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *InventoryService) GetLowStock() ([]models.Inventory, error) {
	return s.repo.FindLowStock()
}
