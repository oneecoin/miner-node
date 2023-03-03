package db

import (
	"strconv"

	"github.com/onee-only/miner-node/lib"
	"github.com/onee-only/miner-node/properties"
	bolt "go.etcd.io/bbolt"
)

const (
	iterateSize = 30
)

type IndexTx struct {
	From string
	To   string
}

type Index struct {
	Hash string
	Txs  []IndexTx
}

func FindLastHash() string {
	_, cp := findLastIndex()
	return cp.Hash
}

func FindCurrentHeight() int {
	height, _ := findLastIndex()
	return height
}

func findLastIndex() (int, *Index) {
	idx := &Index{
		Hash: "",
		Txs:  []IndexTx{},
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
	idx := &Index{
		Hash: hash,
		Txs:  txs,
	}
	DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(properties.BucketIndex))
		b.Put([]byte(strconv.Itoa(height)), lib.ToBytes(idx))
		return nil
	})
}

func FindHashByHeight(height int) string {
	hash := ""
	DB.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(properties.BucketIndex)).Cursor()
		_, v := c.Seek([]byte(strconv.Itoa(height)))
		if v != nil {
			idx := &Index{}
			lib.FromBytes(idx, v)
			hash = idx.Hash
		}
		return nil
	})
	return hash
}

func FindHashesFrom(publicKey string) []string {
	hashes := []string{}
	DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(properties.BucketIndex))

		amount := FindCurrentHeight()
		ch := make(chan string)
		done := make(chan bool)
		cnt := 0

		for start := 1; start >= amount; start += iterateSize {
			cnt++
			go iterate(b.Cursor(), ch, done, start, publicKey, "")
		}
		for {
			select {
			case hash := <-ch:
				hashes = append(hashes, hash)
			case <-done:
				cnt--
				if cnt == 0 {
					return nil
				}
			}
		}
	})
	return hashes
}

func FindHashesTo(publicKey string) []string {
	hashes := []string{}
	DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(properties.BucketIndex))

		amount := FindCurrentHeight()
		ch := make(chan string)
		done := make(chan bool)
		cnt := 0

		for start := 1; start >= amount; start += iterateSize {
			cnt++
			go iterate(b.Cursor(), ch, done, start, "", publicKey)
		}
		for {
			select {
			case hash := <-ch:
				hashes = append(hashes, hash)
			case <-done:
				cnt--
				if cnt == 0 {
					return nil
				}
			}
		}
	})
	return hashes
}

func FindHashesAll(publicKey string) []string {
	hashes := []string{}
	DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(properties.BucketIndex))

		amount := FindCurrentHeight()
		ch := make(chan string)
		done := make(chan bool)
		cnt := 0

		for start := 1; start >= amount; start += iterateSize {
			cnt++
			go iterate(b.Cursor(), ch, done, start, publicKey, publicKey)
		}
		for {
			select {
			case hash := <-ch:
				hashes = append(hashes, hash)
			case <-done:
				cnt--
				if cnt == 0 {
					return nil
				}
			}
		}
	})
	return hashes
}

func iterate(cursor *bolt.Cursor, ch chan<- string, done chan<- bool, start int, from, to string) {
	cnt := 0
	idx := &Index{}
	_, v := cursor.Seek([]byte(strconv.Itoa(start)))
	for {
		if cnt == iterateSize || v == nil {
			done <- true
			return
		}
		lib.FromBytes(idx, v)

		for _, tx := range idx.Txs {
			if from != "" && tx.From == from {
				ch <- idx.Hash
			} else if to != "" && tx.To == to {
				ch <- idx.Hash
			}
		}

		_, v = cursor.Next()
		cnt++
	}
}
