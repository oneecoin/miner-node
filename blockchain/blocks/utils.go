package blocks

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"fmt"
	"os"

	"atomicgo.dev/cursor"
	"github.com/olekukonko/tablewriter"
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

		// should see if there is actual transaction

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

var tbl = tablewriter.NewWriter(os.Stdout)

func printTable(txsCount int, prevHash string) {

	data := []string{
		fmt.Sprintf("%d", difficulty),
		fmt.Sprintf("%d", txsCount),
		fmt.Sprintf("%d", currentHeight),
		prevHash,
	}

	fmt.Println("Start hashing with:\n")
	tbl.SetHeader([]string{"Difficulty", "Transactions", "Height", "PrevHash"})
	tbl.ClearRows()
	tbl.Append(data)
	tbl.Render()
	fmt.Printf("\n\n\n\n")
}

func printBlockStatus(nonce int, hash string) {
	cursor.ClearLinesUp(2)
	fmt.Println("NONCE\t\tHASH")
	fmt.Printf("%d\t\t%s\n", nonce, hash)
}
