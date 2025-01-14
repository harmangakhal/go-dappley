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

	"github.com/dappley/go-dappley/common"
	"github.com/dappley/go-dappley/core/pb"
	"github.com/gogo/protobuf/proto"
)

type TXOutput struct {
	Value      *common.Amount
	PubKeyHash PubKeyHash
	Contract   string
}

func (out *TXOutput) Lock(address Address) {
	hash, _ := address.GetPubKeyHash()
	out.PubKeyHash = PubKeyHash{hash}
}

func (out *TXOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash.GetPubKeyHash(), pubKeyHash) == 0
}

func NewTXOutput(value *common.Amount, address Address) *TXOutput {
	txo := &TXOutput{value, PubKeyHash{}, ""}
	txo.Lock(address)
	return txo
}

func NewContractTXOutput(contract string) *TXOutput {
	txo := &TXOutput{common.NewAmount(0), NewContractPubKeyHash(), contract}
	return txo
}

func (out *TXOutput) ToProto() proto.Message {
	return &corepb.TXOutput{
		Value:      out.Value.Bytes(),
		PubKeyHash: out.PubKeyHash.GetPubKeyHash(),
		Contract:   out.Contract,
	}
}

func (out *TXOutput) FromProto(pb proto.Message) {
	out.Value = common.NewAmountFromBytes(pb.(*corepb.TXOutput).Value)
	out.PubKeyHash = PubKeyHash{pb.(*corepb.TXOutput).PubKeyHash}
	out.Contract = pb.(*corepb.TXOutput).Contract
}
