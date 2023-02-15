package cli

import (
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/deiwin/interact"
	"github.com/mbndr/figlet4go"
	"github.com/onee-only/miner-node/config"
)

func checkHexString(input string) error {
	_, err := hex.DecodeString(input)
	if err != nil {
		return errors.New(config.ErrorStr("input should be hex-encoded"))
	}
	return nil
}

func checkInt(input string) error {
	_, err := strconv.Atoi(input)
	if err != nil {
		return errors.New(config.ErrorStr("input should be integer"))
	}
	return nil
}

func Setup() {

	// print large letter
	ascii := figlet4go.NewAsciiRender()
	renderStr, _ := ascii.Render("OneeCoin")
	fmt.Print(renderStr)

	// interact
	actor := interact.NewActor(os.Stdin, os.Stdout)
	publicKey, err := actor.PromptAndRetry("Please enter your public key", checkHexString)
	if err != nil {
		os.Exit(0)
	}
	config.PublicKey = publicKey
	port, err := actor.PromptOptionalAndRetry("Please enter the port number", "4000", checkInt)
	if err != nil {
		os.Exit(0)
	}
	config.Port, _ = strconv.Atoi(port)
	
}
