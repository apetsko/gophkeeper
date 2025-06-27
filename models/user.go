// Package models defines data structures used throughout the GophKeeper application.
package models

// User represents a user registration or login request.
//
// Fields:
//   - Username: The user's login name (required).
//   - Password: The user's plaintext password (required).
type User struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// UserEntry represents a user record stored in the database.
//
// Fields:
//   - ID: Unique user identifier.
//   - Username: The user's login name.
//   - PasswordHash: The hashed password.
type UserEntry struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
}
