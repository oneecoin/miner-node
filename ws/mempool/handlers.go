package mempool

import (
	"encoding/json"

	"github.com/onee-only/miner-node/lib"
	"github.com/onee-only/miner-node/ws/messages"
)

func handleMessage(m *messages.Message) {
	switch m.Kind {
	case messages.MessageBlocksRequest:
		payload := &messages.PayloadPage{}
		err := json.Unmarshal(m.Payload, payload)
		lib.HandleErr(err)
		sendBlocks(payload.Page)

	case messages.MessageBlockRequest:
		payload := &messages.PayloadHash{}
		err := json.Unmarshal(m.Payload, payload)
		lib.HandleErr(err)
		sendBlock(payload.Hash)

	case messages.MessageUTxOutsRequest:
		payload := &messages.PayloadUTxOutsFilter{}
		err := json.Unmarshal(m.Payload, payload)
		lib.HandleErr(err)
		sendUTxOuts(payload.PublicKey, payload.Amount)

	case messages.MessagePeerRejected:
		payload := &messages.PayloadPeer{}
		err := json.Unmarshal(m.Payload, payload)
		lib.HandleErr(err)
		rejectPeer(payload.PeerAddress)
	}
}

func sendBlocks(page int) {

}

func sendBlock(hash string) {

}

func sendUTxOuts(publicKey string, amount int) {

}

func rejectPeer(address string) {

}
