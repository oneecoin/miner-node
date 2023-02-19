package blocks

import (
	"strings"
	"time"

	"github.com/onee-only/miner-node/blockchain/transactions"
	"github.com/onee-only/miner-node/properties"
	"github.com/onee-only/miner-node/ws/peers"
)

const (
	difficulty = 4
)

func MineTxs(txs transactions.TxS) (*Block, transactions.TxS) {

	// validate the signatures
	invalidTxs := transactions.TxS{}
	for _, tx := range txs {
		if ok := validateTx(tx); !ok {
			invalidTxs = append(invalidTxs, tx)
		}
	}
	if len(invalidTxs) != 0 {
		// send some request
		return nil, invalidTxs
	}

	block := &Block{
		PrevHash:     lastHash,
		Height:       currentHeight,
		Hash:         "",
		Nonce:        0,
		Timestamp:    0,
		Transactions: txs,
	}

	target := strings.Repeat("0", difficulty)

	for {
		select {
		case m := <-peers.Peers.C:
			if m == properties.MessageBlockchainUploading {
				WaitForUpload()
			} else if m == properties.MessageNewBlock {
				// if new block is here
				// set the nonce to 0, and set the height and prevHash again
			}
		}
		hash := HashBlock(block)

		if strings.HasPrefix(hash, target) {
			block.Timestamp = int(time.Now().Local().Unix())
			// create and broadcast block
			printBlockStatus(block)
			break
		}
		printBlockStatus(block)
		block.Nonce++
	}

	return block, nil
}
