package blocks

import (
	"crypto/sha256"
	"fmt"

	"github.com/onee-only/miner-node/db"
)

var (
	lastHash      string
	currentHeight int
)

func HashBlock(block *Block) string {
	hash := sha256.Sum256([]byte(fmt.Sprintf("%v", block)))
	return fmt.Sprintf("%x", hash)
}

func getLastHash() string {
	if lastHash == "" {
		lastHash = db.FindLastHash()
	}
	return lastHash
}

func getCurrentHeight() int {
	if currentHeight == 0 {
		currentHeight = db.FindCurrentHeight()
	}
	return currentHeight
}

func SaveBroadcastedBlocks() {

}

func AddBlock(block *Block) {

}
