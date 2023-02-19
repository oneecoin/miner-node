package db

import (
	"strconv"

	"github.com/onee-only/miner-node/lib"
	"github.com/onee-only/miner-node/properties"
	bolt "go.etcd.io/bbolt"
)

type checkpoint struct {
	hash     string
	prevhash string
}

func FindLastHash() string {
	_, cp := findLatestCheckpoint()
	return cp.hash
}

func FindCurrentHeight() int {
	height, _ := findLatestCheckpoint()
	return height
}

func findLatestCheckpoint() (int, *checkpoint) {
	cp := &checkpoint{
		hash:     "",
		prevhash: "",
	}
	height := 0
	DB.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(properties.BucketCheckpoint)).Cursor()
		k, v := c.Last()
		if k != nil {
			height, _ = strconv.Atoi(string(k))
			lib.FromBytes(cp, v)
		}
		return nil
	})
	return height, cp
}
