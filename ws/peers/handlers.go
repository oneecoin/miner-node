package peers

import (
	"github.com/onee-only/miner-node/blockchain/blocks"
	"github.com/onee-only/miner-node/lib"
	"github.com/onee-only/miner-node/properties"
	"github.com/onee-only/miner-node/ws/messages"
)

func (*TPeers) handleMessage(m *messages.Message, p *Peer) {
	switch m.Kind {
	case messages.MessageDownloadRequest:
		uploadBlockchain(p)
		blocks.SaveBroadcastedBlocks()
	case messages.MessageNewBlock:
		var block *blocks.Block
		lib.FromBytes(block, m.Payload)

		if valid := blocks.ValidateBlock(block); valid {
			properties.C <- properties.MessageNewBlock
			properties.BlockReceiveInbox <- m.Payload
			properties.NewBlockInbox <- p.GetAddress()
		} else {
			properties.RejectPeerInbox <- p.GetAddress()
		}
	}
}
