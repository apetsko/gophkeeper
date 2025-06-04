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

//go:generate mockery --dir ./internal/crypto --name=IEnvelopeStorage --output=.../mocks/ --case=underscore
type IEnvelope interface {
	EncryptUserData(ctx context.Context, masterKey []byte, data []byte) (*models.EncryptedData, error)
	DecryptUserData(ctx context.Context, userData models.DBUserData, masterKey []byte) ([]byte, error)
}

type EnvelopStorage interface {
	SaveUserData(ctx context.Context, userData *models.DBUserData) (int, error)
}

type Envelope struct {
	storage EnvelopStorage
}

func (e *Envelope) SaveUserData(ctx context.Context, userData *models.DBUserData) (int, error) {
	//TODO implement me
	panic("implement me")
}

func NewEnvelope(
	storage EnvelopStorage,
) *Envelope {
	return &Envelope{
		storage: storage,
	}
}

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
