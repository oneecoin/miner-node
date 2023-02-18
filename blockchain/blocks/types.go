package blocks

import "github.com/onee-only/miner-node/blockchain/transactions"

type Block struct {
	Hash         string             `json:"hash"`
	PrevHash     string             `json:"prevHash,omitempty"`
	Height       int                `json:"height"`
	Difficulty   int                `json:"difficulty"`
	Nonce        int                `json:"nonce"`
	Timestamp    int                `json:"timestamp"`
	Transactions []*transactions.Tx `json:"transactions"`
}

