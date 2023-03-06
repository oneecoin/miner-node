package db

import (
	"strconv"

	"github.com/onee-only/miner-node/lib"
	"github.com/onee-only/miner-node/properties"
	bolt "go.etcd.io/bbolt"
)

func AddBlock(hash string, block []byte) {
	DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(properties.BucketBlockchain))
		b.Put([]byte(hash), block)
		return nil
	})
}

func FindBlockByHash(hash string) []byte {
	var block []byte
	DB.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(properties.BucketBlockchain)).Cursor()
		if _, v := c.Seek([]byte(hash)); v != nil {
			block = v
		}
		return nil
	})
	return block
}

func FindBlocksPageByHeight(cur int) [][]byte {
	var blocks [][]byte
	DB.View(func(tx *bolt.Tx) error {
		bc := tx.Bucket([]byte(properties.BucketBlockchain)).Cursor()
		ic := tx.Bucket([]byte(properties.BucketIndex)).Cursor()

		cnt := 0

		for {
			if _, idxBytes := ic.Seek([]byte(strconv.Itoa(cur))); idxBytes != nil {
				var idx Index
				lib.FromBytes(&idx, idxBytes)
				_, blockBytes := bc.Seek([]byte(idx.Hash))
				blocks = append(blocks, blockBytes)
			}
			cur--
			cnt++
			if cnt == properties.DefaultPageSize || cur == 0 {
				break
			}
		}
		return nil
	})
	return blocks
}

func FindBlocksByHashes(hashes []string) [][]byte {
	var blocks [][]byte
	DB.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(properties.BucketBlockchain)).Cursor()

		for _, hash := range hashes {
			_, v := c.Seek([]byte(hash))
			if v != nil {
				blocks = append(blocks, v)
			}
		}
		return nil
	})
	return blocks
}
