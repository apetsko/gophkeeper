package crypto

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"testing"

	"github.com/apetsko/gophkeeper/models"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/argon2"
)

type mockKeyStorage struct {
	storedMK       *models.EncryptedMK
	saveShouldFail bool
	getShouldFail  bool
	getErr         error
}

func (m *mockKeyStorage) GetMasterKey(ctx context.Context, userID int) (*models.EncryptedMK, error) {
	if m.getShouldFail {
		return nil, m.getErr
	}
	if m.storedMK == nil {
		return nil, models.ErrMasterKeyNotFound
	}
	return m.storedMK, nil
}

func (m *mockKeyStorage) SaveMasterKey(ctx context.Context, userID int, encryptedMK, nonce []byte) (int, error) {
	if m.saveShouldFail {
		return 0, errors.New("save failed")
	}
	m.storedMK = &models.EncryptedMK{
		EncryptedMK: encryptedMK,
		Nonce:       nonce,
	}
	return 1, nil
}

func generateEncryptedMK(serverKey, mk []byte) (*models.EncryptedMK, error) {
	block, err := aes.NewCipher(serverKey)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}
	encrypted := gcm.Seal(nil, nonce, mk, nil)
	return &models.EncryptedMK{
		EncryptedMK: encrypted,
		Nonce:       nonce,
	}, nil
}

func TestGetMasterKey_Success(t *testing.T) {
	ctx := context.Background()
	serverKey := []byte("01234567890123456789012345678901")
	userID := 123
	mk := []byte("verysecretmasterkeyverysecretmas") // 32 bytes

	encryptedMK, err := generateEncryptedMK(serverKey, mk)
	require.NoError(t, err)

	storage := &mockKeyStorage{
		storedMK: encryptedMK,
	}

	km := NewKeyManager(storage, serverKey)

	gotMK, err := km.GetMasterKey(ctx, userID)
	require.NoError(t, err)
	require.Equal(t, mk, gotMK)
}

func TestGetMasterKey_NotFound(t *testing.T) {
	ctx := context.Background()
	serverKey := []byte("01234567890123456789012345678901")
	userID := 123

	storage := &mockKeyStorage{}

	km := NewKeyManager(storage, serverKey)

	_, err := km.GetMasterKey(ctx, userID)
	require.ErrorIs(t, err, models.ErrMasterKeyNotFound)
}

func TestGetOrCreateMasterKey_CreateNew_Success(t *testing.T) {
	ctx := context.Background()
	serverKey := []byte("01234567890123456789012345678901")
	userID := 123
	userPassword := "strongpassword"
	userSalt := []byte("saltysalt")

	storage := &mockKeyStorage{}

	km := NewKeyManager(storage, serverKey)

	mk, err := km.GetOrCreateMasterKey(ctx, userID, userPassword, userSalt)
	require.NoError(t, err)
	require.NotNil(t, mk)
	require.Len(t, mk, 32)

	// Check that key was saved in storage
	require.NotNil(t, storage.storedMK)
}

func TestGetOrCreateMasterKey_InvalidPassword(t *testing.T) {
	ctx := context.Background()
	serverKey := []byte("01234567890123456789012345678901")
	userID := 123
	userPassword := "correctpassword"
	userSalt := []byte("saltysalt")

	// Создаём ключ из пароля, но в тесте будем использовать неправильный
	km := NewKeyManager(nil, serverKey)

	// Создаём зашифрованный MK для правильного пароля
	mk := km.generateMasterKeyForTest(userPassword, userSalt, serverKey, t)

	storage := &mockKeyStorage{
		storedMK: mk,
	}

	km.storage = storage

	// Пытаемся получить с неправильным паролем
	_, err := km.GetOrCreateMasterKey(ctx, userID, "wrongpassword", userSalt)
	require.Error(t, err)
	require.Equal(t, "invalid password", err.Error())
}

func TestGenerateAndStoreMasterKey_SaveFail(t *testing.T) {
	ctx := context.Background()
	serverKey := []byte("01234567890123456789012345678901")
	userID := 123
	userPassword := "password"
	userSalt := []byte("saltysalt")

	storage := &mockKeyStorage{saveShouldFail: true}

	km := NewKeyManager(storage, serverKey)

	_, err := km.GetOrCreateMasterKey(ctx, userID, userPassword, userSalt)
	require.Error(t, err)
	require.Contains(t, err.Error(), "save failed")
}

// Вспомогательная функция, чтобы сгенерировать валидный EncryptedMK для теста invalid password
func (m *KeyManager) generateMasterKeyForTest(userPassword string, userSalt, serverKey []byte, t *testing.T) *models.EncryptedMK {
	mk := argon2.IDKey([]byte(userPassword), userSalt, 3, 64*1024, 4, 32)
	block, err := aes.NewCipher(serverKey)
	require.NoError(t, err)
	gcm, err := cipher.NewGCM(block)
	require.NoError(t, err)
	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	require.NoError(t, err)
	encryptedMK := gcm.Seal(nil, nonce, mk, nil)
	return &models.EncryptedMK{
		EncryptedMK: encryptedMK,
		Nonce:       nonce,
	}
}
