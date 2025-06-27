// Package constants defines common constants and utility functions used across the GophKeeper application.
package constants

import pbc "github.com/apetsko/gophkeeper/protogen/api/proto/v1/common"

// mdKey is a custom type for metadata keys used in context.
type mdKey string

// contextKey is a custom type for context keys.
type contextKey string

const (
	// JWT is the context key for storing JWT tokens in metadata.
	JWT mdKey = "jwt"
	// UserID is the context key for storing user ID in context.
	UserID contextKey = "userID"
)

const (
	// BankCard represents the data type for bank card information.
	BankCard string = "bank_card"
	// Credentials represents the data type for user credentials.
	Credentials string = "credentials"
	// BinaryData represents the data type for binary data.
	BinaryData string = "binary_data"
)

const (
	// KeyLength is the required length (in bytes) for encryption keys.
	KeyLength int = 32
	// Mem is the memory parameter for cryptographic operations (e.g., Argon2).
	Mem uint32 = 64 * 1024
	// Threads is the number of threads for cryptographic operations.
	Threads uint8 = 4
)

// MapDataTypeToString maps a protobuf DataType to its string representation.
//
// Returns the corresponding string constant for the given DataType.
// If the type is unknown, returns "unknown".
func MapDataTypeToString(dt pbc.DataType) string {
	switch dt {
	case pbc.DataType_DATA_TYPE_BANK_CARD:
		return BankCard
	case pbc.DataType_DATA_TYPE_CREDENTIALS:
		return Credentials
	case pbc.DataType_DATA_TYPE_BINARY_DATA:
		return BinaryData
	default:
		return "unknown"
	}
}
