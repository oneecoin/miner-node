package mempool

import (
	"fmt"
	"time"

	"atomicgo.dev/cursor"
	"github.com/onee-only/miner-node/blockchain/blocks"
	"github.com/onee-only/miner-node/lib"
	"github.com/onee-only/miner-node/properties"
)

func ListenForMining() {
	interval := time.Minute * time.Duration(properties.CheckInterval)
	shouldWait := true
	for {
		var spent time.Duration = 0

		if shouldWait {
		intervalLoop:
			for {
				select {
				case m := <-properties.C:
					if m == properties.MessageBlockchainUploading {
						blocks.WaitForUpload()
					} else if m == properties.MessageNewBlock {
						blocks.HandleNewBlock()
					}
				default:
					fmt.Printf("Waiting to mine blocks... %s / %s",
						properties.WarningStr(fmtDuration(spent)),
						fmtDuration(interval))
					if interval == spent {
						break intervalLoop
					}
					time.Sleep(time.Second)
					spent += time.Second
				}
				cursor.ClearLine()
				cursor.StartOfLine()
			}
			cursor.ClearLine()
			cursor.StartOfLine()
			fmt.Println("Done waiting")
		}
		// time to mine some blocks

		s := lib.CreateSpinner(
			"Requesting transactions to mine",
			"Transactions received!",
		)

		// request transactions
		requestTxs()

		txs := <-mempool.transactionInbox
		if txs[0].ID == "" {
			s.FinalMSG = properties.ErrorStr("Rejected due to lack of transactions")
			s.Stop()
			continue
		}
		s.Stop()

		block, invalidTxs := blocks.MineTxs(txs)
		if block == nil {
			requestInvalidTxs(invalidTxs)
			time.Sleep(time.Second) // just in case
			shouldWait = false
			continue
		}

		blocks.AddBlock(block)
		properties.BlockBroadcastInbox <- lib.ToBytes(block)

		shouldWait = true
	}
}

func fmtDuration(d time.Duration) string {
	d = d.Round(time.Second)
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	return fmt.Sprintf("%02d:%02d", m, s)
}
