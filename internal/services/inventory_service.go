package services

import (
	"github.com/afdhalpower/golanglaundry/internal/models"
	"github.com/afdhalpower/golanglaundry/internal/repositories"
)

type InventoryService struct {
	repo                  *repositories.InventoryRepository
	stockMovementRepo     *repositories.StockMovementRepository
	stockMovementSvc      *StockMovementService
}

func NewInventoryService(repo *repositories.InventoryRepository, stockMovementRepo *repositories.StockMovementRepository) *InventoryService {
	return &InventoryService{
		repo:              repo,
		stockMovementRepo: stockMovementRepo,
		stockMovementSvc:  NewStockMovementService(stockMovementRepo),
	}
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

func (s *InventoryService) Create(item *models.Inventory, userID uint) error {
	if err := s.repo.Create(item); err != nil {
		return err
	}

	// Record initial stock as "in" movement
	if item.Stock > 0 {
		_ = s.stockMovementSvc.RecordMovement(item.ID, "in", item.Stock, 0, item.Stock, "Stok awal", userID)
	}

	return nil
}

func (s *InventoryService) Update(item *models.Inventory, userID uint) error {
	// Fetch old version to compare stock
	old, err := s.repo.FindByID(item.ID)
	if err != nil {
		return err
	}

	oldStock := old.Stock

	if err := s.repo.Update(item); err != nil {
		return err
	}

	// Record stock movement if stock changed
	diff := item.Stock - oldStock
	if diff > 0 {
		_ = s.stockMovementSvc.RecordMovement(item.ID, "in", diff, oldStock, item.Stock, "Penambahan stok", userID)
	} else if diff < 0 {
		_ = s.stockMovementSvc.RecordMovement(item.ID, "out", -diff, oldStock, item.Stock, "Pengurangan stok", userID)
	}

	return nil
}

func (s *InventoryService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *InventoryService) GetLowStock() ([]models.Inventory, error) {
	return s.repo.FindLowStock()
}

// RecordMovement allows manual stock movement recording (e.g. via handler)
func (s *InventoryService) RecordMovement(inventoryID uint, movementType string, quantity, prevStock, newStock int, note string, userID uint) error {
	return s.stockMovementSvc.RecordMovement(inventoryID, movementType, quantity, prevStock, newStock, note, userID)
}

// GetMovements returns stock movements for a given inventory item
func (s *InventoryService) GetMovements(inventoryID uint) ([]models.StockMovement, error) {
	return s.stockMovementSvc.GetMovements(inventoryID)
}
