package blocks

import (
	"crypto/sha256"
	"fmt"
	"sort"

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

	var blocks []BlockSummary

	for _, blockBytes := range bytes {
		var block Block
		lib.FromBytes(block, blockBytes)
		blockSummary := BlockSummary{
			Hash:              block.Hash,
			Height:            block.Height,
			Timestamp:         block.Timestamp,
			TransactionsCount: len(block.Transactions),
		}
		blocks = append(blocks, blockSummary)
	}

	sort.Slice(blocks, func(i, j int) bool {
		return blocks[i].Timestamp > blocks[j].Timestamp
	})

	blocksJson := lib.ToJSON(blocks)
	return blocksJson
}

func FindBlock(hash string) []byte {
	var block Block
	bytes := db.FindBlockByHash(hash)
	lib.FromBytes(block, bytes)
	return lib.ToJSON(block)
}

func FindLatestBlock() []byte {
	var block Block
	bytes := db.FindBlockByHash(lastHash)
	lib.FromBytes(block, bytes)
	return lib.ToJSON(block)
}

func FindTxsByPublicKey(publicKey string) transactions.TxS {
	hashes := db.FindHashesAll(publicKey)
	bytes := db.FindBlocksByHashes(hashes)

	txs := transactions.TxS{}

	for _, blockBytes := range bytes {
		var block Block
		lib.FromBytes(block, blockBytes)
		for _, tx := range block.Transactions {
			if tx.TxIns.From == publicKey || tx.TxOuts[0].PublicKey == publicKey {
				txs = append(txs, tx)
			}
		}
	}

	sort.Slice(txs, func(i, j int) bool {
		return txs[i].Timestamp > txs[j].Timestamp
	})

	return txs
}

func FindUTxOutsByPublicKey(publicKey string, amount int) (transactions.UTxOutS, bool, int) {
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
		return nil, false, 0
	}
	return uTxOuts, true, got
}

func FindBalanceByPublicKey(publicKey string) int {
	_, _, balance := FindUTxOutsByPublicKey(publicKey, 0)
	return balance
}
