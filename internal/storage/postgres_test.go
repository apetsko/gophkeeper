package storage

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/apetsko/gophkeeper/models"
	"github.com/apetsko/gophkeeper/pkg/logging"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	logger      = logging.NewLogger(slog.LevelDebug)
	connStr     string
	pgContainer *postgres.PostgresContainer
)

func startTestDB() {
	ctx := context.Background()
	var err error
	pgContainer, err = postgres.Run(ctx,
		"postgres:15.3-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(15*time.Second)),
	)
	if err != nil {
		logger.Error("failed to start postgres container", slog.Any("err", err))
		os.Exit(1)
	}

	connStr, err = pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		logger.Error("failed to get connection string", slog.Any("err", err))
		os.Exit(1)
	}
	logger.Info("✅ PostgreSQL test container started")
}

func stopTestDB() {
	if pgContainer != nil {
		logger.Info("🛑 Stopping test database container...")
		_ = pgContainer.Terminate(context.Background())
	}
}

func setupTestStorage(t *testing.T) IStorage {
	storage, err := NewPostgresClient(connStr)
	require.NoError(t, err)
	t.Cleanup(func() { _ = storage.Close() })
	return storage
}

func TestMain(m *testing.M) {
	startTestDB()
	_, terminateMinio, err := startTestMinio(nil)
	if err != nil {
		logger.Error("failed to start minio container", slog.Any("err", err))
		stopTestDB()
		os.Exit(1)
	}
	time.Sleep(5 * time.Second)
	code := m.Run()
	stopTestDB()
	terminateMinio()
	os.Exit(code)
}

func TestStorage_AddUser_and_GetUser(t *testing.T) {
	st := setupTestStorage(t)
	ctx := context.Background()

	user := &models.UserEntry{
		Username:     "testuser",
		PasswordHash: "hashedpassword",
	}
	id, err := st.AddUser(ctx, user)
	require.NoError(t, err)
	require.NotZero(t, id)

	got, err := st.GetUser(ctx, user.Username)
	require.NoError(t, err)
	require.Equal(t, user.Username, got.Username)
	require.Equal(t, user.PasswordHash, got.PasswordHash)
}

func TestStorage_AddUser_Conflict(t *testing.T) {
	st := setupTestStorage(t)
	ctx := context.Background()

	user := &models.UserEntry{
		Username:     "conflictuser",
		PasswordHash: "hash",
	}
	id1, err := st.AddUser(ctx, user)
	require.NoError(t, err)
	require.NotZero(t, id1)

	id2, err := st.AddUser(ctx, user)
	require.ErrorIs(t, err, models.ErrUserExists)
	require.Zero(t, id2)
}

func TestStorage_GetUser_NotFound(t *testing.T) {
	st := setupTestStorage(t)
	ctx := context.Background()

	_, err := st.GetUser(ctx, "ghostuser")
	require.ErrorIs(t, err, models.ErrUserNotFound)
}

func TestStorage_SaveMasterKey_and_GetMasterKey(t *testing.T) {
	st := setupTestStorage(t)
	ctx := context.Background()

	user := &models.UserEntry{
		Username:     "mkuser",
		PasswordHash: "hash",
	}
	uid, err := st.AddUser(ctx, user)
	require.NoError(t, err)

	encMK := []byte("encrypted-mk")
	nonce := []byte("nonce")
	mkID, err := st.SaveMasterKey(ctx, uid, encMK, nonce)
	require.NoError(t, err)
	require.NotZero(t, mkID)

	got, err := st.GetMasterKey(ctx, uid)
	require.NoError(t, err)
	require.Equal(t, encMK, got.EncryptedMK)
	require.Equal(t, nonce, got.Nonce)
}

func TestStorage_GetMasterKey_NotFound(t *testing.T) {
	st := setupTestStorage(t)
	ctx := context.Background()

	_, err := st.GetMasterKey(ctx, 999999)
	require.ErrorIs(t, err, models.ErrMasterKeyNotFound)
}

func TestStorage_GetUserData_NotFound(t *testing.T) {
	st := setupTestStorage(t)
	ctx := context.Background()

	_, err := st.GetUserData(ctx, 999999)
	require.Error(t, err)
}

func TestStorage_DeleteUserData_NotFound(t *testing.T) {
	st := setupTestStorage(t)
	ctx := context.Background()

	err := st.DeleteUserData(ctx, 999999)
	require.Error(t, err)
}

func TestStorage_ContextCancelled(t *testing.T) {
	st := setupTestStorage(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := st.GetUser(ctx, "any")
	require.Error(t, err)
}

func TestStorage_Ping(t *testing.T) {
	st := setupTestStorage(t)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	require.NoError(t, st.(*Storage).DB.Ping(ctx))
}

func TestNewPostgresClient_BadConnString(t *testing.T) {
	_, err := NewPostgresClient("invalid-conn-string")
	require.Error(t, err)
}

func TestStorage_AddUser_DBError(t *testing.T) {
	st := setupTestStorage(t)
	st.(*Storage).DB.Close() // Close pool to force error
	_, err := st.AddUser(context.Background(), &models.UserEntry{Username: "x", PasswordHash: "y"})
	require.Error(t, err)
}

func TestStorage_SaveMasterKey_DBError(t *testing.T) {
	st := setupTestStorage(t)
	st.(*Storage).DB.Close()
	_, err := st.SaveMasterKey(context.Background(), 1, []byte("a"), []byte("b"))
	require.Error(t, err)
}

func TestStorage_GetMasterKey_DBError(t *testing.T) {
	st := setupTestStorage(t)
	st.(*Storage).DB.Close()
	_, err := st.GetMasterKey(context.Background(), 1)
	require.Error(t, err)
}

func TestStorage_SaveUserData_DBError(t *testing.T) {
	st := setupTestStorage(t)
	st.(*Storage).DB.Close()
	_, err := st.SaveUserData(context.Background(), &models.DBUserData{})
	require.Error(t, err)
}

func TestStorage_GetUserData_DBError(t *testing.T) {
	st := setupTestStorage(t)
	st.(*Storage).DB.Close()
	_, err := st.GetUserData(context.Background(), 1)
	require.Error(t, err)
}

func TestStorage_GetUserDataList_DBError(t *testing.T) {
	st := setupTestStorage(t)
	st.(*Storage).DB.Close()
	_, err := st.GetUserDataList(context.Background(), 1)
	require.Error(t, err)
}

func TestStorage_DeleteUserData_DBError(t *testing.T) {
	st := setupTestStorage(t)
	st.(*Storage).DB.Close()
	err := st.DeleteUserData(context.Background(), 1)
	require.Error(t, err)
}
