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
	"bytes"
	"crypto/ecdsa"
	"encoding/binary"
	"testing"

	"github.com/dappley/go-dappley/common"
	"github.com/dappley/go-dappley/core/pb"
	"github.com/dappley/go-dappley/crypto/keystore/secp256k1"
	"github.com/dappley/go-dappley/util"
	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"sync"
)

func getAoB(length int64) []byte {
	return util.GenerateRandomAoB(length)
}

func GenerateFakeTxInputs() []TXInput {
	return []TXInput{
		{getAoB(2), 10, getAoB(2), getAoB(2)},
		{getAoB(2), 5, getAoB(2), getAoB(2)},
	}
}

func GenerateFakeTxOutputs() []TXOutput {
	return []TXOutput{
		{common.NewAmount(1), PubKeyHash{getAoB(2)}, ""},
		{common.NewAmount(2), PubKeyHash{getAoB(2)}, ""},
	}
}

func TestTrimmedCopy(t *testing.T) {
	var tx1 = Transaction{
		ID:   util.GenerateRandomAoB(1),
		Vin:  GenerateFakeTxInputs(),
		Vout: GenerateFakeTxOutputs(),
		Tip:  2,
	}

	t2 := tx1.TrimmedCopy()

	t3 := NewCoinbaseTX(NewAddress("13ZRUc4Ho3oK3Cw56PhE5rmaum9VBeAn5F"), "", 0, common.NewAmount(0))
	t4 := t3.TrimmedCopy()
	assert.Equal(t, tx1.ID, t2.ID)
	assert.Equal(t, tx1.Tip, t2.Tip)
	assert.Equal(t, tx1.Vout, t2.Vout)
	for index, vin := range t2.Vin {
		assert.Nil(t, vin.Signature)
		assert.Nil(t, vin.PubKey)
		assert.Equal(t, tx1.Vin[index].Txid, vin.Txid)
		assert.Equal(t, tx1.Vin[index].Vout, vin.Vout)
	}

	assert.Equal(t, t3.ID, t4.ID)
	assert.Equal(t, t3.Tip, t4.Tip)
	assert.Equal(t, t3.Vout, t4.Vout)
	for index, vin := range t4.Vin {
		assert.Nil(t, vin.Signature)
		assert.Nil(t, vin.PubKey)
		assert.Equal(t, t3.Vin[index].Txid, vin.Txid)
		assert.Equal(t, t3.Vin[index].Vout, vin.Vout)
	}
}

func TestSign(t *testing.T) {
	// Fake a key pair
	privKey, _ := ecdsa.GenerateKey(secp256k1.S256(), bytes.NewReader([]byte("fakefakefakefakefakefakefakefakefakefake")))
	ecdsaPubKey, _ := secp256k1.FromECDSAPublicKey(&privKey.PublicKey)
	pubKey := append(privKey.PublicKey.X.Bytes(), privKey.PublicKey.Y.Bytes()...)
	pubKeyHash, _ := NewUserPubKeyHash(pubKey)

	// Previous transactions containing UTXO of the Address
	prevTXs := []*UTXO{
		{TXOutput{common.NewAmount(13), pubKeyHash, ""}, []byte("01"), 0},
		{TXOutput{common.NewAmount(13), pubKeyHash, ""}, []byte("02"), 0},
		{TXOutput{common.NewAmount(13), pubKeyHash, ""}, []byte("03"), 0},
	}

	// New transaction to be signed (paid from the fake account)
	txin := []TXInput{
		{[]byte{1}, 0, nil, pubKey},
		{[]byte{3}, 0, nil, pubKey},
		{[]byte{3}, 2, nil, pubKey},
	}
	txout := []TXOutput{
		{common.NewAmount(19), pubKeyHash, ""},
	}
	tx := Transaction{nil, txin, txout, 0}

	// Sign the transaction
	err := tx.Sign(*privKey, prevTXs)
	if assert.Nil(t, err) {
		// Assert that the signatures were created by the fake key pair
		for i, vin := range tx.Vin {

			if assert.NotNil(t, vin.Signature) {
				txCopy := tx.TrimmedCopy()
				txCopy.Vin[i].Signature = nil
				txCopy.Vin[i].PubKey = pubKeyHash.GetPubKeyHash()

				verified, err := secp256k1.Verify(txCopy.Hash(), vin.Signature, ecdsaPubKey)
				assert.Nil(t, err)
				assert.True(t, verified)
			}
		}
	}
}

