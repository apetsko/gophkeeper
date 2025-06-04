package constants

import pbc "github.com/apetsko/gophkeeper/protogen/api/proto/v1/common"

type mdKey string
type contextKey string

const (
	JWT    mdKey      = "jwt"
	UserID contextKey = "userID"
)

const (
	BankCard    string = "bank_card"
	Credentials string = "credentials"
	BinaryData  string = "binary_data"
)

const (
	KeyLength int    = 32
	Mem       uint32 = 64 * 1024
	Threads   uint8  = 4
)

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
