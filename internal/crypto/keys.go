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

// Интерфейс для удобства мокирования и внедрения зависимостей
//
//go:generate mockery --dir ./internal/crypto --name=KeyManagerInterface --output=../mocks/ --case=underscore
type KeyManagerInterface interface {
	GetMasterKey(ctx context.Context, userID int) ([]byte, error)
	GetOrCreateMasterKey(ctx context.Context, userID int, userPassword string, userSalt []byte) ([]byte, error)
}

// Убедимся, что KeyManager реализует интерфейс
var _ KeyManagerInterface = (*KeyManager)(nil)

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
