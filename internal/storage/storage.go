// Package storage defines the storage interface for GophKeeper.
//
// This package provides the IStorage interface, which abstracts the storage layer for user management,
// master key handling, and encrypted user data operations. Implementations may use different backends
// such as PostgreSQL or in-memory storage.
package storage

import (
	"context"

	"github.com/apetsko/gophkeeper/models"
)

// IStorage defines the contract for persistent storage operations in GophKeeper.
//
// Implementations of this interface provide methods for user management, master key storage,
// and CRUD operations on encrypted user data.
//
//go:generate mockery --name=IStorage --output=../mocks/ --case=underscore
type IStorage interface {
	// Close releases any resources held by the storage implementation.
	Close() error

	// AddUser adds a new user to the storage.
	// Returns the new user's ID or an error if the user already exists or the operation fails.
	AddUser(ctx context.Context, u *models.UserEntry) (int, error)

	// GetUser retrieves a user by username.
	// Returns the user entry or an error if not found.
	GetUser(ctx context.Context, username string) (*models.UserEntry, error)

	// SaveMasterKey stores an encrypted master key for a user.
	// Returns the new record's ID or an error if the operation fails.
	SaveMasterKey(ctx context.Context, userID int, encryptedMK []byte, nonce []byte) (int, error)

	// GetMasterKey retrieves the encrypted master key for a user.
	// Returns the encrypted master key or an error if not found.
	GetMasterKey(ctx context.Context, userID int) (*models.EncryptedMK, error)

	// SaveUserData stores encrypted user data in the storage.
	// Returns the new record's ID or an error if the operation fails.
	SaveUserData(ctx context.Context, userData *models.DBUserData) (int, error)

	// GetUserData retrieves a user data record by its ID.
	// Returns the user data or an error if not found.
	GetUserData(ctx context.Context, userDataID int) (*models.DBUserData, error)

	// GetUserDataList returns a list of user data items for a given user.
	// Returns the list or an error if the query fails.
	GetUserDataList(ctx context.Context, userID int) ([]models.UserDataListItem, error)

	// DeleteUserData deletes a user data record by its ID.
	// Returns an error if not found or deletion fails.
	DeleteUserData(ctx context.Context, userDataID int) error
}
