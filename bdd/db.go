package bdd

import (
	"encoding/binary"
	"log"
	"os"
	"time"

	"github.com/boltdb/bolt"
	"github.com/spf13/viper"
)

var db *bolt.DB
var inserts *buffer
var cleaner *clean

// Init database
func Init() {
	inserts = &buffer{ticker: time.NewTicker(viper.GetDuration("db.bulk.freq"))}
	cleaner = &clean{ticker: time.NewTicker(viper.GetDuration("db.cleaning.freq"))}

	cleaner.addBucket(&DNS{})
	cleaner.addBucket(&DHCP{})

	inserts.start()
	cleaner.start()
}

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
}

// Close db and stop routines
func Close() {
	stopServices()

	if len(inserts.buffer) > 0 {
		inserts.insert()
	}

	if db != nil {
		db.Close()
	}
}

func stopServices() {
	inserts.ticker.Stop()
	cleaner.ticker.Stop()
}

// Insert append to insertion buffer
func Insert(m Serializable) {
	inserts.append(m)
}

func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}
