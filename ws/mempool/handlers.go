package mempool

import (
	"github.com/onee-only/miner-node/lib"
	"github.com/onee-only/miner-node/ws/messages"
)

func handleMessage(m *messages.Message) {
	switch m.Kind {
	case messages.MessageBlocksRequest:
		payload := &messages.PayloadPage{}
		lib.FromJSON(m.Payload, payload)
		sendBlocks(payload.Page)

	case messages.MessageBlockRequest:
		payload := &messages.PayloadHash{}
		lib.FromJSON(m.Payload, payload)
		sendBlock(payload.Hash)

	case messages.MessageUTxOutsRequest:
		payload := &messages.PayloadUTxOutsFilter{}
		lib.FromJSON(m.Payload, payload)
		sendUTxOuts(payload.PublicKey, payload.Amount)

	case messages.MessagePeerRejected:
		payload := &messages.PayloadPeer{}
		lib.FromJSON(m.Payload, payload)
		rejectPeer(payload.PeerAddress)

	case messages.MessageTxsMempoolResponse:
		payload := &messages.PayloadTxs{}
		lib.FromJSON(m.Payload, payload)
		mempool.transactionInbox <- payload.Txs

	case messages.MessageTxsDeclined:
		mempool.transactionInbox <- nil
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
