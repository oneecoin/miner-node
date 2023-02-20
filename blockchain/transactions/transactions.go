package transactions

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"

	"github.com/onee-only/miner-node/lib"
)

func ValidateTx(tx *Tx) bool {

	x, y, err := lib.RestoreBigInts(tx.TxIns.From)
	if err != nil {
		return false
	}

	hash, err := hex.DecodeString(tx.ID)
	if err != nil {
		return false
	}

	if tx.TxIns.From == "COINBASE" {
		return true
	}

	for _, txIn := range tx.TxIns.V {

		// should see if there is actual transaction

		r, s, err := lib.RestoreBigInts(txIn.Signature)
		if err != nil {
			return false
		}

		valid := ecdsa.Verify(&ecdsa.PublicKey{
			Curve: elliptic.P256(),
			X:     x,
			Y:     y,
		}, hash, r, s)
		if !valid {
			return false
		}
	}
	return true
}
