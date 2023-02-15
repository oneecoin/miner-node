package main

import (
	"github.com/onee-only/miner-node/cli"
	"github.com/onee-only/miner-node/config"
	"github.com/onee-only/miner-node/http"
)

func main() {
	cli.Setup()
	http.InitServer(config.Port)
}
