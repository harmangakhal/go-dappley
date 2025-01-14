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
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/jinzhu/copier"
	logger "github.com/sirupsen/logrus"

	"github.com/dappley/go-dappley/storage"
	"github.com/dappley/go-dappley/util"
)

var tipKey = []byte("tailBlockHash")

const LengthForBlockToBeConsideredHistory = 100

var (
	ErrBlockDoesNotExist   = errors.New("block does not exist")
	ErrTransactionNotFound = errors.New("transaction not found")
)

type Blockchain struct {
	tailBlockHash []byte
	db            storage.Storage
	consensus     Consensus
	txPool        *TransactionPool
}

// CreateBlockchain creates a new blockchain db
func CreateBlockchain(address Address, db storage.Storage, consensus Consensus, transactionPoolLimit uint32) *Blockchain {
	genesis := NewGenesisBlock(address)
	bc := &Blockchain{
		genesis.GetHash(),
		db,
		consensus,
		NewTransactionPool(transactionPoolLimit),
	}
	err := bc.AddBlockToTail(genesis)
	if err != nil {
		logger.Panic("CreateBlockchain: failed to add genesis block!")
	}
	return bc
}

func GetBlockchain(db storage.Storage, consensus Consensus, transactionPoolLimit uint32) (*Blockchain, error) {
	var tip []byte
	tip, err := db.Get(tipKey)
	if err != nil {
		return nil, err
	}

	bc := &Blockchain{
		tip,
		db,
		consensus,
		NewTransactionPool(transactionPoolLimit), //TODO: Need to retrieve transaction pool from db
	}
	if err != nil {
		return nil, err
	}
	return bc, nil
}

func (bc *Blockchain) GetDb() storage.Storage {
	return bc.db
}

func (bc *Blockchain) GetTailBlockHash() Hash {
	return bc.tailBlockHash
}

func (bc *Blockchain) GetConsensus() Consensus {
	return bc.consensus
}

func (bc *Blockchain) GetTxPool() *TransactionPool {
	return bc.txPool
}

func (bc *Blockchain) GetTailBlock() (*Block, error) {
	hash := bc.GetTailBlockHash()
	return bc.GetBlockByHash(hash)
}

func (bc *Blockchain) GetMaxHeight() uint64 {
	block, err := bc.GetTailBlock()
	if err != nil {
		return 0
	}
	return block.GetHeight()
}

func (bc *Blockchain) GetBlockByHash(hash Hash) (*Block, error) {
	rawBytes, err := bc.db.Get(hash)
	if err != nil {
		return nil, ErrBlockDoesNotExist
	}
	return Deserialize(rawBytes), nil
}

func (bc *Blockchain) GetBlockByHeight(height uint64) (*Block, error) {
	hash, err := bc.db.Get(util.UintToHex(height))
	if err != nil {
		return nil, ErrBlockDoesNotExist
	}

	return bc.GetBlockByHash(hash)
}

func (bc *Blockchain) SetTailBlockHash(tailBlockHash Hash) {
	bc.tailBlockHash = tailBlockHash
}

func (bc *Blockchain) SetConsensus(consensus Consensus) {
	bc.consensus = consensus
}

func (bc *Blockchain) AddBlockToTail(block *Block) error {
	blockLogger := logger.WithFields(logger.Fields{
		"height": block.GetHeight(),
		"hash":   hex.EncodeToString(block.GetHash()),
	})

	// Atomically set tail block hash and update UTXO index in db
	bcTemp := bc.deepCopy()

	bcTemp.db.EnableBatch()
	defer bcTemp.db.DisableBatch()

	err := bcTemp.setTailBlockHash(block.GetHash())
	if err != nil {
		blockLogger.Error("Blockchain: failed to set tail block hash!")
		return err
	}

	utxoIndex := LoadUTXOIndex(bcTemp.db)
	utxoIndex.UpdateUtxoState(block.GetTransactions())

	err = utxoIndex.Save(bcTemp.db)
	if err != nil {
		blockLogger.Error("Blockchain: failed to update UTXO index!")
		return err
	}

	err = bcTemp.AddBlockToDb(block)
	if err != nil {
		blockLogger.Warn("Blockchain: failed to add block to database")
		return err
	}

	// Flush batch changes to storage
	err = bcTemp.db.Flush()
	if err != nil {
		blockLogger.Error("Blockchain: failed to update tail block hash and UTXO index!")
		return err
	}

	// Assign changes to receiver
	*bc = *bcTemp

	blockLogger.Info("Blockchain: added a new block to tail.")

	return nil
}

//TODO: optimize performance
func (bc *Blockchain) FindTransaction(ID []byte) (Transaction, error) {
	bci := bc.Iterator()

	for {
		block, err := bci.Next()
		if err != nil {
			return Transaction{}, err
		}

		for _, tx := range block.GetTransactions() {
			if bytes.Compare(tx.ID, ID) == 0 {
				return *tx, nil
			}
		}

		if len(block.GetPrevHash()) == 0 {
			break
		}
	}

	return Transaction{}, ErrTransactionNotFound
}

func (bc *Blockchain) FindTransactionFromIndexBlock(txID []byte, blockId []byte) (Transaction, error) {
	bci := bc.Iterator()

	for {
		block, err := bci.NextFromIndex(blockId)
		if err != nil {
			return Transaction{}, err
		}

		for _, tx := range block.GetTransactions() {
			if bytes.Compare(tx.ID, txID) == 0 {
				return *tx, nil
			}
		}

		if len(block.GetPrevHash()) == 0 {
			break
		}
	}

	return Transaction{}, ErrTransactionNotFound
}

