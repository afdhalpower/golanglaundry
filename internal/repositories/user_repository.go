package repositories

import (
	"github.com/afdhalpower/golanglaundry/internal/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *UserRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	return &user, err
}

func (r *UserRepository) UpdatePassword(id uint, hash string) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).Update("password_hash", hash).Error
}

func (r *UserRepository) UpdateProfile(id uint, name, email string) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"name":  name,
		"email": email,
	}).Error
}

func (r *UserRepository) FindAll(page, limit int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	r.db.Model(&models.User{}).Count(&total)

	offset := (page - 1) * limit
	err := r.db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&users).Error
	return users, total, err
}

func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}
