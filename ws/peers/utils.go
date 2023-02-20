package peers

import (
	"errors"
	"io"
	"math/rand"

	"github.com/gorilla/websocket"
	"github.com/onee-only/miner-node/lib"
	"github.com/onee-only/miner-node/properties"
	"github.com/onee-only/miner-node/ws/messages"
)

func getRandomPeer() *Peer {
	countLeft := 0
	if len(Peers.V) != 0 {
		countLeft = rand.Intn(len(Peers.V))
	}
	var v *Peer
	for _, peer := range Peers.V {
		if countLeft == 0 {
			v = peer
			break
		}
		countLeft--
	}
	return v
}

type WebSocketReader struct {
	conn *websocket.Conn
}

func (r *WebSocketReader) Read(p []byte) (int, error) {
	messageType, payload, err := r.conn.ReadMessage()
	if messageType != websocket.BinaryMessage || err != nil {
		lib.HandleErr(errors.New("wtf this should not happen"))
	}
	if len(payload) == 0 {
		return 0, io.EOF
	}
	return copy(p, payload), nil
}

func listenBlockBroadcast() {
	for {
		block := <-properties.BlockBroadcastInbox

		m := messages.Message{
			Kind:    messages.MessageNewBlock,
			Payload: block,
		}

		for _, p := range Peers.V {
			p.Inbox <- lib.ToJSON(m)
		}
	}
}

func listenPeerRejected() {
	for {
		address := <-properties.PeerRejectedInbox
		peer := Peers.V[address]
		close(peer.Inbox)
	}
}
