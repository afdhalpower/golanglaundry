package repositories

import (
	"github.com/afdhalpower/golanglaundry/internal/models"
	"gorm.io/gorm"
)

type ExpenseCategoryRepository struct {
	db *gorm.DB
}

func NewExpenseCategoryRepository(db *gorm.DB) *ExpenseCategoryRepository {
	return &ExpenseCategoryRepository{db: db}
}

func (r *ExpenseCategoryRepository) FindAll() ([]models.ExpenseCategory, error) {
	var categories []models.ExpenseCategory
	err := r.db.Order("name ASC").Find(&categories).Error
	return categories, err
}

func (r *ExpenseCategoryRepository) Create(name string) (*models.ExpenseCategory, error) {
	cat := &models.ExpenseCategory{Name: name}
	err := r.db.Create(cat).Error
	return cat, err
}

type ExpenseRepository struct {
	db *gorm.DB
}

func NewExpenseRepository(db *gorm.DB) *ExpenseRepository {
	return &ExpenseRepository{db: db}
}

func (r *ExpenseRepository) FindAll(page, limit int, categoryID string) ([]models.Expense, int64, error) {
	var expenses []models.Expense
	var total int64
	query := r.db.Model(&models.Expense{})
	if categoryID != "" {
		query = query.Where("expense_category_id = ?", categoryID)
	}
	query.Count(&total)
	offset := (page - 1) * limit
	err := query.Preload("ExpenseCategory").Offset(offset).Limit(limit).Order("date DESC").Find(&expenses).Error
	return expenses, total, err
}

func (r *ExpenseRepository) FindByID(id uint) (*models.Expense, error) {
	var expense models.Expense
	err := r.db.Preload("ExpenseCategory").First(&expense, id).Error
	return &expense, err
}

func (r *ExpenseRepository) Create(expense *models.Expense) error {
	return r.db.Create(expense).Error
}

func (r *ExpenseRepository) Update(expense *models.Expense) error {
	return r.db.Save(expense).Error
}

func (r *ExpenseRepository) Delete(id uint) error {
	return r.db.Delete(&models.Expense{}, id).Error
}

func (r *ExpenseRepository) GetTotalThisMonth() (float64, error) {
	var total float64
	// Use NowDate from helper - this uses raw SQL for the month aggregation
	// We'll compute it through the service instead
	return total, nil
}
