package bdd

import (
	"log"
	"os"

	"github.com/boltdb/bolt"
	"github.com/spf13/viper"
)

var db *bolt.DB

// Open and Init db
func Open() {
	var err error
	dbFile, dbFileMode := viper.GetString("db.file.path"), viper.GetInt("db.file.mode")
	options := bolt.Options{ReadOnly: false}
	db, err = bolt.Open(dbFile, os.FileMode(dbFileMode), &options)
	if err != nil {
		log.Fatal(err)
	}

	// DB Config
	db.MaxBatchSize = viper.GetInt("db.bulk.size")
	db.MaxBatchDelay = viper.GetDuration("db.cleaning.freq")

	// Insert goroutine
	go insertService()

	// DB Cleaning goroutine
	go cleaningService()
}

// Close db and stop routines
func Close() {
	stopServices()

	if len(inserts.buffer) > 0 {
		insertBuffer()
	}

	if db != nil {
		db.Close()
	}
}
