package main

import (
	"github.com/onee-only/miner-node/cli"
	"github.com/onee-only/miner-node/http"
	"github.com/onee-only/miner-node/properties"
	"github.com/onee-only/miner-node/ws/mempool"
	"github.com/onee-only/miner-node/ws/peers"
)

func main() {
	cli.Setup()
	go http.InitServer(properties.Port)
	mempool.Connect()
	peers.Connect()
	peers.StartDownloadingBlockChain()
	blockchain.InitDB()
	mempool.ListenForMining()
}
