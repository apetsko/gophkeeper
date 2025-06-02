package storage

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pressly/goose/v3"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/apetsko/gophkeeper/models"
)

type PgxPoolIface interface {
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	Begin(context.Context) (pgx.Tx, error)
	Close()
}

//go:embed migrations/*.sql
var migrations embed.FS

type Storage struct {
	DB PgxPoolIface
}

func migrate(conn string) error {
	goose.SetBaseFS(migrations)
	db, err := sql.Open("pgx", conn)
	if err != nil {
		return fmt.Errorf("goose: failed to open DB: %w", err)
	}
	defer db.Close()

	err = goose.Up(db, "migrations")
	if err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}

func NewPostrgesClient(conn string) (*Storage, error) {
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

func (p *Storage) Close() error {
	p.DB.Close()
	return nil
}

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
		return id, nil // Успешное создание
	case errors.Is(err, pgx.ErrNoRows):
		return 0, models.ErrUserExists // Конфликт по username
	default:
		return 0, fmt.Errorf("failed to insertUser user: %w", err)
	}
}

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

func (p *Storage) GetMasterKey(ctx context.Context, userID int) (*models.EncryptedMK, error) {
	const selectSQL = `
        SELECT encrypted_master_key, nonce FROM user_keys 
        WHERE user_id = $1;
    `

	var encryptedMK models.EncryptedMK
	err := p.DB.QueryRow(ctx, selectSQL, userID).Scan(&encryptedMK.EncryptedMK, &encryptedMK.Nonce)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.MasterKeyNotFound
		}
		return nil, err
	}

	return &encryptedMK, err
}

func (p *Storage) SaveUserData(ctx context.Context, userData *models.DbUserData) (int, error) {
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

func (p *Storage) GetUserData(ctx context.Context, userDataID int) (*models.DbUserData, error) {
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

	var userData models.DbUserData

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
