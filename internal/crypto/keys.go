// Package crypto provides cryptographic utilities for data encryption and decryption using envelope encryption.
package crypto

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/subtle"
	"errors"

	"github.com/apetsko/gophkeeper/internal/constants"
	"github.com/apetsko/gophkeeper/models"
	"golang.org/x/crypto/argon2"
)

// KeyManagerInterface defines methods for managing user master keys.
//
//go:generate mockery --dir ./internal/crypto --name=KeyManagerInterface --output=../mocks/ --case=underscore
type KeyManagerInterface interface {
	GetMasterKey(ctx context.Context, userID int) ([]byte, error)
	GetOrCreateMasterKey(ctx context.Context, userID int, userPassword string, userSalt []byte) ([]byte, error)
}

// Убедимся, что KeyManager реализует интерфейс
var _ KeyManagerInterface = (*KeyManager)(nil)

// KeyStorage defines the interface for persisting and retrieving encrypted master keys.
type KeyStorage interface {
	// GetMasterKey fetches the encrypted master key for the user.
	GetMasterKey(ctx context.Context, userID int) (*models.EncryptedMK, error)
	// SaveMasterKey stores the encrypted master key and its nonce for the user.
	SaveMasterKey(ctx context.Context, userID int, encryptedMK, nonce []byte) (int, error)
}

// KeyManager implements KeyManagerInterface and handles master key encryption and decryption.
type KeyManager struct {
	storage             KeyStorage
	serverEncryptionKey []byte
}

// NewKeyManager creates a new KeyManager with the given storage and server encryption key.
func NewKeyManager(
	storage KeyStorage,
	serverEncryptionKey []byte,
) *KeyManager {
	return &KeyManager{
		storage:             storage,
		serverEncryptionKey: serverEncryptionKey,
	}
}

// GetMasterKey retrieves and decrypts the user's master key using the server encryption key.
func (m *KeyManager) GetMasterKey(ctx context.Context, userID int) ([]byte, error) {
	encryptedMK, err := m.storage.GetMasterKey(ctx, userID)
	if err != nil {
		return nil, err
	}

	block, _ := aes.NewCipher(m.serverEncryptionKey)
	gcm, _ := cipher.NewGCM(block)

	mk, err := gcm.Open(nil, encryptedMK.Nonce, encryptedMK.EncryptedMK, nil)
	if err != nil {
		return nil, err
	}

	return mk, nil
}

// GetOrCreateMasterKey retrieves the user's master key if it exists and validates the password,
// or generates and stores a new master key if not found.
func (m *KeyManager) GetOrCreateMasterKey(
	ctx context.Context,
	userID int,
	userPassword string,
	userSalt []byte,
) ([]byte, error) {
	encryptedMK, err := m.storage.GetMasterKey(ctx, userID)
	if err != nil {
		if errors.Is(err, models.ErrMasterKeyNotFound) {
			return m.generateAndStoreMasterKey(ctx, userID, userPassword, userSalt)
		}
		return nil, err
	}

	block, _ := aes.NewCipher(m.serverEncryptionKey)
	gcm, _ := cipher.NewGCM(block)

	mk, err := gcm.Open(nil, encryptedMK.Nonce, encryptedMK.EncryptedMK, nil)
	if err != nil {
		return nil, err
	}

	computedMK := argon2.IDKey([]byte(userPassword), userSalt, 3, constants.Mem, constants.Threads, uint32(constants.KeyLength))
	if subtle.ConstantTimeCompare(mk, computedMK) != 1 {
		return nil, errors.New("invalid password")
	}

	return mk, nil
}

// generateAndStoreMasterKey generates a new master key from the user's password and salt,
// encrypts it with the server key, stores it, and returns the plaintext master key.
func (m *KeyManager) generateAndStoreMasterKey(
	ctx context.Context,
	userID int,
	userPassword string,
	userSalt []byte,
) ([]byte, error) {
	mk := argon2.IDKey([]byte(userPassword), userSalt, 3, constants.Mem, constants.Threads, uint32(constants.KeyLength))

	block, errBlock := aes.NewCipher(m.serverEncryptionKey)
	if errBlock != nil {
		return nil, errBlock
	}

	gcm, _ := cipher.NewGCM(block)
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, errors.New("error generate nonce")
	}

	encryptedMK := gcm.Seal(nil, nonce, mk, nil)

	_, err := m.storage.SaveMasterKey(ctx, userID, encryptedMK, nonce)

	return mk, err
}
