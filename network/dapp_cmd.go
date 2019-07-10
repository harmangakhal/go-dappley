// Copyright (C) 2018 go-dappley authors
//
// This file is part of the go-dappley library.
//
// the go-dappley library is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// the go-dappley library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with the go-dappley library.  If not, see <http://www.gnu.org/licenses/>.
//

package network

import (
	"github.com/golang/protobuf/proto"
	logger "github.com/sirupsen/logrus"

	"github.com/dappley/go-dappley/network/pb"
)

type DappCmd struct {
	cmd         string
	data        []byte
	isBroadcast bool
}

func NewDapCmd(cmd string, data []byte, isBroadcast bool) *DappCmd {
	return &DappCmd{cmd, data, isBroadcast}
}

func (dc *DappCmd) GetCmd() string {
	return dc.cmd
}

func (dc *DappCmd) GetData() []byte {
	return dc.data
}

func ParseDappMsgFromDappPacket(packet *DappPacket) *DappCmd {
	return ParseDappMsgFromRawBytes(packet.GetData())
}

func ParseDappMsgFromRawBytes(bytes []byte) *DappCmd {
	dmpb := &networkpb.DappCmd{}

	//unmarshal byte to proto
	if err := proto.Unmarshal(bytes, dmpb); err != nil {
		logger.WithError(err).Warn("Stream: Unable to")
	}

	dm := &DappCmd{}
	dm.FromProto(dmpb)
	return dm
}

func (dc *DappCmd) GetRawBytes() []byte {
	data, err := proto.Marshal(dc.ToProto())
	if err != nil {
		logger.WithError(err).Error("DappCmd: Dapp Command can not be converted into raw bytes")
	}
	return data
}

func (dc *DappCmd) ToProto() proto.Message {
	return &networkpb.DappCmd{
		Cmd:         dc.cmd,
		Data:        dc.data,
		IsBroadcast: dc.isBroadcast,
	}
}

func (dc *DappCmd) FromProto(pb proto.Message) {
	dc.cmd = pb.(*networkpb.DappCmd).GetCmd()
	dc.data = pb.(*networkpb.DappCmd).GetData()
	dc.isBroadcast = pb.(*networkpb.DappCmd).GetIsBroadcast()
}