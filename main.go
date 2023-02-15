package main

import (
	"github.com/onee-only/miner-node/config"
	"github.com/onee-only/miner-node/http"
)

func main() {
	http.InitServer(config.Port)
}
