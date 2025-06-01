package models

type BankCard struct {
	CardNumber string `json:"card_number"`
	ExpiryDate string `json:"expiry_date"`
	Cvv        string `json:"cvv"`
	Cardholder string `json:"cardholder"`
}

type UserData struct {
	UserID        int    `json:"user_id"`
	Type          string `json:"type"`
	MinioObjectID string `json:"minio_object_id"`
	Meta          string `json:"meta"`
}

type SaveUserData struct {
	UserID        int    `json:"user_id"`
	Type          string `json:"type"`
	MinioObjectID string `json:"minio_object_id"`
	EncryptedData []byte `json:"encrypted_data"`
	DataNonce     []byte `json:"data_nonce"`
	EncryptedDek  []byte `json:"encrypted_dek"`
	DekNonce      []byte `json:"dek_nonce"`
	Meta          string `json:"meta"`
}
