package models

type EncryptedMK struct {
	EncryptedMK []byte `json:"encrypted_mk"`
	Nonce       []byte `json:"nonce"`
}

type EncryptedData struct {
	EncryptedData []byte
	DataNonce     []byte
	EncryptedDek  []byte
	DekNonce      []byte
}
