// Package crypto provides cryptographic utilities for data encryption and decryption using envelope encryption.
package crypto

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"

	"github.com/apetsko/gophkeeper/internal/constants"
	"github.com/apetsko/gophkeeper/models"
)

// IEnvelope defines the interface for envelope encryption and decryption of user data.
//
//go:generate mockery --dir ./internal/crypto --name=IEnvelopeStorage --output=.../mocks/ --case=underscore
type IEnvelope interface {
	// EncryptUserData encrypts the given data with a randomly generated DEK, which is itself encrypted with the master key.
	EncryptUserData(ctx context.Context, masterKey []byte, data []byte) (*models.EncryptedData, error)
	// DecryptUserData decrypts the user data using the provided master key.
	DecryptUserData(ctx context.Context, userData models.DBUserData, masterKey []byte) ([]byte, error)
}

// EnvelopStorage defines the interface for persisting user data.
type EnvelopStorage interface {
	// SaveUserData saves the user data and returns the new record's ID.
	SaveUserData(ctx context.Context, userData *models.DBUserData) (int, error)
}

// Envelope implements the IEnvelope interface and provides envelope encryption logic.
type Envelope struct {
	storage EnvelopStorage
}

// NewEnvelope creates a new Envelope with the given storage backend.
func NewEnvelope(storage EnvelopStorage) *Envelope {
	return &Envelope{
		storage: storage,
	}
}

// EncryptUserData encrypts the provided data using a randomly generated DEK, which is then encrypted with the master key.
// Returns the encrypted data, encrypted DEK, and their nonces.
func (e *Envelope) EncryptUserData(
	ctx context.Context,
	masterKey []byte,
	data []byte,
) (*models.EncryptedData, error) {
	// 1. Генерируем случайный DEK
	dek := make([]byte, constants.KeyLength)
	if _, err := rand.Read(dek); err != nil {
		return nil, fmt.Errorf("failed to generate DEK: %w", err)
	}

	// 2. Шифруем данные DEK
	block, err := aes.NewCipher(dek)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher for DEK: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM for DEK: %w", err)
	}
	dataNonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(dataNonce); err != nil {
		return nil, fmt.Errorf("error generate dataNonce: %w", err)
	}

	encryptedData := gcm.Seal(nil, dataNonce, data, nil)

	// 3. Шифруем DEK Master Key
	mkBlock, _ := aes.NewCipher(masterKey)
	mkGCM, _ := cipher.NewGCM(mkBlock)
	dekNonce := make([]byte, mkGCM.NonceSize())
	if _, err := rand.Read(dekNonce); err != nil {
		return nil, fmt.Errorf("error generate dekNonce: %w", err)
	}

	encryptedDEK := mkGCM.Seal(nil, dekNonce, dek, nil)

	return &models.EncryptedData{
		EncryptedData: encryptedData,
		DataNonce:     dataNonce,
		EncryptedDek:  encryptedDEK,
		DekNonce:      dekNonce,
	}, nil
}

// DecryptUserData decrypts the encrypted user data using the master key.
// It first decrypts the DEK, then uses it to decrypt the actual data.
func (e *Envelope) DecryptUserData(
	ctx context.Context,
	userData models.DBUserData,
	masterKey []byte,
) ([]byte, error) {
	// 1. Расшифровываем DEK с помощью Master Key
	mkBlock, err := aes.NewCipher(masterKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher for master key: %w", err)
	}

	mkGCM, err := cipher.NewGCM(mkBlock)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM for master key: %w", err)
	}

	dek, err := mkGCM.Open(nil, userData.DekNonce, userData.EncryptedDek, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt DEK: %w", err)
	}

	// 2. Расшифровываем данные с помощью DEK
	block, err := aes.NewCipher(dek)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher for DEK: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM for DEK: %w", err)
	}

	decryptData, err := gcm.Open(nil, userData.DataNonce, userData.EncryptedData, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %w", err)
	}

	return decryptData, nil
}
