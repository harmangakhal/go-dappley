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
	"github.com/dappley/go-dappley/core"
	"github.com/libp2p/go-libp2p-peer"
	"github.com/multiformats/go-multiaddr"
)

func FakeNodeWithPeer(pid, addr string) *Node {

	node := NewNode(nil, nil)
	peerid, _ := peer.IDB58Decode(pid)
	maddr, _ := multiaddr.NewMultiaddr(addr)
	p := &Peer{peerid, maddr}
	node.GetPeerList().Add(p)

	return node
}

func FakeNodeWithPidAndAddr(pool *core.BlockPool, bc *core.Blockchain, pid, addr string) *Node {

	node := NewNode(bc, pool)
	peerid, _ := peer.IDB58Decode(pid)
	maddr, _ := multiaddr.NewMultiaddr(addr)
	p := &Peer{peerid, maddr}
	node.info = p

	return node
}
