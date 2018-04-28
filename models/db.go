package models

import (
	"log"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	// SQlite3 dialect for gORM
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/spf13/viper"
)

const (
	bulkSize     = 5000
	bulkFreq     = 1 * time.Second
	serviceFreq  = 50 * time.Millisecond
	cleaningFreq = 15 * time.Second
	day          = 24 * time.Hour
)

var (
	// Db Database
	Db *gorm.DB

	insertMutx sync.Mutex
	updateMutx sync.Mutex
	deleteMutx sync.Mutex

	tables = []interface{}{
		DhcpLease{},
		DnsmasqLog{},
	}
	tablesToClean = []interface{}{
		DnsmasqLog{},
	}

	insertBuffer = make([]interface{}, 0)
	updateBuffer = make([]interface{}, 0)
	deleteBuffer = make([]interface{}, 0)

	lastInsert = time.Now()
	lastUpdate = time.Now()
	lastDelete = time.Now()
)

// InitDB Initialize database
func InitDB() {
	var err error

	if Db, err = gorm.Open("sqlite3", viper.GetString("db_file")); err != nil {
		log.Fatalln(err)
	}

	Db.Exec("PRAGMA auto_vacuum=FULL;")

	for _, t := range tables {
		// Auto migrate schema
		Db.AutoMigrate(&t)

		if !Db.HasTable(&t) {
			Db.CreateTable(&t)
		}
	}

	go insertService()
	go updateService()
	go deleteService()
	// Init cleaning service
	go startDbCleaningService()
}

// Insert insert data
func Insert(data interface{}) {
	insertMutx.Lock()
	insertBuffer = append(insertBuffer, data)
	insertMutx.Unlock()
	if len(insertBuffer) >= bulkSize {
		insert()
	}
}

func insert() {
	tx := Db.Begin()
	insertMutx.Lock()
	buffer := insertBuffer
	insertBuffer = make([]interface{}, 0)
	insertMutx.Unlock()

	for _, item := range buffer {
		if tx.NewRecord(item) {
			tx.Create(item)
		}
	}
	if tx.Error != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}

	lastInsert = time.Now()
}

func insertService() {
	for {
		if time.Now().Sub(lastInsert) >= bulkFreq {
			insert()
		}
		time.Sleep(serviceFreq)
	}
}

// Update update data
func Update(data interface{}) {
	updateMutx.Lock()
	updateBuffer = append(updateBuffer, data)
	updateMutx.Unlock()
	if len(updateBuffer) >= bulkSize {
		update()
	}
}

func update() {

}

func updateService() {
	for {
		if time.Now().Sub(lastUpdate) >= bulkFreq {
			update()
		}
		time.Sleep(serviceFreq)
	}
}

// Delete delete data
func Delete(data interface{}) {
	deleteMutx.Lock()
	deleteBuffer = append(deleteBuffer, data)
	deleteMutx.Unlock()
	if len(deleteBuffer) >= bulkSize {
		delete()
	}
}

func delete() {

}

func deleteService() {
	for {
		if time.Now().Sub(lastDelete) >= bulkFreq {
			delete()
		}
		time.Sleep(serviceFreq)
	}
}

func startDbCleaningService() {
	time.Sleep(cleaningFreq)
	for {
		for _, t := range tablesToClean {
			tx := Db.Begin()

			tx.Delete(t, "date < ?", time.Now().Add(-1*day*7))

			if tx.Error != nil {
				log.Println(tx.Error)
				tx.Rollback()
			} else {
				tx.Commit()
			}
		}

		time.Sleep(cleaningFreq)
	}
}
