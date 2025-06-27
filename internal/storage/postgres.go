// Package storage provides PostgreSQL-backed storage implementation for GophKeeper.
//
// This package defines the Storage type and related interfaces for managing users, master keys,
// and user data in a PostgreSQL database. It handles database migrations, connection pooling,
// and CRUD operations for application data.
package storage

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pressly/goose/v3"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/apetsko/gophkeeper/models"
)

// PgxPoolIface abstracts a subset of pgxpool.Pool methods for database operations.
//
// This interface allows for easier testing and mocking of database interactions.
type PgxPoolIface interface {
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	Begin(context.Context) (pgx.Tx, error)
	Close()
	Ping(context.Context) error
}

// migrations embeds SQL migration files for database schema management.
//
//go:embed migrations/*.sql
var migrations embed.FS

// Storage implements the IStorage interface using a PostgreSQL backend.
//
// It provides methods for user management, master key storage, and user data operations.
type Storage struct {
	DB PgxPoolIface
}

// migrate applies database migrations using goose and the embedded migration files.
//
// Parameters:
//   - conn: PostgreSQL connection string.
//
// Returns:
//   - error: An error if migrations fail, otherwise nil.
func migrate(conn string) error {
	goose.SetBaseFS(migrations)

	var db *sql.DB
	var err error
	for i := 0; i < 10; i++ {
		db, err = sql.Open("pgx", conn)
		if err == nil {
			err = db.Ping()
			if err == nil {
				break // подключились!
			}
		}
		log.Printf("waiting for DB... (%d/10)", i+1)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		return fmt.Errorf("goose: failed to connect to DB after retries: %w", err)
	}
	defer db.Close()

	err = goose.Up(db, "migrations")
	if err != nil {
		return fmt.Errorf("goose: migration failed: %w", err)
	}
	return nil
}

// NewPostgresClient creates a new Storage instance with a PostgreSQL connection pool.
//
// It applies database migrations before establishing the connection.
//
// Parameters:
//   - conn: PostgreSQL connection string.
//
// Returns:
//   - IStorage: The initialized storage instance.
//   - error: An error if migrations or connection fail.
func NewPostgresClient(conn string) (IStorage, error) {
	if err := migrate(conn); err != nil {
		return nil, err
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, conn)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	return &Storage{DB: pool}, nil
}

// Close closes the underlying database connection pool.
//
// Returns:
//   - error: Always nil.
func (p *Storage) Close() error {
	p.DB.Close()
	return nil
}

// AddUser inserts a new user into the database.
//
// Parameters:
//   - ctx: Context for the operation.
//   - u: Pointer to the UserEntry to add.
//
// Returns:
//   - int: The new user's ID.
//   - error: An error if the user exists or insertion fails.
func (p *Storage) AddUser(ctx context.Context, u *models.UserEntry) (int, error) {
	const insertUser = `
        INSERT INTO users (username, password_hash, created_at, updated_at)
        VALUES ($1, $2, NOW(), NOW())
        ON CONFLICT (username) DO NOTHING
        RETURNING id;
    `

	var id int
	err := p.DB.QueryRow(ctx, insertUser, u.Username, u.PasswordHash).Scan(&id)

	switch {
	case err == nil:
		return id, nil // Successfully created
	case errors.Is(err, pgx.ErrNoRows):
		return 0, models.ErrUserExists // Username conflict
	default:
		return 0, fmt.Errorf("failed to insert user: %w", err)
	}
}

// GetUser retrieves a user by username.
//
// Parameters:
//   - ctx: Context for the operation.
//   - username: Username to search for.
//
// Returns:
//   - *models.UserEntry: The found user entry.
//   - error: An error if not found or query fails.
func (p *Storage) GetUser(ctx context.Context, username string) (*models.UserEntry, error) {
	const getUser = `
		SELECT id, username, password_hash FROM users
		WHERE username = $1;
	`

	var u models.UserEntry

	err := p.DB.QueryRow(ctx, getUser, username).Scan(&u.ID, &u.Username, &u.PasswordHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

// SaveMasterKey stores an encrypted master key for a user.
//
// Parameters:
//   - ctx: Context for the operation.
//   - userID: User ID.
//   - encryptedMK: Encrypted master key bytes.
//   - nonce: Nonce used for encryption.
//
// Returns:
//   - int: The new record's ID.
//   - error: An error if the operation fails.
func (p *Storage) SaveMasterKey(
	ctx context.Context,
	userID int,
	encryptedMK []byte,
	nonce []byte,
) (int, error) {
	const insertSQL = `
        INSERT INTO user_keys (user_id, encrypted_master_key, nonce) 
        VALUES ($1, $2, $3)
        RETURNING id;
    `

	var id int

	err := p.DB.QueryRow(ctx, insertSQL, userID, encryptedMK, nonce).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to save master key: %w", err)
	}

	return id, err
}

// GetMasterKey retrieves the encrypted master key for a user.
//
// Parameters:
//   - ctx: Context for the operation.
//   - userID: User ID.
//
// Returns:
//   - *models.EncryptedMK: The encrypted master key and nonce.
//   - error: An error if not found or query fails.
func (p *Storage) GetMasterKey(ctx context.Context, userID int) (*models.EncryptedMK, error) {
	const selectSQL = `
        SELECT encrypted_master_key, nonce FROM user_keys 
        WHERE user_id = $1;
    `

	var encryptedMK models.EncryptedMK
	err := p.DB.QueryRow(ctx, selectSQL, userID).Scan(&encryptedMK.EncryptedMK, &encryptedMK.Nonce)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrMasterKeyNotFound
		}
		return nil, err
	}

	return &encryptedMK, err
}

