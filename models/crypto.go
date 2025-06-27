package models

// EncryptedMK holds an encrypted master key and its nonce.
//
// Fields:
//   - EncryptedMK: The encrypted master key bytes.
//   - Nonce: The nonce used for encryption.
type EncryptedMK struct {
	EncryptedMK []byte `json:"encrypted_mk"`
	Nonce       []byte `json:"nonce"`
}

// EncryptedData contains encrypted user data and associated nonces.
//
// Fields:
//   - EncryptedData: The encrypted data bytes.
//   - DataNonce: Nonce for the encrypted data.
//   - EncryptedDek: The encrypted data encryption key.
//   - DekNonce: Nonce for the encrypted DEK.
type EncryptedData struct {
	EncryptedData []byte
	DataNonce     []byte
	EncryptedDek  []byte
	DekNonce      []byte
}
