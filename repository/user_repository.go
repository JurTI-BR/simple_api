package repository

import (
	"errors"
	"strings"

	"books_api/models"
	"gorm.io/gorm"
)

var ErrUsernameInUse = errors.New("nome de usuário já está em uso")

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	if err := r.DB.Where("LOWER(username) = ?", strings.ToLower(username)).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("usuário não encontrado")
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Create(user *models.User) error {
	if err := user.Validate(); err != nil {
		return err
	}

	if existingUser, _ := r.FindByUsername(user.Username); existingUser != nil {
		return ErrUsernameInUse
	}

	return r.DB.Create(user).Error
}
