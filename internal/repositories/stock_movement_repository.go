package repositories

import (
	"github.com/afdhalpower/golanglaundry/internal/models"
	"gorm.io/gorm"
)

type StockMovementRepository struct {
	db *gorm.DB
}

func NewStockMovementRepository(db *gorm.DB) *StockMovementRepository {
	return &StockMovementRepository{db: db}
}

func (r *StockMovementRepository) Create(movement *models.StockMovement) error {
	return r.db.Create(movement).Error
}

func (r *StockMovementRepository) FindByInventoryID(inventoryID uint) ([]models.StockMovement, error) {
	var movements []models.StockMovement
	err := r.db.Where("inventory_id = ?", inventoryID).
		Preload("User").
		Order("created_at DESC").
		Find(&movements).Error
	return movements, err
}

func (r *StockMovementRepository) FindAll() ([]models.StockMovement, error) {
	var movements []models.StockMovement
	err := r.db.Preload("Inventory").Preload("User").
		Order("created_at DESC").
		Find(&movements).Error
	return movements, err
}
