package mempool

import (
	"github.com/onee-only/miner-node/blockchain/blocks"
	"github.com/onee-only/miner-node/lib"
	"github.com/onee-only/miner-node/ws/messages"
	"github.com/onee-only/miner-node/ws/peers"
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
	m := messages.Message{
		Kind:    messages.MessageBlocksResponse,
		Payload: blocks.FindBlocksWithPage(page),
	}

	mempool.inbox <- lib.ToJSON(m)
}

func sendBlock(hash string) {
	m := messages.Message{
		Kind:    messages.MessageBlockResponse,
		Payload: blocks.FindBlock(hash),
	}

	mempool.inbox <- lib.ToJSON(m)
}

func sendUTxOuts(publicKey string, amount int) {

	uTxOuts, available := blocks.FindUTxOutsByPublicKey(publicKey, amount)

	payload := messages.PayloadUTxOuts{
		Available: available,
		UTxOuts:   uTxOuts,
	}

	m := messages.Message{
		Kind:    messages.MessageUTxOutsResponse,
		Payload: lib.ToJSON(payload),
	}

	mempool.inbox <- lib.ToJSON(m)
}

func rejectPeer(address string) {
	peer := peers.Peers.V[address]
	close(peer.Inbox)
}
