package messages

type MempoolMessage struct {
	Kind    MempoolMessageKind
	Payload []byte
}

type PeerMessage struct {
	Kind    PeerMessageKind
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
	UTxOuts   transaction_model.UTxOutS
}

type PayloadTxs struct {
	Txs transaction_model.TxS
}
