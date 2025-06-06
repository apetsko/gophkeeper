package models

import "time"

// DBUserData represents a user data record stored in the database.
//
// Fields:
//   - UserID: The ID of the user who owns the data.
//   - Type: The type/category of the data.
//   - MinioObjectID: The S3/MinIO object identifier.
//   - Meta: Metadata associated with the data.
//   - EncryptedData: The encrypted data bytes.
//   - DataNonce: Nonce for the encrypted data.
//   - EncryptedDek: The encrypted data encryption key.
//   - DekNonce: Nonce for the encrypted DEK.
type DBUserData struct {
	UserID        int    `json:"user_id"`
	Type          string `json:"type"`
	MinioObjectID string `json:"minio_object_id"`
	Meta          string `json:"meta"`
	EncryptedData []byte `json:"encrypted_data"`
	DataNonce     []byte `json:"data_nonce"`
	EncryptedDek  []byte `json:"encrypted_dek"`
	DekNonce      []byte `json:"dek_nonce"`
}

// UserDataListItem represents a summary of a user data record for listing purposes.
//
// Fields:
//   - ID: Unique identifier of the data record.
//   - UserID: The ID of the user who owns the data.
//   - Type: The type/category of the data.
//   - Meta: Metadata associated with the data.
//   - CreatedAt: Timestamp when the data was created.
type UserDataListItem struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Type      string    `json:"type"`
	Meta      string    `json:"meta"`
	CreatedAt time.Time `json:"created_at"`
}
