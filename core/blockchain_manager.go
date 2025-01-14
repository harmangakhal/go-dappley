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
	"encoding/hex"

	"github.com/libp2p/go-libp2p-peer"
	logger "github.com/sirupsen/logrus"

	"github.com/dappley/go-dappley/common"
)

type BlockChainManager struct {
	blockchain *Blockchain
	blockPool  *BlockPool
}

func NewBlockChainManager() *BlockChainManager {
	return &BlockChainManager{}
}

func (bm *BlockChainManager) SetblockPool(blockPool *BlockPool) {
	bm.blockPool = blockPool
}

func (bm *BlockChainManager) Setblockchain(blockchain *Blockchain) {
	bm.blockchain = blockchain
}

func (bm *BlockChainManager) Getblockchain() *Blockchain {
	return bm.blockchain
}

func (bm *BlockChainManager) GetblockPool() *BlockPool {
	return bm.blockPool
}

func (bm *BlockChainManager) Push(block *Block, pid peer.ID) {
	logger.WithFields(logger.Fields{
		"from": pid.String(),
		"hash": hex.EncodeToString(block.GetHash()),
	}).Info("BlockChainManager: received a new block.")
	if !bm.blockPool.Verify(block) {
		return
	}
	logger.Debug("BlockChainManager: block is verified.")
	if !(bm.blockchain.GetConsensus().Validate(block)) {
		logger.Warn("BlockChainManager: block is invalid according to consensus!")
		return
	}
	logger.Debug("BlockChainManager: block is valid according to consensus.")
	tree, _ := common.NewTree(block.GetHash().String(), block)
	forkheadParentHash := bm.blockPool.HandleRecvdBlock(tree, bm.blockchain.GetMaxHeight())
	if forkheadParentHash == nil {
		return
	}
	if parent, _ := bm.blockchain.GetBlockByHash(forkheadParentHash); parent == nil {
		logger.WithFields(logger.Fields{
			"parent_hash":   hex.EncodeToString(tree.GetValue().(*Block).GetPrevHash()),
			"parent_height": tree.GetValue().(*Block).GetHeight() - 1,
		}).Info("BlockChainManager: cannot find the parent of the received block from blockchain.")
		bm.blockPool.requestPrevBlock(tree, pid)
		return
	}
	forkBlks := bm.blockPool.GenerateForkBlocks(tree, bm.blockchain.GetMaxHeight())
	bm.MergeFork(forkBlks, forkheadParentHash)
	bm.blockPool.CleanCache(tree)
}

func (bm *BlockChainManager) MergeFork(forkBlks []*Block, forkParentHash Hash) {

	//find parent block
	if len(forkBlks) == 0 {
		return
	}
	forkHeadBlock := forkBlks[len(forkBlks)-1]
	if forkHeadBlock == nil {
		return
	}

	//verify transactions in the fork
	utxo, err := GetUTXOIndexAtBlockHash(bm.blockchain.db, bm.blockchain, forkParentHash)
	if err != nil {
		logger.Error("BlockChainManager: blockchain is corrupted! Delete the database file and resynchronize to the network.")
		return
	}

	if !VerifyTransactions(*utxo, forkBlks) {
		return
	}

	bm.blockchain.Rollback(forkParentHash)

	//add all blocks in fork from head to tail
	bm.blockchain.addBlocksToTail(forkBlks)

}
