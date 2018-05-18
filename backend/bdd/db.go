package bdd

import (
	"encoding/binary"
	"errors"
	"os"
	"time"

	"cm-cloud.fr/go-pihole/backend/log"
	"github.com/boltdb/bolt"
	"github.com/spf13/viper"
)

var db *bolt.DB
var inserts *buffer
var cleaner *clean

// Errors
var (
	ErrDbClosed = errors.New("Database is closed")
)

// Init database
func Init() {
	log.Debug().Println("Initing database")
	inserts = &buffer{ticker: time.NewTicker(viper.GetDuration("db.bulk.freq"))}
	cleaner = &clean{ticker: time.NewTicker(viper.GetDuration("db.cleaning.freq"))}

	cleaner.addBucket(&DNS{})
	cleaner.addBucket(&DHCP{})

	log.Debug().Println("Starting database services")
	inserts.start()
	cleaner.start()
}

// Open and Init db
func Open() {
	log.Debug().Println("Openning database")
	var err error
	dbFile, dbFileMode := viper.GetString("db.file.path"), viper.GetInt("db.file.mode")
	options := bolt.Options{ReadOnly: false}
	db, err = bolt.Open(dbFile, os.FileMode(dbFileMode), &options)
	if err != nil {
		log.Error().Fatal(err)
	}

	// DB Config
	db.MaxBatchSize = viper.GetInt("db.bulk.size")
	db.MaxBatchDelay = viper.GetDuration("db.cleaning.freq")
}

// Close db and stop routines
func Close() {
	log.Debug().Println("Closing database")
	stopServices()

	if len(inserts.buffer) > 0 {
		inserts.insert()
	}

	if db != nil {
		db.Close()
	}
}

func stopServices() {
	log.Debug().Println("Stopping database services")
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
