package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestBcryptHasher_HashPassword(t *testing.T) {
	hasher := &BcryptHasher{}
	password := "mySecret"
	hash, err := hasher.HashPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NoError(t, bcrypt.CompareHashAndPassword(hash, []byte(password)))
}

func TestComparePassword(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("pass123"), bcrypt.DefaultCost)
	assert.True(t, ComparePassword(string(hash), "pass123"))
	assert.False(t, ComparePassword(string(hash), "wrongpass"))
}

func TestHashPassword(t *testing.T) {
	password := "anotherSecret"
	hash, err := HashPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NoError(t, bcrypt.CompareHashAndPassword(hash, []byte(password)))
}

func TestGenerateID(t *testing.T) {
	id10 := GenerateID("input", 10)
	id20 := GenerateID("input", 20)
	assert.Len(t, id10, 10)
	assert.Len(t, id20, 20)
	assert.NotEqual(t, id10, id20)
}

func TestValidateStruct(t *testing.T) {
	type testStruct struct {
		Field1 string `validate:"required"`
		Field2 int    `validate:"min=1"`
	}
	valid := testStruct{Field1: "ok", Field2: 2}
	invalid := testStruct{Field1: "", Field2: 0}

	assert.NoError(t, ValidateStruct(valid))
	assert.Error(t, ValidateStruct(invalid))
}
