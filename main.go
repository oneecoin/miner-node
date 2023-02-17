package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/deiwin/interact"
	"github.com/mbndr/figlet4go"
	"github.com/onee-only/miner-node/db"
	"github.com/onee-only/miner-node/http"
	"github.com/onee-only/miner-node/lib"
	"github.com/onee-only/miner-node/properties"
	"github.com/onee-only/miner-node/ws/mempool"
	"github.com/onee-only/miner-node/ws/peers"
)

func main() {
	setup()
	go http.InitServer(properties.Port)
	mempool.Connect()
	peers.Connect()
	peers.StartDownloadingBlockChain()
	db.Init()
	defer db.Close()
	mempool.ListenForMining()
}

func setup() {

	// print large letter
	ascii := figlet4go.NewAsciiRender()
	renderStr, _ := ascii.Render("OneeCoin")
	fmt.Print(renderStr)

	// interact
	actor := interact.NewActor(os.Stdin, os.Stdout)
	publicKey, err := actor.PromptAndRetry("Please enter your public key", lib.CheckNotEmpty, lib.CheckHexString)
	if err != nil {
		os.Exit(0)
	}
	properties.PublicKey = publicKey

	port, err := actor.PromptOptionalAndRetry("Please enter the port number", "4000", lib.CheckInt)
	if err != nil {
		os.Exit(0)
	}
	properties.Port, _ = strconv.Atoi(port)

	checkInterval, err := actor.PromptOptionalAndRetry("Please enter the check interval (m) (5 ~ 60)", "5", lib.CheckIntRange5to60)
	if err != nil {
		os.Exit(0)
	}
	properties.CheckInterval, _ = strconv.Atoi(checkInterval)

	minTxs, err := actor.PromptOptionalAndRetry("Please enter the minimum transaction count per block (max 4)", "2", lib.CheckIntRange1to4)
	if err != nil {
		os.Exit(0)
	}
	properties.MinTxs, _ = strconv.Atoi(minTxs)

	// finish
	fmt.Printf(
		"\nStarting miner with\n"+
			"\tport: %d\n"+
			"\tcheck interval: %dm\n"+
			"\tmin transaction count: %d\n"+
			"\n",
		properties.Port, properties.CheckInterval, properties.MinTxs)
}
