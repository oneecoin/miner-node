package db

import (
	"time"

	"github.com/briandowns/spinner"
	"github.com/onee-only/miner-node/lib"
	"github.com/onee-only/miner-node/properties"
	bolt "go.etcd.io/bbolt"
)

var DB *bolt.DB

func Init() {
	// spinner
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Prefix = "Configuring DB "
	s.FinalMSG = "DB configured!"
	s.Start()

	// actual work
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

	// if new block broadcasted while downloading
	// retrieve all blocks from queue
	// and save it
}

func Close() {
	DB.Close()
}
