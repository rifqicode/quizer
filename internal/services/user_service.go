package services

import (
	"errors"
	"quizer/datasource"
	"quizer/internal/models"

	"gorm.io/gorm"
)

type UserService interface {
	CreateUser(user *models.User) error
	GetUserByID(id uint) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetAllUsers() ([]models.User, error)
	UpdateUser(id uint, user *models.User) error
	DeleteUser(id uint) error
}

type userService struct {
	db *gorm.DB
}

func NewUserService() UserService {
	return &userService{
		db: datasource.GetDB(),
	}
}

func (s *userService) CreateUser(user *models.User) error {
	if user.Name == "" || user.Email == "" {
		return errors.New("name and email are required")
	}

	result := s.db.Create(user)
	return result.Error
}

func (s *userService) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	result := s.db.First(&user, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}

	return &user, nil
}

func (s *userService) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	result := s.db.Where("email = ?", email).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}

	return &user, nil
}

func (s *userService) GetAllUsers() ([]models.User, error) {
	var users []models.User
	result := s.db.Find(&users)
	return users, result.Error
}

func (s *userService) UpdateUser(id uint, user *models.User) error {
	existingUser, err := s.GetUserByID(id)
	if err != nil {
		return err
	}

	// Update only non-zero fields
	result := s.db.Model(existingUser).Updates(user)
	return result.Error
}

func (s *userService) DeleteUser(id uint) error {
	result := s.db.Delete(&models.User{}, id)

	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}

	return result.Error
}
