package blocks

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"

	"github.com/onee-only/miner-node/blockchain/transactions"
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

func validateTx(tx *transactions.Tx) bool {
	
	x, y, err := lib.RestoreBigInts(tx.TxIns.From)
	if err != nil {
		return false
	}

	txId

	for _, txIn := range tx.TxIns.V {
		ecdsa.Verify(&ecdsa.PublicKey{
			Curve: elliptic.P256(),
			X: x,
			Y: y,
		}, ,)
	}
}

func printBlockStatus(block *Block) {

}
