package mempool

import (
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/onee-only/miner-node/blockchain/transactions"
	"github.com/onee-only/miner-node/lib"
	"github.com/onee-only/miner-node/properties"
	"github.com/onee-only/miner-node/ws/messages"
)

type tMempool struct {
	conn             *websocket.Conn
	inbox            chan []byte
	transactionInbox chan transactions.TxS
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

func Connect() {
	s := lib.CreateSpinner(
		"Connecting to mempool server",
		"Mempool server connected!",
	)

	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s/ws?port=%d&publicKey=%s", properties.MempoolAddress, properties.Port, properties.PublicKey), nil)
	lib.HandleErr(err)
	mempool.conn = conn
	go mempool.read()
	go mempool.write()
	s.Stop()
}
