package blocks

import (
	"fmt"

	"github.com/onee-only/miner-node/lib"
	"github.com/onee-only/miner-node/properties"
	"github.com/onee-only/miner-node/ws/peers"
)

func WaitForUpload() {
	fmt.Println("Blockchain upload requested")
	s := lib.CreateSpinner(
		"Uploading blockchain",
		"Blockchain successfully uploaded!",
	)
	m := <-peers.Peers.C
	if m == properties.MessageBlockchainUploaded {
		s.Stop()
	}
}

func printBlockStatus(block *Block) {

}
