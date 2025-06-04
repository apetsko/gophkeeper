package crypto

import (
	"context"
	"testing"

	"github.com/apetsko/gophkeeper/models"
	"github.com/stretchr/testify/require"
)

type mockStorage struct{}

func (m *mockStorage) SaveUserData(ctx context.Context, userData *models.DBUserData) (int, error) {
	return 0, nil
}

func TestEncryptDecryptUserData(t *testing.T) {
	ctx := context.Background()
	masterKey := []byte("01234567890123456789012345678901") // 32 bytes master key

	e := NewEnvelope(&mockStorage{})

	plaintext := []byte("some very secret data")

	encryptedData, err := e.EncryptUserData(ctx, masterKey, plaintext)
	require.NoError(t, err)
	require.NotNil(t, encryptedData)
	require.NotEmpty(t, encryptedData.EncryptedData)
	require.NotEmpty(t, encryptedData.DataNonce)
	require.NotEmpty(t, encryptedData.EncryptedDek)
	require.NotEmpty(t, encryptedData.DekNonce)

	// Превратим EncryptedData в модель DBUserData для дешифровки
	dbUserData := models.DBUserData{
		EncryptedData: encryptedData.EncryptedData,
		DataNonce:     encryptedData.DataNonce,
		EncryptedDek:  encryptedData.EncryptedDek,
		DekNonce:      encryptedData.DekNonce,
	}

	decrypted, err := e.DecryptUserData(ctx, dbUserData, masterKey)
	require.NoError(t, err)
	require.Equal(t, plaintext, decrypted)
}

func TestDecryptUserData_WrongMasterKey(t *testing.T) {
	ctx := context.Background()
	masterKey := []byte("01234567890123456789012345678901") // 32 bytes master key
	wrongKey := []byte("99999999999999999999999999999999")

	e := NewEnvelope(&mockStorage{})

	plaintext := []byte("some very secret data")

	encryptedData, err := e.EncryptUserData(ctx, masterKey, plaintext)
	require.NoError(t, err)

	dbUserData := models.DBUserData{
		EncryptedData: encryptedData.EncryptedData,
		DataNonce:     encryptedData.DataNonce,
		EncryptedDek:  encryptedData.EncryptedDek,
		DekNonce:      encryptedData.DekNonce,
	}

	_, err = e.DecryptUserData(ctx, dbUserData, wrongKey)
	require.Error(t, err)
}

func TestDecryptUserData_CorruptedData(t *testing.T) {
	ctx := context.Background()
	masterKey := []byte("01234567890123456789012345678901") // 32 bytes master key

	e := NewEnvelope(&mockStorage{})

	plaintext := []byte("some very secret data")

	encryptedData, err := e.EncryptUserData(ctx, masterKey, plaintext)
	require.NoError(t, err)

	dbUserData := models.DBUserData{
		EncryptedData: encryptedData.EncryptedData,
		DataNonce:     encryptedData.DataNonce,
		EncryptedDek:  encryptedData.EncryptedDek,
		DekNonce:      encryptedData.DekNonce,
	}

	// Искажём зашифрованные данные
	dbUserData.EncryptedData[0] ^= 0xFF

	_, err = e.DecryptUserData(ctx, dbUserData, masterKey)
	require.Error(t, err)
}
