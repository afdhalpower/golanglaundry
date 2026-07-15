package repositories

import (
	"github.com/afdhalpower/golanglaundry/internal/models"
	"gorm.io/gorm"
)

type SettingRepository struct {
	db *gorm.DB
}

func NewSettingRepository(db *gorm.DB) *SettingRepository {
	return &SettingRepository{db: db}
}

func (r *SettingRepository) Get(key string) (string, error) {
	var setting models.Setting
	err := r.db.Where("key = ?", key).First(&setting).Error
	if err != nil {
		return "", err
	}
	return setting.Value, nil
}

func (r *SettingRepository) Set(key, value string) error {
	var setting models.Setting
	err := r.db.Where("key = ?", key).First(&setting).Error
	if err != nil {
		return r.db.Create(&models.Setting{Key: key, Value: value}).Error
	}
	setting.Value = value
	return r.db.Save(&setting).Error
}

func (r *SettingRepository) GetAll() (map[string]string, error) {
	var settings []models.Setting
	if err := r.db.Find(&settings).Error; err != nil {
		return nil, err
	}
	result := make(map[string]string)
	for _, s := range settings {
		result[s.Key] = s.Value
	}
	return result, nil
}
