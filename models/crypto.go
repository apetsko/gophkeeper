package models

type EncryptedMK struct {
	EncryptedMK []byte `json:"encrypted_mk"`
	Nonce       []byte `json:"nonce"`
}