// SaveUserData stores encrypted user data in the database.
//
// Parameters:
//   - ctx: Context for the operation.
//   - userData: Pointer to the DBUserData to store.
//
// Returns:
//   - int: The new record's ID.
//   - error: An error if the operation fails.
func (p *Storage) SaveUserData(ctx context.Context, userData *models.DBUserData) (int, error) {
	const insertSQL = `
        INSERT INTO user_data (user_id, type, minio_object_id, encrypted_data, data_nonce, encrypted_dek, dek_nonce, meta) 
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING id;
    `

	var id int

	err := p.DB.QueryRow(
		ctx,
		insertSQL,
		userData.UserID,
		userData.Type,
		userData.MinioObjectID,
		userData.EncryptedData,
		userData.DataNonce,
		userData.EncryptedDek,
		userData.DekNonce,
		userData.Meta,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to save user data: %w", err)
	}

	return id, err
}

// GetUserData retrieves a user data record by its ID.
//
// Parameters:
//   - ctx: Context for the operation.
//   - userDataID: ID of the user data record.
//
// Returns:
//   - *models.DBUserData: The user data record.
//   - error: An error if not found or query fails.
func (p *Storage) GetUserData(ctx context.Context, userDataID int) (*models.DBUserData, error) {
	const selectSQL = `
        SELECT user_id, 
               type,
               minio_object_id,
               encrypted_data,
               data_nonce,
               encrypted_dek,
               dek_nonce,
               meta FROM user_data 
        WHERE id = $1;
    `

	var userData models.DBUserData

	err := p.DB.QueryRow(ctx, selectSQL, userDataID).Scan(
		&userData.UserID,
		&userData.Type,
		&userData.MinioObjectID,
		&userData.EncryptedData,
		&userData.DataNonce,
		&userData.EncryptedDek,
		&userData.DekNonce,
		&userData.Meta,
	)
	if err != nil {
		return &userData, fmt.Errorf("failed to get user data: %w", err)
	}

	return &userData, err
}

// GetUserDataList returns a list of user data items for a given user.
//
// Parameters:
//   - ctx: Context for the operation.
//   - userID: User ID.
//
// Returns:
//   - []models.UserDataListItem: List of user data items.
//   - error: An error if the query fails.
func (p *Storage) GetUserDataList(ctx context.Context, userID int) ([]models.UserDataListItem, error) {
	const selectSQL = `
        SELECT id,
               user_id, 
               type,
               meta,
               created_at
        FROM user_data 
        WHERE user_id = $1
        ORDER BY id DESC;
    `

	rows, err := p.DB.Query(ctx, selectSQL, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query user data list: %w", err)
	}
	defer rows.Close()

	var result []models.UserDataListItem
	for rows.Next() {
		var data models.UserDataListItem
		err := rows.Scan(
			&data.ID,
			&data.UserID,
			&data.Type,
			&data.Meta,
			&data.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		result = append(result, data)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return result, nil
}

// DeleteUserData deletes a user data record by its ID.
//
// Parameters:
//   - ctx: Context for the operation.
//   - userDataID: ID of the user data record to delete.
//
// Returns:
//   - error: An error if not found or deletion fails.
func (p *Storage) DeleteUserData(ctx context.Context, userDataID int) error {
	const deleteSQL = `
        DELETE FROM user_data 
        WHERE id = $1
        RETURNING id;
    `

	var deletedID int
	err := p.DB.QueryRow(ctx, deleteSQL, userDataID).Scan(&deletedID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("user data with ID %d not found", userDataID)
		}
		return fmt.Errorf("failed to delete user data: %w", err)
	}

	return nil
}
