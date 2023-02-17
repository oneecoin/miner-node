package peers

import "github.com/onee-only/miner-node/ws/messages"

func (*TPeers) handleMessage(m *messages.Message, p *Peer) {
	switch m.Kind {
	case messages.MessageDownloadRequest:
		uploadBlockchain(p)
	}
}
