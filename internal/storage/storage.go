// storage/interface.go

package storage

import (
	"context"

	"github.com/apetsko/gophkeeper/models"
)

//go:generate mockery --name=IStorage --output=../mocks/ --case=underscore

type IStorage interface {
	Close() error
	AddUser(ctx context.Context, u *models.UserEntry) (int, error)
	GetUser(ctx context.Context, username string) (*models.UserEntry, error)
	SaveMasterKey(ctx context.Context, userID int, encryptedMK []byte, nonce []byte) (int, error)
	GetMasterKey(ctx context.Context, userID int) (*models.EncryptedMK, error)
	SaveUserData(ctx context.Context, userData *models.DBUserData) (int, error)
	GetUserData(ctx context.Context, userDataID int) (*models.DBUserData, error)
	GetUserDataList(ctx context.Context, userID int) ([]models.UserDataListItem, error)
	DeleteUserData(ctx context.Context, userDataID int) error
}
