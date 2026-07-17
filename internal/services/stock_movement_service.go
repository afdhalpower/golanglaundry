package services

import (
	"github.com/afdhalpower/golanglaundry/internal/models"
	"github.com/afdhalpower/golanglaundry/internal/repositories"
)

type StockMovementService struct {
	repo *repositories.StockMovementRepository
}

func NewStockMovementService(repo *repositories.StockMovementRepository) *StockMovementService {
	return &StockMovementService{repo: repo}
}

func (s *StockMovementService) RecordMovement(inventoryID uint, movementType string, quantity, prevStock, newStock int, note string, userID uint) error {
	movement := &models.StockMovement{
		InventoryID:   inventoryID,
		Type:          movementType,
		Quantity:      quantity,
		PreviousStock: prevStock,
		NewStock:      newStock,
		Note:          note,
		CreatedBy:     userID,
	}
	return s.repo.Create(movement)
}

func (s *StockMovementService) GetMovements(inventoryID uint) ([]models.StockMovement, error) {
	return s.repo.FindByInventoryID(inventoryID)
}
