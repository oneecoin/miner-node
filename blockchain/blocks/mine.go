package blocks

import (
	"fmt"
	"strings"
	"time"

	"github.com/onee-only/miner-node/blockchain/transactions"
	"github.com/onee-only/miner-node/lib"
	"github.com/onee-only/miner-node/properties"
)

const (
	difficulty = 5
)

func MineTxs(txs transactions.TxS) (*Block, transactions.TxS) {

	s := lib.CreateSpinner(
		"Validating transactions",
		"Transactions validated!",
	)

	// validate the signatures
	invalidTxs := transactions.TxS{}
	for _, tx := range txs {
		if ok := transactions.ValidateTx(tx); !ok {
			invalidTxs = append(invalidTxs, tx)
		}
	}
	if len(invalidTxs) != 0 {
		s.FinalMSG = properties.ErrorStr("Invalid transaction(s) found and rejected")
		s.Stop()
		return nil, invalidTxs
	}
	s.Stop()

	block := &Block{
		PrevHash:     getLastHash(),
		Height:       getCurrentHeight() + 1,
		Hash:         "",
		Nonce:        0,
		Timestamp:    0,
		Transactions: txs,
	}

	target := strings.Repeat("0", difficulty)

	printTable(len(txs), getLastHash())

	for {
		select {
		case m := <-properties.C:
			if m == properties.MessageBlockchainUploading {
				fmt.Println()
				WaitForUpload()
			} else if m == properties.MessageNewBlock {

				HandleNewBlock()

				// set things of my block
				block.PrevHash = getLastHash()
				block.Height = getCurrentHeight() + 1
				block.Nonce = 0

				fmt.Printf("\n%s\n", properties.WarningStr("New block broadcasted. Should reset config"))
				printTable(len(txs), getLastHash())
			}
		default:

		}
		block.Hash = HashBlock(block)

		if strings.HasPrefix(block.Hash, target) {
			block.Timestamp = int(time.Now().Local().Unix())
			printBlockStatus(block.Nonce, block.Hash)
			break
		}
		printBlockStatus(block.Nonce, block.Hash)
		block.Nonce++
	}

	return block, nil
}