func TestVerifyCoinbaseTransaction(t *testing.T) {
	var prevTXs = map[string]Transaction{}

	var tx1 = Transaction{
		ID:   util.GenerateRandomAoB(1),
		Vin:  GenerateFakeTxInputs(),
		Vout: GenerateFakeTxOutputs(),
		Tip:  2,
	}

	var tx2 = Transaction{
		ID:   util.GenerateRandomAoB(1),
		Vin:  GenerateFakeTxInputs(),
		Vout: GenerateFakeTxOutputs(),
		Tip:  5,
	}
	var tx3 = Transaction{
		ID:   util.GenerateRandomAoB(1),
		Vin:  GenerateFakeTxInputs(),
		Vout: GenerateFakeTxOutputs(),
		Tip:  10,
	}
	var tx4 = Transaction{
		ID:   util.GenerateRandomAoB(1),
		Vin:  GenerateFakeTxInputs(),
		Vout: GenerateFakeTxOutputs(),
		Tip:  20,
	}
	prevTXs[string(tx1.ID)] = tx2
	prevTXs[string(tx2.ID)] = tx3
	prevTXs[string(tx3.ID)] = tx4

	// test verifying coinbase transactions
	var t5 = NewCoinbaseTX(NewAddress("13ZRUc4Ho3oK3Cw56PhE5rmaum9VBeAn5F"), "", 5, common.NewAmount(0))
	bh1 := make([]byte, 8)
	binary.BigEndian.PutUint64(bh1, 5)
	txin1 := TXInput{nil, -1, bh1, []byte(nil)}
	txout1 := NewTXOutput(common.NewAmount(10), NewAddress("13ZRUc4Ho3oK3Cw56PhE5rmaum9VBeAn5F"))
	var t6 = Transaction{nil, []TXInput{txin1}, []TXOutput{*txout1}, 0}

	// test valid coinbase transaction
	assert.True(t, t5.Verify(UTXOIndex{}, NewTransactionPool(128), 5))
	assert.True(t, t6.Verify(UTXOIndex{}, NewTransactionPool(128), 5))

	// test coinbase transaction with incorrect blockHeight
	assert.False(t, t5.Verify(UTXOIndex{}, NewTransactionPool(128), 10))

	// test coinbase transaction with incorrect subsidy
	bh2 := make([]byte, 8)
	binary.BigEndian.PutUint64(bh2, 5)
	txin2 := TXInput{nil, -1, bh2, []byte(nil)}
	txout2 := NewTXOutput(common.NewAmount(20), NewAddress("13ZRUc4Ho3oK3Cw56PhE5rmaum9VBeAn5F"))
	var t7 = Transaction{nil, []TXInput{txin2}, []TXOutput{*txout2}, 0}
	assert.False(t, t7.Verify(UTXOIndex{}, NewTransactionPool(128), 5))

}

