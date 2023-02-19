package db

import (
	"github.com/onee-only/miner-node/lib"
	"github.com/onee-only/miner-node/properties"
	bolt "go.etcd.io/bbolt"
)

var DB *bolt.DB

func Init() {
	s := lib.CreateSpinner(
		"Configuring DB",
		"DB configured!",
	)

	db, err := bolt.Open(properties.DBName, 0600, nil)
	lib.HandleErr(err)
	DB = db

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(properties.BucketBlockchain))
		lib.HandleErr(err)
		_, err = tx.CreateBucketIfNotExists([]byte(properties.BucketCheckpoint))
		return err
	})
	lib.HandleErr(err)

	s.Stop()
	properties.IsDownloading = false
}

func Close() {
	DB.Close()
}
