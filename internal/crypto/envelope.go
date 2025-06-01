package crypto

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"

	"github.com/apetsko/gophkeeper/models"
)

type EnvelopStorage interface {
	SaveUserData(ctx context.Context, userData *models.SaveUserData) (int, error)
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

	saveUserData := &models.SaveUserData{
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
