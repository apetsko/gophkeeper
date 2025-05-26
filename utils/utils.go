package utils

import (
	"crypto/sha256"
	"encoding/base64"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type PasswordHasher interface {
	HashPassword(password string) ([]byte, error)
}

type BcryptHasher struct{}

func (b *BcryptHasher) HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func ComparePassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func GenerateID(s string, length int) (id string) {
	hash := sha256.Sum256([]byte(s))
	id = base64.RawURLEncoding.EncodeToString(hash[:length])[:length]
	return
}

func ValidateStruct(a any) error {
	return validator.New().Struct(a)
}
