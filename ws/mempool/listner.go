package mempool

import (
	"fmt"
	"time"

	"atomicgo.dev/cursor"
	"github.com/onee-only/miner-node/properties"
	"github.com/onee-only/miner-node/ws/peers"
)

func ListenForMining() {
	for {
		interval := time.Minute * time.Duration(properties.CheckInterval)
		var spent time.Duration = 0
		for {
			select {
			case <-peers.Peers.C:

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
	}
}

func fmtDuration(d time.Duration) string {
	d = d.Round(time.Second)
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	return fmt.Sprintf("%02d:%02d", m, s)
}
