package peers

import (
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/onee-only/miner-node/lib"
	"github.com/onee-only/miner-node/properties"
	"github.com/onee-only/miner-node/ws/messages"
)

type TAddress struct {
	Host string
	Port string
}

type Peer struct {
	Conn    *websocket.Conn
	Inbox   chan []byte
	Address TAddress
}

func (p Peer) GetAddress() string {
	return fmt.Sprintf("%s:%s", p.Address.Host, p.Address.Port)
}

func (p *Peer) closeConn() {
	Peers.M.Lock()
	defer Peers.M.Unlock()
	p.Conn.Close()
	delete(Peers.V, p.GetAddress())
}

func (p *Peer) read() {
	defer p.closeConn()
	for {
		m := &messages.Message{}
		messageType, payload, err := p.Conn.ReadMessage()
		if err != nil {
			break
		}

		if messageType == websocket.BinaryMessage && len(payload) == 8 {
			downloadBlockchain(payload, p.Conn)
		} else {
			if properties.IsDownloading && m.Kind == messages.MessageNewBlock {
				// take this in the queue or something.
			} else {
				lib.FromJSON(payload, m)
				Peers.handleMessage(m, p)
			}
		}
	}
}

func (p *Peer) write() {
	defer p.closeConn()
	for {
		m, ok := <-p.Inbox
		if !ok {
			break
		}
		p.Conn.WriteMessage(websocket.TextMessage, m)
	}
}
