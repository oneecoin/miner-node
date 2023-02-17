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
	checkInterval, err := actor.PromptOptionalAndRetry("Please enter the check interval (m)", "5", lib.CheckInt)
	if err != nil {
		os.Exit(0)
	}
	properties.CheckInterval, _ = strconv.Atoi(checkInterval)

	// finish
	fmt.Printf("\nStarting miner with\n\tport: %d\n\tcheck interval: %dm\n\n", properties.Port, properties.CheckInterval)
}
