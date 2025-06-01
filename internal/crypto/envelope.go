package crypto

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"

	"github.com/apetsko/gophkeeper/models"
)

type EnvelopStorage interface {
	SaveUserData(ctx context.Context, userData *models.DbUserData) (int, error)
}

type Envelop struct {
	storage EnvelopStorage
}

func NewEnvelop(
	storage EnvelopStorage,
) *Envelop {
	return &Envelop{
		storage: storage,
	}
}

func (e *Envelop) EncryptUserData(
	ctx context.Context,
	userData models.UserData,
	masterKey []byte,
	data []byte,
) ([]byte, error) {
	// 1. Генерируем случайный DEK
	dek := make([]byte, 32)
	if _, err := rand.Read(dek); err != nil {
		return nil, errors.New("error generate DEK")
	}

	// 2. Шифруем данные DEK
	block, _ := aes.NewCipher(dek)
	gcm, _ := cipher.NewGCM(block)
	dataNonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(dataNonce); err != nil {
		return nil, errors.New("error generate dataNonce")
	}

	encryptedData := gcm.Seal(nil, dataNonce, data, nil)

	// 3. Шифруем DEK Master Key
	mkBlock, _ := aes.NewCipher(masterKey)
	mkGCM, _ := cipher.NewGCM(mkBlock)
	dekNonce := make([]byte, mkGCM.NonceSize())
	if _, err := rand.Read(dekNonce); err != nil {
		return nil, errors.New("error generate dekNonce")
	}

	encryptedDEK := mkGCM.Seal(nil, dekNonce, dek, nil)

	saveUserData := &models.DbUserData{
		UserID:        userData.UserID,
		Type:          userData.Type,
		MinioObjectID: userData.MinioObjectID,
		EncryptedData: encryptedData,
		DataNonce:     dataNonce,
		EncryptedDek:  encryptedDEK,
		DekNonce:      dekNonce,
		Meta:          userData.Meta,
	}

	_, err := e.storage.SaveUserData(ctx, saveUserData)
	if err != nil {
		return nil, err
	}

	return encryptedData, nil
}

func (e *Envelop) DecryptUserData(
	ctx context.Context,
	userData models.DbUserData,
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
