package repository

import (
	"errors"
	"gorm.io/gorm"

	"books_api/models"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	if err := r.DB.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("usuário não encontrado")
		}
		return nil, err
	}
	return &user, nil
}

// Create cria um novo usuário no banco de dados.
func (r *UserRepository) Create(user *models.User) error {
	// Valida os dados do usuário antes de salvar
	if err := user.Validate(); err != nil {
		return err
	}

	existingUser, _ := r.FindByUsername(user.Username)
	if existingUser != nil {
		return errors.New("nome de usuário já está em uso")
	}

	return r.DB.Create(user).Error
}
