package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(userID int, username, jwtSecret string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"name":    username,
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(jwtSecret))
}
