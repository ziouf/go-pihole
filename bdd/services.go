package bdd

import (
	"encoding/binary"
	"fmt"
	"log"
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
	buckets []Decodable
	timer   *time.Ticker
}

func (c *clean) addBucket(e Decodable) {
	cleaner.buckets = append(cleaner.buckets, e)
}

// Init database
func Init() {
	inserts = &buffer{timer: time.NewTicker(viper.GetDuration("db.bulk.freq"))}
	cleaner = &clean{timer: time.NewTicker(viper.GetDuration("db.cleaning.freq"))}

	cleaner.addBucket(&DNS{})
	cleaner.addBucket(&DHCP{})
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
	inserts.mtx.Lock()
	buffer := inserts.buffer
	inserts.buffer = make([]Encodable, 0)
	inserts.mtx.Unlock()

	if len(buffer) <= 0 {
		return fmt.Errorf("Buffer is empty : nothing to persist")
	}

	t := reflect.TypeOf(buffer[0]).Elem()
	return db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(t.Name()))
		if err != nil {
			return fmt.Errorf("Failed to get/create bucket %s: %s", t.Name(), err)
		}
		for _, b := range buffer {
			id, _ := bucket.NextSequence()
			if err := bucket.Put(itob(id), b.Encode()); err != nil {
				return fmt.Errorf("Failed to insert data: %s", err)
			}
		}
		return nil
	})
}

func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}

func insertService() {
	for range inserts.timer.C {
		if err := insertBuffer(); err != nil {
			log.Println(err)
		}
	}
}

// Cleaning
func cleaning() error {
	if !viper.GetBool(`db.cleaning.enable`) {
		return fmt.Errorf(`Db cleaning service is disabled`)
	}

	if db == nil {
		return fmt.Errorf("Db is not open")
	}

	return db.Update(func(tx *bolt.Tx) error {
		for _, b := range cleaner.buckets {
			t := reflect.TypeOf(b).Elem()
			bucket := tx.Bucket([]byte(t.Name()))
			if bucket == nil {
				continue
			}
			c := bucket.Cursor()
			stamp := time.Now().Add(-1 * viper.GetDuration("db.cleaning.keep"))
			for k, v := c.First(); k != nil; k, v = c.Next() {
				switch t {
				case reflect.TypeOf(&DNS{}).Elem():
					pDNS := new(DNS)
					pDNS.Decode(v)
					if pDNS.Date.Before(stamp) {
						bucket.Delete(k)
					}
				case reflect.TypeOf(&DHCP{}).Elem():
					pDHCP := new(DHCP)
					pDHCP.Decode(v)
					if pDHCP.Date.Before(stamp) {
						bucket.Delete(k)
					}
				default:
				}
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
