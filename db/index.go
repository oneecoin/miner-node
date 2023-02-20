package db

import (
	"strconv"

	"github.com/onee-only/miner-node/lib"
	"github.com/onee-only/miner-node/properties"
	bolt "go.etcd.io/bbolt"
)

type indexTx struct {
	from string
	to   string
}

type index struct {
	hash string
	txs  []indexTx
}

func FindLastHash() string {
	_, cp := findLastIndex()
	return cp.hash
}

func FindCurrentHeight() int {
	height, _ := findLastIndex()
	return height
}

func findLastIndex() (int, *index) {
	idx := &index{
		hash: "",
		txs:  []indexTx{},
	}
	height := 0
	DB.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(properties.BucketIndex)).Cursor()
		k, v := c.Last()
		if k != nil {
			height, _ = strconv.Atoi(string(k))
			lib.FromBytes(idx, v)
		}
		return nil
	})
	return height, idx
}
