package main

import (
	"github.com/onee-only/miner-node/cli"
	"github.com/onee-only/miner-node/config"
	"github.com/onee-only/miner-node/http"
	"github.com/onee-only/miner-node/ws/mempool"
)

func main() {
	cli.Setup()
	go http.InitServer(config.Port)
	mempool.Connect()
}
