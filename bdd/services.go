package bdd

import (
	"bytes"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/boltdb/bolt"
	"github.com/spf13/viper"
)

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

func (c *clean) add(e interface{}) {
	t := reflect.TypeOf(e).Elem()
	cleaner.buckets = append(cleaner.buckets, t.Name())
}

// Init database
func Init() {
	inserts = &buffer{timer: time.NewTicker(viper.GetDuration("db.bulk.freq"))}
	cleaner = &clean{timer: time.NewTicker(viper.GetDuration("db.cleaning.freq"))}

	cleaner.add(&DNS{})
	cleaner.add(&DHCP{})
}

func stopServices() {
	inserts.timer.Stop()
	cleaner.timer.Stop()
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
			if err := bucket.Put(b.Encode()); err != nil {
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

// Cleaning
func cleaning() error {
	if db == nil {
		return fmt.Errorf("Db is not open")
	}

	return db.Update(func(tx *bolt.Tx) error {
		for _, b := range cleaner.buckets {
			bucket := tx.Bucket([]byte(b))
			if bucket == nil {
				continue
			}
			c := bucket.Cursor()
			stamp := time.Now().Add(-1 * viper.GetDuration("db.cleaning.keep"))
			for k, _ := c.First(); k != nil && bytes.Compare(k, []byte(stamp.Format(time.Stamp))) <= 0; k, _ = c.Next() {
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
