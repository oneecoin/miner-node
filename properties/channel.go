package properties

type ChanMessageType int

const (
	MessageBlockchainDownloaded ChanMessageType = iota
	MessageBlockchainUploading
	MessageBlockchainUploaded

	MessageNewBlock
)

var (
	C                   = make(chan ChanMessageType)
	BlockReceiveInbox   = make(chan []byte)
	BlockBroadcastInbox = make(chan []byte)
	NewBlockInbox       = make(chan string)
	RejectPeerInbox     = make(chan string)
	PeerRejectedInbox   = make(chan string)
)
