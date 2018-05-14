package bdd

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"
	"time"

	"cm-cloud.fr/go-pihole/log"
	"github.com/boltdb/bolt"
	"github.com/spf13/viper"
)

// Errors ...
var (
	ErrTickerNil               = errors.New(`Ticker must no be nil`)
	ErrEmptyBuffer             = errors.New(`Buffer is empty : nothing to persist`)
	ErrCleaningServiceDisabled = errors.New(`Db cleaning service is disabled`)
)

type buffer struct {
	mtx    sync.Mutex
	buffer []Encodable
	ticker *time.Ticker
}

func (b *buffer) append(data Encodable) {
	b.mtx.Lock()
	b.buffer = append(b.buffer, data)
	b.mtx.Unlock()

	if len(b.buffer) >= viper.GetInt("db.bulk.size") {
		b.insert()
	}
}

func (b *buffer) start() error {
	if b.ticker == nil {
		return ErrTickerNil
	}
	go func() {
		for range b.ticker.C {
			b.insert()
		}
	}()
	return nil
}

func (b *buffer) insert() error {
	b.mtx.Lock()
	buffer := b.buffer
	b.buffer = make([]Encodable, 0)
	b.mtx.Unlock()

	if len(buffer) <= 0 {
		return ErrEmptyBuffer
	}

	return db.Update(func(tx *bolt.Tx) error {
		for _, b := range buffer {
			t := reflect.TypeOf(b).Elem()
			bucket, err := tx.CreateBucketIfNotExists([]byte(t.Name()))
			if err != nil {
				return fmt.Errorf("Failed to get/create bucket %s: %s", t.Name(), err)
			}
			id, _ := bucket.NextSequence()
			if err := bucket.Put(itob(id), b.Encode()); err != nil {
				return fmt.Errorf("Failed to insert data: %s", err)
			}
		}
		return nil
	})
}

type clean struct {
	buckets []Decodable
	ticker  *time.Ticker
}

func (c *clean) addBucket(e Decodable) {
	log.Debug().Printf("Appending bucket [%s] to auto clean list", reflect.TypeOf(e).Elem().Name())
	c.buckets = append(c.buckets, e)
}

func (c *clean) start() error {
	if c.ticker == nil {
		return ErrTickerNil
	}
	go func() {
		for range c.ticker.C {
			c.clean()
		}
	}()
	return nil
}

func (c *clean) clean() error {
	if !viper.GetBool(`db.cleaning.enable`) {
		return ErrCleaningServiceDisabled
	}

	if db == nil {
		return ErrDbClosed
	}

	return db.Update(func(tx *bolt.Tx) error {
		count := uint64(0)
		for _, b := range c.buckets {
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
						atomic.AddUint64(&count, 1)
					}
				case reflect.TypeOf(&DHCP{}).Elem():
					pDHCP := new(DHCP)
					pDHCP.Decode(v)
					if pDHCP.Date.Before(stamp) {
						bucket.Delete(k)
						atomic.AddUint64(&count, 1)
					}
				default:
				}
			}
		}
		log.Verbose().Println(`Deleted`, count, `item(s)`)
		return nil
	})
}
