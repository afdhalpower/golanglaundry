package repositories

import (
	"github.com/afdhalpower/golanglaundry/internal/models"
	"gorm.io/gorm"
)

type InventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) *InventoryRepository {
	return &InventoryRepository{db: db}
}

func (r *InventoryRepository) FindAll(page, limit int, search string) ([]models.Inventory, int64, error) {
	var items []models.Inventory
	var total int64
	query := r.db.Model(&models.Inventory{})
	if search != "" {
		query = query.Where("name ILIKE ? OR category ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	query.Count(&total)
	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Order("name ASC").Find(&items).Error
	return items, total, err
}

func (r *InventoryRepository) FindByID(id uint) (*models.Inventory, error) {
	var item models.Inventory
	err := r.db.First(&item, id).Error
	return &item, err
}

func (r *InventoryRepository) Create(item *models.Inventory) error {
	return r.db.Create(item).Error
}

func (r *InventoryRepository) Update(item *models.Inventory) error {
	return r.db.Save(item).Error
}

func (r *InventoryRepository) Delete(id uint) error {
	return r.db.Delete(&models.Inventory{}, id).Error
}

func (r *InventoryRepository) FindLowStock() ([]models.Inventory, error) {
	var items []models.Inventory
	err := r.db.Where("stock <= min_stock AND min_stock > 0").Order("stock ASC").Find(&items).Error
	return items, err
}
