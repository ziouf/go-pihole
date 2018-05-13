package bdd

import (
	"sync/atomic"
	"fmt"
	"reflect"

	"github.com/boltdb/bolt"
)

func getBucket(tx *bolt.Tx, name string) (*bolt.Bucket, error) {
	bucket := tx.Bucket([]byte(name))
	if bucket == nil {
		return bucket, fmt.Errorf("Bucket [%s] not found", name)
	}
	return bucket, nil
}

func GetLast(d Decodable) (Decodable, error) {
	t := reflect.TypeOf(d).Elem()
	err := db.View(func(tx *bolt.Tx) error {
		bucket, e := getBucket(tx, t.Name())
		if e != nil {
			return e
		}
		c := bucket.Cursor()
		_, v := c.Last()
		if v == nil {
			return fmt.Errorf("Entry not found")
		}
		if e := d.Decode(v); e != nil {
			return fmt.Errorf("Can't decode value: %s", e)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return d, nil
}

func Count(d Decodable) (uint64, error) {
	t := reflect.TypeOf(d).Elem()
	count := uint64(0)
	err := db.View(func(tx *bolt.Tx) error {
		bucket, e := getBucket(tx, t.Name())
		if e != nil {
			return e
		}
		c := bucket.Cursor()
		for k,_ := c.First(); k != nil; k,_ = c.Next() {
			atomic.AddUint64(&count, 1)
		}

		return nil
	})
	return count, err
}
