package mempool

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/gorilla/websocket"
	"github.com/onee-only/miner-node/config"
)

var mempoolConn *websocket.Conn

func Connect() {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Prefix = "Connecting mempool server "
	s.FinalMSG = "Mempool server connected!\n"
	s.Start()
	time.Sleep(5 * time.Second)
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s/ws?port=%d&publicKey=%s", config.MempoolAddress, config.Port, config.PublicKey), nil)
	if err != nil {
		panic(err)
	}
	mempoolConn = conn
	s.Stop()
}
