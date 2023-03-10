package mempool

import (
	"github.com/onee-only/miner-node/blockchain/transactions"
	"github.com/onee-only/miner-node/lib"
	"github.com/onee-only/miner-node/properties"
	"github.com/onee-only/miner-node/ws/messages"
)

func requestTxs() {
	payload := messages.PayloadCount{
		Count: properties.MinTxs,
	}
	payloadBytes := lib.ToJSON(&payload)

	m := messages.Message{
		Kind:    messages.MessageMempoolTxsRequest,
		Payload: payloadBytes,
	}

	mBytes := lib.ToJSON(&m)

	mempool.inbox <- mBytes
}

func requestInvalidTxs(txs transactions.TxS) {
	payload := messages.PayloadTxs{
		Txs: txs,
	}

	m := messages.Message{
		Kind:    messages.MessageInvalidTxsRequest,
		Payload: lib.ToJSON(payload),
	}

	mempool.inbox <- lib.ToJSON(m)
}

func listenRequestRejectPeer() {
	for {
		address := <-properties.RejectPeerInbox

		m := messages.Message{
			Kind:    messages.MessageRejectPeer,
			Payload: lib.ToJSON(messages.PayloadPeer{PeerAddress: address}),
		}

		mempool.inbox <- lib.ToJSON(m)
	}
}

func listenRequestNewBlock() {
	for {
		address := <-properties.NewBlockInbox

		m := messages.Message{
			Kind: messages.MessageNewBlock,
			Payload: lib.ToJSON(messages.PayloadPeer{
				PeerAddress: address,
			}),
		}

		mempool.inbox <- lib.ToJSON(m)
	}
}
