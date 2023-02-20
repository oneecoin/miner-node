package blocks

import (
	"crypto/sha256"
	"fmt"

	"github.com/onee-only/miner-node/blockchain/transactions"
	"github.com/onee-only/miner-node/db"
	"github.com/onee-only/miner-node/lib"
)

var BlocksQueue []*Block = []*Block{}

func HashBlock(block *Block) string {
	hash := sha256.Sum256([]byte(fmt.Sprintf("%v", block)))
	return fmt.Sprintf("%x", hash)
}

func ValidateBlock(block *Block) bool {

	// validate transactions
	for _, tx := range block.Transactions {
		valid := transactions.ValidateTx(tx)
		if !valid {
			return false
		}
	}

	if block.PrevHash != getLastHash() {
		return false
	}
	if block.Height != getCurrentHeight()+1 {
		return false
	}

	// hash it to validate
	copyBlock := *block
	copyBlock.Hash = ""
	copyBlock.Timestamp = 0

	hash := HashBlock(&copyBlock)
	return hash == block.Hash
}

func SaveBroadcastedBlocks() {
	if len(BlocksQueue) != 0 {
		for _, block := range BlocksQueue {
			if valid := ValidateBlock(block); valid {
				AddBlock(block)
				updateCurrentHeight(block.Height)
				updateLastHash(block.Hash)
			}
		}
	}
}

func AddBlock(block *Block) {
	// block
	db.AddBlock(block.Hash, lib.ToBytes(block))

	// index
	txs := []db.IndexTx{}

	for _, tx := range block.Transactions {
		txs = append(txs, db.IndexTx{
			From: tx.TxIns.From,
			To:   tx.TxOuts[0].PublicKey,
		})
	}

	db.AddIndex(block.Height, block.Hash, txs)
}
