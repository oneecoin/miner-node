package properties

type ChanMessageType int

const (
	MessageBlockchainDownloaded ChanMessageType = iota
	MessageBlockchainUploading
	MessageBlockchainUploaded

	MessageNewBlock
)
