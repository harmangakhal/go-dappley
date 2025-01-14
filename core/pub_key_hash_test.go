package core

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewUserPubKeyHash(t *testing.T) {
	expect := []uint8([]byte{versionUser,0xb1, 0x34, 0x4c, 0x17, 0x67, 0x4c, 0x18, 0xd1, 0xa2, 0xdc, 0xea, 0x9f, 0x17, 0x16, 0xe0, 0x49, 0xf4, 0xa0, 0x5e, 0x6c})

	publicKey := []uint8([]byte{0xd7, 0x23, 0x82, 0x25, 0xaa, 0x81, 0x1f, 0x4d, 0xf6, 0xae, 0x31, 0x35, 0x60, 0xfc, 0x81, 0x7, 0x8, 0x8b, 0x3b, 0x87, 0x25, 0xae, 0xf3, 0xec, 0x62, 0xde, 0xa8, 0x88, 0xbc, 0x1e, 0x93, 0xa4, 0xc9, 0xac, 0xfa, 0x27, 0x83, 0xf4, 0x69, 0x61, 0x57, 0xb5, 0x82, 0xe6, 0x62, 0xd0, 0x18, 0x5c, 0xdd, 0x28, 0xbf, 0xe4, 0x5c, 0xb5, 0xd7, 0xe3, 0xb5, 0x43, 0xd, 0x20, 0xac, 0x73, 0x58, 0x15})
	content, _ := NewUserPubKeyHash(publicKey)
	assert.Equal(t, expect, content.GetPubKeyHash())
}

func TestNewUserPubKeyHash_Fail(t *testing.T) {
	content, err := NewUserPubKeyHash(nil)
	assert.Nil(t, content.GetPubKeyHash())
	assert.NotNil(t, err)
}

func TestNewContractPubKeyHash(t *testing.T) {
	pkh := NewContractPubKeyHash()
	assert.Equal(t, versionContract, pkh.PubKeyHash[0])
}

func TestPubKeyHash_IsContract(t *testing.T) {
	tests := []struct {
		name        string
		pubKeyHash  []byte
		expectedRes bool
		expectedErr error
	}{
		{
			name:        "ContractAddress",
			pubKeyHash:  []byte{versionContract},
			expectedRes: true,
			expectedErr: nil,
		},
		{
			name:        "UserAddress",
			pubKeyHash:  []byte{versionUser},
			expectedRes: false,
			expectedErr: nil,
		},
		{
			name:        "InvalidAddress",
			pubKeyHash:  []byte{0x00},
			expectedRes: false,
			expectedErr: ErrInvalidPubKeyHashVersion,
		},
		{
			name:        "EmptyAddress",
			pubKeyHash:  []byte{},
			expectedRes: false,
			expectedErr: ErrEmptyPublicKeyHash,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkh := PubKeyHash{tt.pubKeyHash}
			res,err := pkh.IsContract()
			assert.Equal(t, res, tt.expectedRes)
			assert.Equal(t, err, tt.expectedErr)
		})
	}
}