package blocks

import "github.com/onee-only/miner-node/db"

var (
	lastHash      string
	currentHeight int
)

func getLastHash() string {
	if lastHash == "" {
		lastHash = db.FindLastHash()
	}
	return lastHash
}

func updateLastHash(hash string) {
	lastHash = hash
}

func getCurrentHeight() int {
	if currentHeight == 0 {
		currentHeight = db.FindCurrentHeight()
	}
	return currentHeight
}

func updateCurrentHeight(height int) {
	currentHeight = height
}
