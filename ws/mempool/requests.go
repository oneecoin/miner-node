package mempool

import (
	"encoding/json"

	"github.com/onee-only/miner-node/blockchain/transactions"
	"github.com/onee-only/miner-node/lib"
	"github.com/onee-only/miner-node/properties"
	"github.com/onee-only/miner-node/ws/messages"
)

func requestTxs() {
	payload := messages.PayloadCount{
		Count: properties.MinTxs,
	}
	payloadBytes, err := json.Marshal(&payload)
	lib.HandleErr(err)

	m := messages.Message{
		Kind:    messages.MessageMempoolTxsRequest,
		Payload: payloadBytes,
	}

	mBytes, err := json.Marshal(&m)
	lib.HandleErr(err)

	mempool.inbox <- mBytes
}

func requestInvalidTxs(txs transactions.TxS) {

}
