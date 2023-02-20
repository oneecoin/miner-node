package blocks

import (
	"crypto/sha256"
	"fmt"

	"github.com/onee-only/miner-node/blockchain/transactions"
	"github.com/onee-only/miner-node/db"
	"github.com/onee-only/miner-node/lib"
	"github.com/onee-only/miner-node/properties"
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

func FindBlocksWithPage(page int) []byte {
	start := getCurrentHeight() - (properties.DefaultPageSize * (page - 1))
	hash := db.FindHashByHeight(start)
	bytes := db.FindBlocksPageByHash(hash)

	var blocks []Block

	for _, blockBytes := range bytes {
		var block Block
		lib.FromBytes(block, blockBytes)
		blocks = append(blocks, block)
	}
	// should sort it
	blocksJson := lib.ToJSON(blocks)
	return blocksJson
}

func FindBlock(hash string) []byte {
	var block Block
	bytes := db.FindBlockByHash(hash)
	lib.FromBytes(block, bytes)
	return lib.ToJSON(block)
}

func FindUTxOutsByPublicKey(publicKey string, amount int) (transactions.UTxOutS, bool) {
	spentAt := db.FindHashesFrom(publicKey)
	earnedAt := db.FindHashesTo(publicKey)

	spentMap := make(map[string]string)
	uTxOuts := transactions.UTxOutS{}
	got := 0

	bytes := db.FindBlocksByHashes(spentAt)
	for _, blockBytes := range bytes {
		var block Block
		lib.FromBytes(block, blockBytes)
		for _, tx := range block.Transactions {
			if tx.TxIns.From == publicKey {
				for _, txIn := range tx.TxIns.V {
					spentMap[txIn.BlockHash] = txIn.TxID
				}
			}
		}
	}

	bytes = db.FindBlocksByHashes(earnedAt)
	for _, blockBytes := range bytes {
		var block Block
		lib.FromBytes(block, blockBytes)
		for _, tx := range block.Transactions {
			for index, txOut := range tx.TxOuts {
				if txOut.PublicKey == publicKey {
					txId, exists := spentMap[block.Hash]
					if exists && tx.ID == txId {
						continue
					}
					got += txOut.Amount
					uTxOuts = append(uTxOuts, &transactions.UTxOut{
						BlockHash: block.Hash,
						TxID:      tx.ID,
						Index:     index,
						Amount:    txOut.Amount,
					})
				}
			}
		}
	}

	if amount > got {
		return nil, false
	}
	return uTxOuts, true
}
