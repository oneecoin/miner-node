package mempool

import (
	"github.com/onee-only/miner-node/blockchain/blocks"
	"github.com/onee-only/miner-node/blockchain/transactions"
	"github.com/onee-only/miner-node/lib"
	"github.com/onee-only/miner-node/properties"
	"github.com/onee-only/miner-node/ws/messages"
)

func handleMessage(m *messages.Message) {
	switch m.Kind {
	case messages.MessageBlocksRequest:
		payload := &messages.PayloadPage{}
		lib.FromJSON(m.Payload, payload)
		sendBlocks(payload.Page)

	case messages.MessageBlockRequest:
		if m.Payload != nil {
			payload := &messages.PayloadHash{}
			lib.FromJSON(m.Payload, payload)
			sendBlock(payload.Hash)
		} else {
			sendLatest()
		}

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
	case messages.MessageNodeTxsRequest:
		payload := &messages.PayloadHash{}
		lib.FromJSON(m.Payload, payload)

		txs := blocks.FindTxsByPublicKey(payload.Hash)
		m := messages.Message{
			Kind: messages.MessageNodeTxsResponse,
			Payload: lib.ToJSON(messages.PayloadTxs{
				Txs: txs,
			}),
		}

		mempool.inbox <- lib.ToJSON(m)
	case messages.MessageBalanceRequest:
		payload := &messages.PayloadHash{}
		lib.FromJSON(m.Payload, payload)

		balance := blocks.FindBalanceByPublicKey(payload.Hash)

		m := messages.Message{
			Kind:    messages.MessageBalanceResponse,
			Payload: lib.ToJSON(messages.PayloadCount{Count: balance}),
		}

		mempool.inbox <- lib.ToJSON(m)
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

func sendLatest() {
	m := messages.Message{
		Kind:    messages.MessageBlockResponse,
		Payload: blocks.FindLatestBlock(),
	}

	mempool.inbox <- lib.ToJSON(m)
}

func sendUTxOuts(publicKey string, amount int) {

	uTxOuts, available, _ := blocks.FindUTxOutsByPublicKey(publicKey, amount)
	need := transactions.UTxOutS{}

	if available {
		got := 0
		for _, uTxOut := range uTxOuts {
			if got >= amount {
				break
			}
			got += uTxOut.Amount
			need = append(need, uTxOut)
		}
	}

	payload := messages.PayloadUTxOuts{
		Available: available,
		UTxOuts:   need,
	}

	m := messages.Message{
		Kind:    messages.MessageUTxOutsResponse,
		Payload: lib.ToJSON(payload),
	}

	mempool.inbox <- lib.ToJSON(m)
}

func rejectPeer(address string) {
	properties.PeerRejectedInbox <- address
}
