// Copyright (C) 2018 go-dappley authors
//
// This file is part of the go-dappley library.
//
// the go-dappley library is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either pubKeyHash 3 of the License, or
// (at your option) any later pubKeyHash.
//
// the go-dappley library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with the go-dappley library.  If not, see <http://www.gnu.org/licenses/>.
//

package core

import (
	"crypto/ecdsa"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateAddress(t *testing.T) {
	key1 := &KeyPair{PrivateKey: ecdsa.PrivateKey{},
		PublicKey: []uint8{0x5c, 0x7b, 0x4e, 0x64, 0x19, 0x37, 0xaf, 0x2a, 0x9c, 0x56, 0x84, 0x3, 0x6e, 0x13, 0x3d, 0x92, 0x4, 0x94, 0x32, 0x23, 0xea, 0xe3, 0xcb, 0x6d, 0xf8, 0xb5, 0xf2, 0x92, 0x11, 0x61, 0xd, 0x9, 0xc1, 0x5b, 0x56, 0x17, 0x1d, 0x91, 0xf9, 0x53, 0x76, 0x1a, 0xce, 0x7a, 0x5c, 0xae, 0xe1, 0xc5, 0xa3, 0xbb, 0xcb, 0xd2, 0x5b, 0x6f, 0xf3, 0x4e, 0x1, 0x3b, 0xc1, 0xf8, 0x39, 0xe, 0x90, 0x6}}
	key2 := &KeyPair{PrivateKey: ecdsa.PrivateKey{},
		PublicKey: []uint8{0xff, 0x62, 0x80, 0x2b, 0xec, 0xac, 0x6f, 0x6c, 0x16, 0xda, 0xde, 0x6e, 0xa9, 0x3b, 0x87, 0x8a, 0x17, 0xc7, 0x9c, 0x2e, 0x2e, 0x4c, 0x2f, 0xb9, 0x64, 0xda, 0x12, 0x60, 0x91, 0x82, 0x9a, 0x64, 0x73, 0xd7, 0xd3, 0x4b, 0x51, 0x81, 0x9e, 0xd2, 0x2e, 0xb9, 0x42, 0x1, 0xce, 0xe0, 0x19, 0x97, 0xa0, 0x8e, 0xea, 0x80, 0xb, 0x18, 0x64, 0x8b, 0xf4, 0xd4, 0xd, 0xdc, 0x91, 0x40, 0x37, 0x75}}

	addressExpectForKey1 := Address{Address: "dGDrVKjCG3sdXtDUgWZ7Fp3Q97tLhqWivf"}
	addressExpectForKey2 := Address{Address: "dG6HhzSdA5m7KqvJNszVSf8i5f4neAteSs"}
	address1 := key1.GenerateAddress(false)
	address2 := key2.GenerateAddress(false)

	address3 := key1.GenerateAddress(false)

	assert.NotNil(t, address1)
	assert.NotNil(t, address2)
	assert.NotNil(t, address3)

	assert.Equal(t, address1, addressExpectForKey1)
	assert.Equal(t, address2, addressExpectForKey2)
	assert.Equal(t, address1, address3)
}

func TestNewKeyPair(t *testing.T) {
	key1 := NewKeyPair()

	assert.NotNil(t, key1)
	assert.NotNil(t, key1.PrivateKey)
	assert.NotNil(t, key1.PublicKey)

	assert.Equal(t, 64, len(key1.PublicKey))
	assert.Equal(t, 32, len(key1.PrivateKey.D.Bytes()))

}
