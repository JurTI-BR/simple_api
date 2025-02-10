package service

import (
	"books_api/models"
	"books_api/repository"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

var (
	ErrMissingCredentials   = errors.New("nome de usuário e senha são obrigatórios")
	ErrInvalidCredentials   = errors.New("usuário ou senha inválidos")
	ErrJWTSecretNotProvided = errors.New("JWT_SECRET não configurado")
)

type AuthService struct {
	UserRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{UserRepo: userRepo}
}

// Authenticate realiza a autenticação do usuário e gera um token JWT.
func (s *AuthService) Authenticate(username, password string) (string, error) {
	if username == "" || password == "" {
		return "", ErrMissingCredentials
	}

	user, err := s.UserRepo.FindByUsername(username)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	if !user.CheckPassword(password) {
		return "", ErrInvalidCredentials
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": now.Add(24 * time.Hour).Unix(),
		"iat": now.Unix(),
		"jti": uuid.NewString(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", ErrJWTSecretNotProvided
	}
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("erro ao assinar token: %w", err)
	}

	return tokenString, nil
}

func (s *AuthService) Register(username, password string) error {
	if username == "" || password == "" {
		return ErrMissingCredentials
	}

	user := &models.User{
		Username: username,
		Password: password,
	}

	if err := s.UserRepo.Create(user); err != nil {
		if errors.Is(err, repository.ErrUsernameInUse) {
			return repository.ErrUsernameInUse
		}
		return fmt.Errorf("erro ao criar usuário: %w", err)
	}

	return nil
}
