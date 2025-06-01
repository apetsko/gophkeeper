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

func (p *Storage) SaveUserData(ctx context.Context, userData *models.SaveUserData) (int, error) {
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