func TestVerifyNoCoinbaseTransaction(t *testing.T) {
	// Fake a key pair
	privKey, _ := ecdsa.GenerateKey(secp256k1.S256(), bytes.NewReader([]byte("fakefakefakefakefakefakefakefakefakefake")))
	privKeyByte, _ := secp256k1.FromECDSAPrivateKey(privKey)
	pubKey := append(privKey.PublicKey.X.Bytes(), privKey.PublicKey.Y.Bytes()...)
	pubKeyHash, _ := NewUserPubKeyHash(pubKey)
	//Address := KeyPair{*privKey, pubKey}.GenerateAddress()

	// Fake a wrong key pair
	wrongPrivKey, _ := ecdsa.GenerateKey(secp256k1.S256(), bytes.NewReader([]byte("FAKEfakefakefakefakefakefakefakefakefake")))
	wrongPrivKeyByte, _ := secp256k1.FromECDSAPrivateKey(wrongPrivKey)
	wrongPubKey := append(wrongPrivKey.PublicKey.X.Bytes(), wrongPrivKey.PublicKey.Y.Bytes()...)
	//wrongPubKeyHash, _ := NewUserPubKeyHash(wrongPubKey)
	//wrongAddress := KeyPair{*wrongPrivKey, wrongPubKey}.GenerateAddress()
	utxoIndex := NewUTXOIndex()
	utxoIndex.index = map[string][]*UTXO{
		string(pubKeyHash.GetPubKeyHash()): {
			{TXOutput{common.NewAmount(4), pubKeyHash, ""}, []byte{1}, 0},
			{TXOutput{common.NewAmount(3), pubKeyHash, ""}, []byte{2}, 1},
		},
	}

	// Prepare a transaction to be verified
	txin := []TXInput{{[]byte{1}, 0, nil, pubKey}}
	txin1 := append(txin, TXInput{[]byte{2}, 1, nil, pubKey})      // Normal test
	txin2 := append(txin, TXInput{[]byte{2}, 1, nil, wrongPubKey}) // previous not found with wrong pubkey
	txin3 := append(txin, TXInput{[]byte{3}, 1, nil, pubKey})      // previous not found with wrong Txid
	txin4 := append(txin, TXInput{[]byte{2}, 2, nil, pubKey})      // previous not found with wrong TxIndex
	pbh, _ := NewUserPubKeyHash(pubKey)
	txout := []TXOutput{{common.NewAmount(7), pbh, ""}}
	txout2 := []TXOutput{{common.NewAmount(8), pbh, ""}} //Vout amount > Vin amount

	tests := []struct {
		name     string
		transaction       Transaction
		signWith []byte
		ok       bool
	}{
		{"normal", Transaction{nil, txin1, txout, 0}, privKeyByte, true},
		{"previous tx not found with wrong pubkey", Transaction{nil, txin2, txout, 0}, privKeyByte, false},
		{"previous tx not found with wrong Txid", Transaction{nil, txin3, txout, 0}, privKeyByte, false},
		{"previous tx not found with wrong TxIndex", Transaction{nil, txin4, txout, 0}, privKeyByte, false},
		{"Amount invalid", Transaction{nil, txin1, txout2, 0}, privKeyByte, false},
		{"Sign invalid", Transaction{nil, txin1, txout, 0}, wrongPrivKeyByte, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Generate signatures for all tx inputs
			for i := range tt.transaction.Vin {
				txCopy := tt.transaction.TrimmedCopy()
				txCopy.Vin[i].Signature = nil
				txCopy.Vin[i].PubKey = pubKeyHash.GetPubKeyHash()
				signature, _ := secp256k1.Sign(txCopy.Hash(), tt.signWith)
				tt.transaction.Vin[i].Signature = signature
				tt.transaction.ID = tt.transaction.Hash()
			}

			// Verify the signatures
			result := tt.transaction.Verify(*utxoIndex, NewTransactionPool(128), 0)
			assert.Equal(t, tt.ok, result)
		})
	}
}

func TestNewCoinbaseTX(t *testing.T) {
	t1 := NewCoinbaseTX(NewAddress("dXnq2R6SzRNUt7ZANAqyZc2P9ziF6vYekB"), "", 0, common.NewAmount(0))
	expectVin := TXInput{nil, -1, []byte{0, 0, 0, 0, 0, 0, 0, 0}, []byte("Reward to 'dXnq2R6SzRNUt7ZANAqyZc2P9ziF6vYekB'")}
	expectVout := TXOutput{common.NewAmount(10), PubKeyHash{[]byte{0x5a, 0xc9, 0x85, 0x37, 0x92, 0x37, 0x76, 0x80, 0xb1, 0x31, 0xa1, 0xab, 0xb, 0x5b, 0xa6, 0x49, 0xe5, 0x27, 0xf0, 0x42, 0x5d}}, ""}
	assert.Equal(t, 1, len(t1.Vin))
	assert.Equal(t, expectVin, t1.Vin[0])
	assert.Equal(t, 1, len(t1.Vout))
	assert.Equal(t, expectVout, t1.Vout[0])
	assert.Equal(t, uint64(0), t1.Tip)

	t2 := NewCoinbaseTX(NewAddress("dXnq2R6SzRNUt7ZANAqyZc2P9ziF6vYekB"), "", 0, common.NewAmount(0))

	// Assert that NewCoinbaseTX is deterministic (i.e. >1 coinbaseTXs in a block would have identical txid)
	assert.Equal(t, t1, t2)

	t3 := NewCoinbaseTX(NewAddress("dXnq2R6SzRNUt7ZANAqyZc2P9ziF6vYekB"), "", 1, common.NewAmount(0))

	assert.NotEqual(t, t1, t3)
	assert.NotEqual(t, t1.ID, t3.ID)
}

