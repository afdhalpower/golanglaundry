package services

import "github.com/afdhalpower/golanglaundry/internal/repositories"

type SettingService struct {
	repo *repositories.SettingRepository
}

func NewSettingService(repo *repositories.SettingRepository) *SettingService {
	return &SettingService{repo: repo}
}

func (s *SettingService) GetAll() (map[string]string, error) {
	return s.repo.GetAll()
}

func (s *SettingService) Update(settings map[string]string) error {
	for key, value := range settings {
		if err := s.repo.Set(key, value); err != nil {
			return err
		}
	}
	return nil
}
