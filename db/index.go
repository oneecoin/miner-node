package db

import (
	"strconv"

	"github.com/onee-only/miner-node/lib"
	"github.com/onee-only/miner-node/properties"
	bolt "go.etcd.io/bbolt"
)

type IndexTx struct {
	From string
	To   string
}

type index struct {
	hash string
	txs  []IndexTx
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
		txs:  []IndexTx{},
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

func AddIndex(height int, hash string, txs []IndexTx) {
	idx := &index{
		hash: hash,
		txs:  txs,
	}
	DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(properties.BucketIndex))
		b.Put([]byte(strconv.Itoa(height)), lib.ToBytes(idx))
		return nil
	})
}
