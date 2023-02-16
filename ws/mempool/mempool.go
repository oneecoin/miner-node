package mempool

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/gorilla/websocket"
	"github.com/onee-only/miner-node/config"
	"github.com/onee-only/miner-node/lib"
	"github.com/onee-only/miner-node/ws/messages"
)

type tMempool struct {
	conn  *websocket.Conn
	inbox chan []byte
}

var mempool tMempool = tMempool{
	inbox: make(chan []byte),
}

func (tMempool) read() {
	defer mempool.conn.Close()
	for {
		m := &messages.Message{}
		mempool.conn.ReadJSON(m)
		handleMessage(m)
	}
}

func (tMempool) write() {
	defer mempool.conn.Close()
	for {
		m := <-mempool.inbox
		mempool.conn.WriteMessage(websocket.TextMessage, m)
	}
}

func Connect() *[]string {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Prefix = "Connecting to mempool server "
	s.FinalMSG = "Mempool server connected!\n"
	s.Start()
	conn, res, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s/ws?port=%d&publicKey=%s", config.MempoolAddress, config.Port, config.PublicKey), nil)
	lib.HandleErr(err)
	mempool.conn = conn
	go mempool.read()
	go mempool.write()

	var peerList []string
	err = json.NewDecoder(res.Body).Decode(&peerList)
	lib.HandleErr(err)
	s.Stop()
	return &peerList
}
