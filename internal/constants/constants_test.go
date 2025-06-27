package constants

import (
	"testing"

	pbc "github.com/apetsko/gophkeeper/protogen/api/proto/v1/common"
	"github.com/stretchr/testify/require"
)

func TestMapDataTypeToString(t *testing.T) {
	tests := []struct {
		name     string
		input    pbc.DataType
		expected string
	}{
		{"bank card", pbc.DataType_DATA_TYPE_BANK_CARD, BankCard},
		{"credentials", pbc.DataType_DATA_TYPE_CREDENTIALS, Credentials},
		{"binary data", pbc.DataType_DATA_TYPE_BINARY_DATA, BinaryData},
		{"unknown", pbc.DataType(999), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MapDataTypeToString(tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}
