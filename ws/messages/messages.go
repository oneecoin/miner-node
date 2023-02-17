package messages

import "github.com/onee-only/miner-node/blockchain/transactions"

type Message struct {
	Kind    MessageKind
	Payload []byte
}

type PayloadPeer struct {
	PeerAddress string
}

type PayloadPage struct {
	Page int
}

type PayloadHash struct {
	Hash string
}

type PayloadUTxOutsFilter struct {
	PublicKey string
	Amount    int
}

type PayloadUTxOuts struct {
	Available bool
	UTxOuts   transactions.UTxOutS
}

type PayloadTxs struct {
	Txs transactions.TxS
}

type PayloadCount struct {
	Count int
}
