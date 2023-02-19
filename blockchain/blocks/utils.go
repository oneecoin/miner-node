package blocks

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
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

	hash, err := hex.DecodeString(tx.ID)
	if err != nil {
		return false
	}

	for _, txIn := range tx.TxIns.V {

		// see if there is actual transaction

		r, s, err := lib.RestoreBigInts(txIn.Signature)
		if err != nil {
			return false
		}

		valid := ecdsa.Verify(&ecdsa.PublicKey{
			Curve: elliptic.P256(),
			X:     x,
			Y:     y,
		}, hash, r, s)
		if !valid {
			return false
		}
	}
	return true
}

func printBlockStatus(block *Block) {

}
