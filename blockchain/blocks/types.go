package blocks

import "github.com/onee-only/miner-node/blockchain/transactions"

type Block struct {
	Hash         string           `json:"hash"`
	PrevHash     string           `json:"prevHash,omitempty"`
	Height       int              `json:"height"`
	Nonce        int              `json:"nonce"`
	Timestamp    int              `json:"timestamp"`
	Transactions transactions.TxS `json:"transactions"`
}

type BlockSummary struct {
	Hash              string `json:"hash"`
	Height            int    `json:"height"`
	Timestamp         int    `json:"timestamp"`
	TransactionsCount int    `json:"transactionsCount"`
}
