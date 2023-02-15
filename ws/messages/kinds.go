package messages

type MempoolMessageKind int

const (

	// requests from miner
	MessageTxsRequest MempoolMessageKind = iota

	// responses from miner
	MessageBlocksResponse
	MessageBlockResponse
	MessageUTxOutsResponse

	// requests from mempool
	MessageNewPeer
	MessageBlocksRequest
	MessageBlockRequest
	MessageUTxOutsRequest

	// responses from mempool
	MessageTxsResponse

	// etc.
	MessageRejectPeer
	MessagePeerRejected
)

type PeerMessageKind int

