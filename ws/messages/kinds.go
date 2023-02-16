package messages

type MessageKind int

const (

	// mempool-node

	// requests from miner
	MessageTxsRequest MessageKind = iota

	// responses from miner
	MessageBlocksResponse
	MessageBlockResponse
	MessageUTxOutsResponse

	// requests from mempool
	MessageBlocksRequest
	MessageBlockRequest
	MessageUTxOutsRequest

	// responses from mempool
	MessageTxsResponse

	// etc.
	MessageRejectPeer
	MessagePeerRejected

	// node-node

	// download blockchain
	MessageDownloadRequest

	MessageNewBlock
)
