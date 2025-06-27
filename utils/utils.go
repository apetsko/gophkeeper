// Package utils contains helper functions for password hashing, ID generation, and struct validation.
package utils

import (
	"crypto/sha256"
	"encoding/base64"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

// PasswordHasher defines the interface for password hashing implementations.
//
// Implementations should provide a method to hash a plaintext password.
type PasswordHasher interface {
	// HashPassword hashes the given plaintext password and returns the hash or an error.
	HashPassword(password string) ([]byte, error)
}

// BcryptHasher implements the PasswordHasher interface using bcrypt.
type BcryptHasher struct{}

// HashPassword hashes the given password using bcrypt.
//
// Returns the hashed password as a byte slice or an error if hashing fails.
func (b *BcryptHasher) HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// ComparePassword compares a bcrypt hashed password with its possible plaintext equivalent.
//
// Returns true if the password matches the hash, false otherwise.
func ComparePassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// HashPassword hashes the given password using bcrypt with the default cost.
//
// Returns the hashed password as a byte slice or an error if hashing fails.
func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// GenerateID generates a base64-encoded, URL-safe ID of the specified length from the SHA-256 hash of the input string.
//
// Parameters:
//   - s: The input string to hash.
//   - length: The desired length of the resulting ID.
//
// Returns the generated ID as a string.
func GenerateID(s string, length int) (id string) {
	hash := sha256.Sum256([]byte(s))
	id = base64.RawURLEncoding.EncodeToString(hash[:length])[:length]
	return
}

// ValidateStruct validates a struct using the go-playground/validator package.
//
// Returns an error if validation fails, or nil if the struct is valid.
func ValidateStruct(a any) error {
	return validator.New().Struct(a)
}
