package models

import "time"

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

type UserDataListItem struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Type      string    `json:"type"`
	Meta      string    `json:"meta"`
	CreatedAt time.Time `json:"created_at"`
}
