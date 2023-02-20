package messages

type MessageKind int

const (

	// mempool-node

	// requests from miner
	MessageMempoolTxsRequest MessageKind = iota
	MessageInvalidTxsRequest

	// responses from miner
	MessageBlocksResponse
	MessageBlockResponse
	MessageUTxOutsResponse
	MessageNodeTxsResponse

	// requests from mempool
	MessageBlocksRequest
	MessageBlockRequest
	MessageUTxOutsRequest
	MessageNodeTxsRequest

	// responses from mempool
	MessageTxsMempoolResponse
	MessageTxsDeclined

	// etc.
	MessageRejectPeer
	MessagePeerRejected
	MessageNewBlock

	// node-node

	// download blockchain
	MessageDownloadRequest
)
