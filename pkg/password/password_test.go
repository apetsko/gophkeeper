package password

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHashPasswordAndCheckPasswordHash(t *testing.T) {
	password := "supersecret123"

	hash, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hash)
	require.NotEqual(t, password, hash)

	ok := CheckPasswordHash(password, hash)
	require.True(t, ok, "password should match hash")

	wrong := CheckPasswordHash("wrongpassword", hash)
	require.False(t, wrong, "wrong password should not match hash")
}

func TestHashPassword_ErrorOnEmpty(t *testing.T) {
	hash, err := HashPassword("")
	require.NoError(t, err)
	require.NotEmpty(t, hash)
}
