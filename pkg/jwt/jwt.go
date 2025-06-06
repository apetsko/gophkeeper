// Package jwt offers utilities for generating and handling JSON Web Tokens (JWT)
// for authentication and authorization.
package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateJWT creates a signed JWT token for the given user ID and username.
//
// The token uses HS256 signing and includes user ID, username, and issued-at claims.
//
// Returns the signed JWT string or an error if signing fails.
func GenerateJWT(userID int, username, jwtSecret string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"name":    username,
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(jwtSecret))
}
