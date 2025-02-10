package service

import (
	"books_api/models"
	"books_api/repository" // Certifique-se de que este pacote exista e exporte as funções HashPassword e ComparePassword
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type AuthService struct {
	UserRepo *repository.UserRepository
}

// NewAuthService cria uma nova instância de AuthService.
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
		if err.Error() == "usuário não encontrado" {
			return "", errors.New("usuário ou senha inválidos")
		}
		return "", fmt.Errorf("erro ao buscar usuário: %w", err)
	}

	//if !user.CheckPassword(password) {
	//	return "", errors.New("credenciais inválidas")
	//}

	// Define os claims do token
	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
		"jti": uuid.NewString(),
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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("erro ao gerar hash da senha: %w", err)
	}

	user := &models.User{
		Username: username,
		Password: string(hashedPassword),
	}

	if err := s.UserRepo.Create(user); err != nil {
		if err != nil && err.Error() == "nome de usuário já está em uso" {
			return errors.New("nome de usuário já está em uso")
		}
		return fmt.Errorf("erro ao criar usuário: %w", err)
	}

	return nil
}
