package messages

type MessageKind int

const (

	// mempool-node

	// requests from miner
	MessageMempoolTxsRequest MessageKind = iota

	// responses from miner
	MessageBlocksResponse
	MessageBlockResponse
	MessageUTxOutsResponse

	// requests from mempool
	MessageBlocksRequest
	MessageBlockRequest
	MessageUTxOutsRequest

	// responses from mempool
	MessageTxsMempoolResponse
	MessageTxsDeclined

	// etc.
	MessageRejectPeer
	MessagePeerRejected
	MessageBlockAdded
	MessageNewBlock

	// node-node

	// download blockchain
	MessageDownloadRequest

)
