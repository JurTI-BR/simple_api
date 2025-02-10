package service

import (
	"books_api/models"
	"books_api/repository"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type AuthService struct {
	UserRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{UserRepo: userRepo}
}

// Authenticate verifica as credenciais e gera um token JWT.
func (s *AuthService) Authenticate(username, password string) (string, error) {
	if username == "" || password == "" {
		return "", errors.New("nome de usuário e senha são obrigatórios")
	}

	user, err := s.UserRepo.FindByUsername(username)
	if err != nil {
		return "", fmt.Errorf("credenciais inválidas: %w", err)
	}

	if !user.CheckPassword(password) {
		return "", errors.New("credenciais inválidas")
	}

	// Define os claims do token
	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET não configurado")
	}

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("erro ao assinar token: %w", err)
	}

	return tokenString, nil
}

// Register cria um novo usuário.
func (s *AuthService) Register(username, password string) error {
	if username == "" || password == "" {
		return errors.New("nome de usuário e senha são obrigatórios")
	}

	user := &models.User{Username: username, Password: password}

	// Gera o hash da senha antes de salvar o usuário
	if err := user.HashPassword(); err != nil {
		return fmt.Errorf("erro ao gerar hash da senha: %w", err)
	}

	if err := s.UserRepo.Create(user); err != nil {
		return fmt.Errorf("erro ao criar usuário: %w", err)
	}

	return nil
}
