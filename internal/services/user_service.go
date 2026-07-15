package services

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"github.com/afdhalpower/golanglaundry/internal/models"
	"github.com/afdhalpower/golanglaundry/internal/repositories"
)

type UserService struct {
	repo *repositories.UserRepository
}

func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetAll(page, limit int) ([]models.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	return s.repo.FindAll(page, limit)
}

func (s *UserService) GetByID(id uint) (*models.User, error) {
	return s.repo.FindByID(id)
}

func (s *UserService) Create(name, email, password, role string) (*models.User, error) {
	// Check if email already exists
	existing, _ := s.repo.FindByEmail(email)
	if existing != nil && existing.ID > 0 {
		return nil, errors.New("email sudah digunakan")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("gagal hash password: %w", err)
	}

	user := &models.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(hash),
		Role:         role,
		IsActive:     true,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, fmt.Errorf("gagal menyimpan user: %w", err)
	}

	return user, nil
}

func (s *UserService) Update(id uint, name, email, role string, isActive bool) error {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("user tidak ditemukan")
	}

	user.Name = name
	user.Email = email
	user.Role = role
	user.IsActive = isActive

	return s.repo.Update(user)
}

func (s *UserService) ResetPassword(id uint, newPassword string) error {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("user tidak ditemukan")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("gagal hash password: %w", err)
	}

	user.PasswordHash = string(hash)
	return s.repo.Update(user)
}

func (s *UserService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *UserService) GetRoles() []string {
	return []string{"admin", "kasir", "pegawai"}
}
