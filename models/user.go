package models

import (
	"errors"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User representa um usuário no sistema.
type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Username string `json:"username" gorm:"unique;not null"`
	Password string `json:"-" gorm:"not null"`
}

// Valida os dados antes da persistência no banco.
func (u *User) Validate() error {
	// Remove espaços extras no username
	u.Username = strings.TrimSpace(u.Username)

	// Validação do nome de usuário
	if len(u.Username) < 3 || len(u.Username) > 30 {
		return errors.New("o nome de usuário deve ter entre 3 e 30 caracteres")
	}
	if matched, _ := regexp.MatchString(`^[a-zA-Z0-9_.-]+$`, u.Username); !matched {
		return errors.New("o nome de usuário só pode conter letras, números, pontos e traços")
	}

	// Validação da senha
	if len(u.Password) < 8 {
		return errors.New("a senha deve ter pelo menos 8 caracteres")
	}

	return nil
}

// HashPassword criptografa a senha do usuário antes de salvar.
func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword compara a senha fornecida com a armazenada.
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// Antes de salvar, garante que a senha está criptografada.
func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	return u.HashPassword()
}
