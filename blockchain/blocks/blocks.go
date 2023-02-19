package blocks

import (
	"crypto/sha256"
	"fmt"
)

var (
	lastHash      string
	currentHeight int
)

func HashBlock(block *Block) string {
	hash := sha256.Sum256([]byte(fmt.Sprintf("%v", block)))
	return fmt.Sprintf("%x", hash)
}
