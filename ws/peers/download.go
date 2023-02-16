package peers

import (
	"encoding/json"

	"github.com/onee-only/miner-node/ws/messages"
)

func StartDownloadingBlockChain() {
	peer := getRandomPeer()

	m := messages.Message{
		Kind:    messages.MessageDownloadRequest,
		Payload: nil,
	}

	bytes, err := json.Marshal(m)
	peer.Inbox <- bytes
}
