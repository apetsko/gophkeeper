package crypto

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/subtle"
	"errors"

	"golang.org/x/crypto/argon2"

	"github.com/apetsko/gophkeeper/models"
)

type KeyStorage interface {
	GetMasterKey(ctx context.Context, userID int) (*models.EncryptedMK, error)
	SaveMasterKey(ctx context.Context, userID int, encryptedMK, nonce []byte) (int, error)
}

type KeyManager struct {
	storage             KeyStorage
	serverEncryptionKey []byte
}

func NewKeyManager(
	storage KeyStorage,
	serverEncryptionKey []byte,
) *KeyManager {
	return &KeyManager{
		storage:             storage,
		serverEncryptionKey: serverEncryptionKey,
	}
}

func (m *KeyManager) GetMasterKey(
	ctx context.Context,
	userID int,
	userPassword string,
	userSalt []byte,
) ([]byte, error) {
	// 1. Получаем зашифрованный MK из БД
	encryptedMK, err := m.storage.GetMasterKey(ctx, userID)
	if err != nil {
		if errors.Is(err, models.MasterKeyNotFound) {
			return m.generateAndStoreMasterKey(
				ctx,
				userID,
				userPassword,
				userSalt,
			)
		}

		return nil, err
	}

	// 2. Расшифровываем серверным ключом
	block, _ := aes.NewCipher(m.serverEncryptionKey)
	gcm, _ := cipher.NewGCM(block)

	mk, err := gcm.Open(nil, encryptedMK.Nonce, encryptedMK.EncryptedMK, nil)
	if err != nil {
		return nil, err
	}

	// Верифицируем пароль
	computedMK := argon2.IDKey([]byte(userPassword), userSalt, 3, 64*1024, 4, 32)
	if subtle.ConstantTimeCompare(mk, computedMK) != 1 {
		return nil, errors.New("invalid password")
	}

	return mk, nil
}

func (m *KeyManager) generateAndStoreMasterKey(
	ctx context.Context,
	userID int,
	userPassword string,
	userSalt []byte,
) ([]byte, error) {
	// 1. Генерируем Master Key из пароля пользователя
	mk := argon2.IDKey([]byte(userPassword), userSalt, 3, 64*1024, 4, 32)

	// 2. Шифруем Master Key серверным ключом
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

	// 3. Сохраняем в БД (user_id, encrypted_mk, nonce)
	_, err := m.storage.SaveMasterKey(ctx, userID, encryptedMK, nonce)

	return mk, err
}
