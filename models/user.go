package models

import (
	"errors"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Username string `json:"username" gorm:"unique;not null"`
	Password string `json:"-" gorm:"not null"`
}

func (u *User) Validate() error {

	u.Username = strings.TrimSpace(u.Username)

	if len(u.Username) < 3 || len(u.Username) > 30 {
		return errors.New("o nome de usuário deve ter entre 3 e 30 caracteres")
	}
	if matched, _ := regexp.MatchString(`^[a-zA-Z0-9_.-]+$`, u.Username); !matched {
		return errors.New("o nome de usuário só pode conter letras, números, pontos e traços")
	}

	if len(u.Password) < 8 {
		return errors.New("a senha deve ter pelo menos 8 caracteres")
	}

	return nil
}

func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword compara a senha fornecida com o hash armazenado.
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		// Apenas para debug
		// fmt.Printf("Erro ao comparar hash: %v\n", err)
	}
	return err == nil
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	// Se a senha já tiver 60 caracteres, podemos supor que já foi hasheada.
	if len(u.Password) != 60 {
		return u.HashPassword()
	}
	return nil
}
