package blocks

import (
	"fmt"
	"os"

	"atomicgo.dev/cursor"
	"github.com/olekukonko/tablewriter"
	"github.com/onee-only/miner-node/lib"
	"github.com/onee-only/miner-node/properties"
)

func WaitForUpload() {
	fmt.Println("Blockchain upload requested")
	s := lib.CreateSpinner(
		"Uploading blockchain",
		"Blockchain successfully uploaded!",
	)
	m := <-properties.C
	if m == properties.MessageBlockchainUploaded {
		s.Stop()
	}
}

var tbl = tablewriter.NewWriter(os.Stdout)

func printTable(txsCount int, prevHash string) {

	data := []string{
		fmt.Sprintf("%d", difficulty),
		fmt.Sprintf("%d", txsCount),
		fmt.Sprintf("%d", getCurrentHeight()),
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

func HandleNewBlock() {
	blockBytes := <-properties.BlockReceiveInbox
	newBlock := &Block{}
	lib.FromBytes(newBlock, blockBytes)

	AddBlock(newBlock)
	updateCurrentHeight(newBlock.Height)
	updateLastHash(newBlock.Hash)
}
