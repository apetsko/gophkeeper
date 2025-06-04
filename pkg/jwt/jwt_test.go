package jwt

import (
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func TestGenerateJWT(t *testing.T) {
	userID := 42
	username := "testuser"
	secret := "mysecret"

	tokenStr, err := GenerateJWT(userID, username, secret)
	require.NoError(t, err)
	require.NotEmpty(t, tokenStr)

	// Parse the token to verify claims
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	require.NoError(t, err)
	require.True(t, token.Valid)

	claims, ok := token.Claims.(jwt.MapClaims)
	require.True(t, ok)
	require.Equal(t, float64(userID), claims["user_id"])
	require.Equal(t, username, claims["name"])
	require.NotZero(t, claims["iat"])
}
