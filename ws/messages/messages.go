package messages

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
	// UTxOuts   transaction_model.UTxOutS
}

type PayloadTxs struct {
	// Txs transaction_model.TxS
}
