package blocks

import (
	"crypto/sha256"
	"fmt"

	"github.com/onee-only/miner-node/blockchain/transactions"
)

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

}

func AddBlock(block *Block) {

}
