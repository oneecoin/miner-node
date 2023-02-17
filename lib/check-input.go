package lib

import (
	"encoding/hex"
	"errors"
	"strconv"

	"github.com/onee-only/miner-node/properties"
)

func CheckHexString(input string) error {
	_, err := hex.DecodeString(input)
	if err != nil {
		return errors.New(properties.ErrorStr("input should be hex-encoded"))
	}
	return nil
}

func CheckNotEmpty(input string) error {
	if input == "" {
		return errors.New(properties.ErrorStr("input should not be empty"))
	}
	return nil
}

func CheckInt(input string) error {
	_, err := strconv.Atoi(input)
	if err != nil {
		return errors.New(properties.ErrorStr("input should be integer"))
	}
	return nil
}