//test IsCoinBase function
func TestIsCoinBase(t *testing.T) {
	var tx1 = Transaction{
		ID:   util.GenerateRandomAoB(1),
		Vin:  GenerateFakeTxInputs(),
		Vout: GenerateFakeTxOutputs(),
		Tip:  2,
	}

	assert.False(t, tx1.IsCoinbase())

	t2 := NewCoinbaseTX(NewAddress("13ZRUc4Ho3oK3Cw56PhE5rmaum9VBeAn5F"), "", 0, common.NewAmount(0))

	assert.True(t, t2.IsCoinbase())

}

func TestTransaction_Proto(t *testing.T) {
	tx1 := Transaction{
		ID:   util.GenerateRandomAoB(1),
		Vin:  GenerateFakeTxInputs(),
		Vout: GenerateFakeTxOutputs(),
		Tip:  5,
	}

	pb := tx1.ToProto()
	var i interface{} = pb
	_, correct := i.(proto.Message)
	assert.Equal(t, true, correct)
	mpb, err := proto.Marshal(pb)
	assert.Nil(t, err)

	newpb := &corepb.Transaction{}
	err = proto.Unmarshal(mpb, newpb)
	assert.Nil(t, err)

	tx2 := Transaction{}
	tx2.FromProto(newpb)

	assert.Equal(t, tx1, tx2)
}

func TestTransaction_GetContractAddress(t *testing.T) {

	tests := []struct {
		name        string
		addr        string
		expectedRes string
	}{
		{
			name:        "ContainsContractAddress",
			addr:        "cavQdWxvUQU1HhBg1d7zJFwhf31SUaQwop",
			expectedRes: "cavQdWxvUQU1HhBg1d7zJFwhf31SUaQwop",
		},
		{
			name:        "ContainsUserAddress",
			addr:        "dGDrVKjCG3sdXtDUgWZ7Fp3Q97tLhqWivf",
			expectedRes: "",
		},
		{
			name:        "EmptyInput",
			addr:        "",
			expectedRes: "",
		},
		{
			name:        "InvalidAddress",
			addr:        "dsdGDrVKjCG3sdXtDUgWZ7Fp3Q97tLhqWivf",
			expectedRes: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addr := NewAddress(tt.addr)
			pkh, _ := addr.GetPubKeyHash()
			tx := Transaction{
				nil,
				nil,
				[]TXOutput{
					{nil,
						PubKeyHash{pkh},
						"",
					},
				},
				0,
			}

			assert.Equal(t, NewAddress(tt.expectedRes), tx.GetContractAddress())
		})
	}
}