func (bc *Blockchain) Iterator() *Blockchain {
	return &Blockchain{bc.tailBlockHash, bc.db, bc.consensus, nil}
}

func (bc *Blockchain) Next() (*Block, error) {
	var block *Block

	encodedBlock, err := bc.db.Get(bc.tailBlockHash)
	if err != nil {
		return nil, err
	}

	block = Deserialize(encodedBlock)

	bc.tailBlockHash = block.GetPrevHash()

	return block, nil
}

func (bc *Blockchain) NextFromIndex(indexHash []byte) (*Block, error) {
	var block *Block

	encodedBlock, err := bc.db.Get(indexHash)
	if err != nil {
		return nil, err
	}

	block = Deserialize(encodedBlock)

	bc.tailBlockHash = block.GetPrevHash()
	println(bc.tailBlockHash)
	return block, nil
}

func (bc *Blockchain) String() string {
	var buffer bytes.Buffer

	bci := bc.Iterator()
	for {
		block, err := bci.Next()
		if err != nil {
			logger.Error(err)
		}

		buffer.WriteString(fmt.Sprintf("============ Block %x ============\n", block.GetHash()))
		buffer.WriteString(fmt.Sprintf("Height: %d\n", block.GetHeight()))
		buffer.WriteString(fmt.Sprintf("Prev. block: %x\n", block.GetPrevHash()))
		for _, tx := range block.GetTransactions() {
			buffer.WriteString(tx.String())
		}
		buffer.WriteString(fmt.Sprintf("\n\n"))

		if len(block.GetPrevHash()) == 0 {
			break
		}
	}
	return buffer.String()
}

//AddBlockToDb record the new block in the database
func (bc *Blockchain) AddBlockToDb(block *Block) error {

	err := bc.db.Put(block.GetHash(), block.Serialize())
	if err != nil {
		logger.WithError(err).Warn("Blockchain: failed to add block to database!")
		return err
	}

	err = bc.db.Put(util.UintToHex(block.GetHeight()), block.GetHash())
	if err != nil {
		logger.WithError(err).Warn("Blockchain: failed to index the block by block height in database!")
		return err
	}

	return nil
}

func (bc *Blockchain) IsHigherThanBlockchain(block *Block) bool {
	return block.GetHeight() > bc.GetMaxHeight()
}

func (bc *Blockchain) IsInBlockchain(hash Hash) bool {
	_, err := bc.GetBlockByHash(hash)
	return err == nil
}

//Verify all transactions in a fork
func VerifyTransactions(utxo UTXOIndex, forkBlks []*Block) bool {
	logger.Info("VerifyTransactions: is verifying transactions...")
	for i := len(forkBlks) - 1; i >= 0; i-- {
		logger.WithFields(logger.Fields{
			"height": forkBlks[i].GetHeight(),
			"hash":   hex.EncodeToString(forkBlks[i].GetHash()),
		}).Debug("VerifyTransactions: is verifying a block in the fork.")

		if !forkBlks[i].VerifyTransactions(utxo) {
			return false
		}

		utxo.UpdateUtxoState(forkBlks[i].GetTransactions())
	}
	return true
}

func (bc *Blockchain) addBlocksToTail(blocks []*Block) {
	if len(blocks) > 0 {
		for i := len(blocks) - 1; i >= 0; i-- {
			err := bc.AddBlockToTail(blocks[i])
			if err != nil {
				logger.WithError(err).Error("Blockchain: failed to add block to tail while concatenating fork!")
				return
			}
			//Remove transactions in current transaction pool
			bc.GetTxPool().RemoveMultipleTransactions(blocks[i].GetTransactions())
		}
	}
}

//rollback the blockchain to a block with the targetHash
func (bc *Blockchain) Rollback(targetHash Hash) bool {

	if !bc.IsInBlockchain(targetHash) {
		return false
	}
	parentblockHash := bc.GetTailBlockHash()
	//if is child of tail, skip rollback
	if IsHashEqual(parentblockHash, targetHash) {
		return true
	}

	//keep rolling back blocks until the block with the input hash
loop:
	for {
		if bytes.Compare(parentblockHash, targetHash) == 0 {
			break loop
		}
		block, err := bc.GetBlockByHash(parentblockHash)
		logger.WithFields(logger.Fields{
			"height": block.GetHeight(),
			"hash":   hex.EncodeToString(parentblockHash),
		}).Info("Blockchain: is about to rollback the block...")
		if err != nil {
			return false
		}
		parentblockHash = block.GetPrevHash()
		block.Rollback(bc.txPool)
	}

	err := bc.setTailBlockHash(parentblockHash)
	if err != nil {
		logger.Error("Blockchain: failed to set tail block hash during rollback!")
		return false
	}

	return true
}

func (bc *Blockchain) setTailBlockHash(hash Hash) error {
	err := bc.db.Put(tipKey, hash)
	if err != nil {
		return err
	}
	bc.tailBlockHash = hash
	return nil
}

func (bc *Blockchain) deepCopy() *Blockchain {
	newCopy := &Blockchain{}
	copier.Copy(&newCopy, &bc)
	return newCopy
}
