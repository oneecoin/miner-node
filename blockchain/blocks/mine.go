package blocks

import (
	"fmt"
	"strings"
	"time"

	"github.com/onee-only/miner-node/blockchain/transactions"
	"github.com/onee-only/miner-node/lib"
	"github.com/onee-only/miner-node/properties"
	"github.com/onee-only/miner-node/ws/peers"
)

const (
	difficulty = 4
)

func MineTxs(txs transactions.TxS) (*Block, transactions.TxS) {

	s := lib.CreateSpinner(
		"Validating transactions",
		"Transactions validated!",
	)

	// validate the signatures
	invalidTxs := transactions.TxS{}
	for _, tx := range txs {
		if ok := validateTx(tx); !ok {
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
		PrevHash:     lastHash,
		Height:       currentHeight,
		Hash:         "",
		Nonce:        0,
		Timestamp:    0,
		Transactions: txs,
	}

	target := strings.Repeat("0", difficulty)

	printTable(len(txs), lastHash)

	for {
		select {
		case m := <-peers.Peers.C:
			if m == properties.MessageBlockchainUploading {
				fmt.Println()
				WaitForUpload()
			} else if m == properties.MessageNewBlock {
				// if new block is here
				// set the nonce to 0, and set the height and prevHash again

				fmt.Printf("\n%s\n", properties.WarningStr("New block broadcasted. Should reset config"))
				printTable(len(txs), lastHash)
			}
		default:

		}
		hash := HashBlock(block)

		if strings.HasPrefix(hash, target) {
			block.Timestamp = int(time.Now().Local().Unix())
			// create and broadcast block
			printBlockStatus(block.Nonce, block.Hash)
			break
		}
		printBlockStatus(block.Nonce, block.Hash)
		block.Nonce++
	}

	return block, nil
}
