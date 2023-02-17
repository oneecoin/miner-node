package mempool

import (
	"fmt"
	"time"

	"atomicgo.dev/cursor"
	"github.com/onee-only/miner-node/lib"
	"github.com/onee-only/miner-node/properties"
	"github.com/onee-only/miner-node/ws/peers"
)

func ListenForMining() {
	interval := time.Minute * time.Duration(properties.CheckInterval)
	for {
		var spent time.Duration = 0
		for {
			select {
			case m := <-peers.Peers.C:
				if m == properties.MessageBlockchainUploading {
					fmt.Println("Blockchain upload requested")
					s := lib.CreateSpinner(
						"Uploading blockchain",
						"Blockchain successfully uploaded!",
					)
					m = <-peers.Peers.C
					if m == properties.MessageBlockchainUploaded {
						s.Stop()
					}
				}
			default:
				fmt.Printf("Waiting to mine blocks... %s / %s",
					properties.WarningStr(fmtDuration(spent)),
					fmtDuration(interval))
				if interval == spent {
					break
				}
				time.Sleep(time.Second)
				spent += time.Second
			}
			cursor.ClearLine()
		}
		// time to mine some blocks
	}
}

func fmtDuration(d time.Duration) string {
	d = d.Round(time.Second)
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	return fmt.Sprintf("%02d:%02d", m, s)
}
