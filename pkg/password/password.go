// Package password provides password hashing and verification utilities using bcrypt.
package password

import "golang.org/x/crypto/bcrypt"

const passwordCost = 14

// HashPassword hashes the given password using bcrypt with a predefined cost.
//
// Returns the hashed password as a string or an error if hashing fails.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), passwordCost)
	return string(bytes), err
}

// CheckPasswordHash compares a bcrypt hashed password with its possible plaintext equivalent.
//
// Returns true if the password matches the hash, false otherwise.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
