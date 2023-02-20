package db

import (
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

func FindBlocksPageByHash(hash string) [][]byte {
	var blocks [][]byte
	DB.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(properties.BucketBlockchain)).Cursor()

		cnt := 0
		if _, v := c.Seek([]byte(hash)); v != nil {
			blocks = append(blocks, v)
		}
		for {
			cnt++
			if cnt == properties.DefaultPageSize {
				break
			}
			_, v := c.Prev()
			if v == nil {
				break
			}
			blocks = append(blocks, v)
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
