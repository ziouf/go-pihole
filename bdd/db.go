package bdd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"reflect"
	"sync"
	"time"

	"github.com/boltdb/bolt"
	"github.com/spf13/viper"
)

var db *bolt.DB
var inserts *buffer
var cleaner *clean

type buffer struct {
	mtx    sync.Mutex
	buffer []Encodable
	timer  *time.Ticker
}

type clean struct {
	buckets []string
	timer   *time.Ticker
}

// Init database
func Init() {
	inserts = &buffer{timer: time.NewTicker(time.Second * viper.GetDuration("db.bulk.freq"))}
	cleaner = &clean{timer: time.NewTicker(time.Second * viper.GetDuration("db.cleaning.freq"))}
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

	// Insert goroutine
	go insertService()

	// DB Cleaning goroutine
	go cleaningService()
}

// Close db and stop routines
func Close() {
	inserts.timer.Stop()
	cleaner.timer.Stop()

	if db != nil {
		db.Close()
	}
}

// Insert append to insertion buffer
func Insert(m Encodable) {
	inserts.mtx.Lock()
	inserts.buffer = append(inserts.buffer, m)
	inserts.mtx.Unlock()

	if len(inserts.buffer) >= viper.GetInt("db.bulk.size") {
		insertBuffer()
	}
}

func insertBuffer() error {
	if len(inserts.buffer) <= 0 {
		return fmt.Errorf("Buffer is empty : nothing to persist")
	}

	inserts.mtx.Lock()
	buffer := inserts.buffer
	inserts.buffer = make([]Encodable, 0)
	inserts.mtx.Unlock()

	t := reflect.TypeOf(buffer[0]).Elem()
	return db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(t.Name()))
		if err != nil {
			return fmt.Errorf("Failed to get/create bucket %s: %s", t.Name(), err)
		}
		for _, b := range buffer {
			enc, err := b.Encode()
			if err != nil {
				return fmt.Errorf("Failed to encode item: %s", err)
			}
			if err := bucket.Put(b.StampEncoded(), enc); err != nil {
				return fmt.Errorf("Failed to insert data: %s", err)
			}
		}
		return nil
	})
}

func insertService() {
	for range inserts.timer.C {
		insertBuffer()
	}
}

// AddToClean append interface to cleaned bucket
func AddToClean(e interface{}) {
	t := reflect.TypeOf(e).Elem()
	cleaner.buckets = append(cleaner.buckets, t.Name())
}

func cleaning() error {
	if db == nil {
		return fmt.Errorf("Db is not open")
	}

	return db.Update(func(tx *bolt.Tx) error {
		for _, b := range cleaner.buckets {
			bucket := tx.Bucket([]byte(b))
			c := bucket.Cursor()
			days := viper.GetDuration("db.cleaning.days.to.keep")
			stamp := []byte(time.Now().Add(-1 * time.Hour * 24 * days).Format(time.RFC3339))

			for k, _ := c.First(); k != nil && bytes.Compare(k, stamp) <= 0; k, _ = c.Next() {
				bucket.Delete(k)
			}
		}
		return nil
	})
}

func cleaningService() {
	for range cleaner.timer.C {
		cleaning()
	}
}