func TestTransaction_VerifyDependentTransactions(t *testing.T) {
	var prikey1 = "bb23d2ff19f5b16955e8a24dca34dd520980fe3bddca2b3e1b56663f0ec1aa71"
	var pubkey1 = GetKeyPairByString(prikey1).PublicKey
	var pkHash1, _ = NewUserPubKeyHash(pubkey1)
	var prikey2 = "bb23d2ff19f5b16955e8a24dca34dd520980fe3bddca2b3e1b56663f0ec1aa72"
	var pubkey2 = GetKeyPairByString(prikey2).PublicKey
	var pkHash2, _ = NewUserPubKeyHash(pubkey2)
	var prikey3 = "bb23d2ff19f5b16955e8a24dca34dd520980fe3bddca2b3e1b56663f0ec1aa73"
	var pubkey3 = GetKeyPairByString(prikey3).PublicKey
	var pkHash3, _ = NewUserPubKeyHash(pubkey3)
	var prikey4 = "bb23d2ff19f5b16955e8a24dca34dd520980fe3bddca2b3e1b56663f0ec1aa74"
	var pubkey4 = GetKeyPairByString(prikey4).PublicKey
	var pkHash4, _ = NewUserPubKeyHash(pubkey4)
	var prikey5 = "bb23d2ff19f5b16955e8a24dca34dd520980fe3bddca2b3e1b56663f0ec1aa75"
	var pubkey5 = GetKeyPairByString(prikey5).PublicKey
	var pkHash5, _ = NewUserPubKeyHash(pubkey5)

	var dependentTx1 = Transaction{
		ID: nil,
		Vin: []TXInput{
			{tx1.ID, 1, nil, pubkey1},
		},
		Vout: []TXOutput{
			{common.NewAmount(5), pkHash1,""},
			{common.NewAmount(10), pkHash2,""},
		},
		Tip: 3,
	}
	dependentTx1.ID = dependentTx1.Hash()

	var dependentTx2 = Transaction{
		ID: nil,
		Vin: []TXInput{
			{dependentTx1.ID, 1, nil, pubkey2},
		},
		Vout: []TXOutput{
			{common.NewAmount(5), pkHash3,""},
			{common.NewAmount(3), pkHash4,""},
		},
		Tip: 2,
	}
	dependentTx2.ID = dependentTx2.Hash()

	var dependentTx3 = Transaction{
		ID: nil,
		Vin: []TXInput{
			{dependentTx2.ID, 0, nil, pubkey3},
		},
		Vout: []TXOutput{
			{common.NewAmount(1), pkHash4,""},
		},
		Tip: 4,
	}
	dependentTx3.ID = dependentTx3.Hash()

	var dependentTx4 = Transaction{
		ID: nil,
		Vin: []TXInput{
			{dependentTx2.ID, 1, nil, pubkey4},
			{dependentTx3.ID, 0, nil, pubkey4},
		},
		Vout: []TXOutput{
			{common.NewAmount(3), pkHash1,""},
		},
		Tip: 1,
	}
	dependentTx4.ID = dependentTx4.Hash()

	var dependentTx5 = Transaction{
		ID: nil,
		Vin: []TXInput{
			{dependentTx1.ID, 0, nil, pubkey1},
			{dependentTx4.ID, 0, nil, pubkey1},
		},
		Vout: []TXOutput{
			{common.NewAmount(4), pkHash5,""},
		},
		Tip: 4,
	}
	dependentTx5.ID = dependentTx5.Hash()

	var UtxoIndex = UTXOIndex{
		map[string][]*UTXO{
			string(pkHash2.GetPubKeyHash()): {&UTXO{dependentTx1.Vout[1], dependentTx1.ID, 1}},
			string(pkHash1.GetPubKeyHash()): {&UTXO{dependentTx1.Vout[0], dependentTx1.ID, 0}},
		},
		&sync.RWMutex{},
	}

	tx2Utxo1 := UTXO{dependentTx2.Vout[0], dependentTx2.ID, 0}
	tx2Utxo2 := UTXO{dependentTx2.Vout[1], dependentTx2.ID, 1}
	tx2Utxo3 := UTXO{dependentTx3.Vout[0], dependentTx3.ID, 0}
	tx2Utxo4 := UTXO{dependentTx1.Vout[0], dependentTx1.ID, 0}
	tx2Utxo5 := UTXO{dependentTx4.Vout[0], dependentTx4.ID, 0}
	dependentTx2.Sign(GetKeyPairByString(prikey2).PrivateKey, UtxoIndex.index[string(pkHash2.GetPubKeyHash())])
	dependentTx3.Sign(GetKeyPairByString(prikey3).PrivateKey, []*UTXO{&tx2Utxo1})
	dependentTx4.Sign(GetKeyPairByString(prikey4).PrivateKey, []*UTXO{&tx2Utxo2, &tx2Utxo3})
	dependentTx5.Sign(GetKeyPairByString(prikey1).PrivateKey, []*UTXO{&tx2Utxo4, &tx2Utxo5})

	txPool := NewTransactionPool(6)
	txPool.Push(dependentTx2)
	txPool.Push(dependentTx3)
	txPool.Push(dependentTx4)
	txPool.Push(dependentTx5)

	// verify dependent txs 2,3,4,5 with relation:
	//tx1 (UtxoIndex)
	//|     \
	//tx2    \
	//|  \    \
	//tx3-tx4-tx5

	// test a transaction whose Vin is from UtxoIndex
	assert.Equal(t, true, dependentTx2.Verify(UtxoIndex, txPool, 0))

	// test a transaction whose Vin is from another transaction in transaction pool
	assert.Equal(t, true, dependentTx3.Verify(UtxoIndex, txPool, 0))

	// test a transaction whose Vin is from another two transactions in transaction pool
	assert.Equal(t, true, dependentTx4.Verify(UtxoIndex, txPool, 0))

	// test a transaction whose Vin is from another transaction in transaction pool and UtxoIndex
	assert.Equal(t, true, dependentTx5.Verify(UtxoIndex, txPool, 0))

	// test UTXOs not found for parent transactions
	assert.Equal(t, false, dependentTx3.Verify(UTXOIndex{make(map[string][]*UTXO), &sync.RWMutex{}}, txPool, 0))

	// test a standalone transaction
	txPool.Push(tx1)
	assert.Equal(t, false, tx1.Verify(UtxoIndex, txPool, 0))
}
