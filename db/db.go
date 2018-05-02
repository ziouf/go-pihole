package db

import (
	"log"
	"reflect"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	// SQlite3 dialect for gORM
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/spf13/viper"
)

const serviceFreq = 1 * time.Millisecond

var (
	// Db Database
	Db *gorm.DB

	tablesToAutoClean = make([]interface{}, 0)

	insertMutx sync.Mutex
	updateMutx sync.Mutex
	deleteMutx sync.Mutex

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

	if Db == nil {
		if Db, err = gorm.Open("sqlite3", viper.GetString("db_file")); err != nil {
			log.Fatalln(err)
		}
	}

	// Debug config
	// Db.LogMode(true)

	// PRAGMA configs
	Db.Exec("PRAGMA auto_vacuum=FULL;")

	// DB Configuration
	Db.DB().SetMaxIdleConns(1)
	Db.DB().SetMaxOpenConns(1)

}

func InitDbServices() {
	// Init background services
	go insertService()
	go updateService()
	go deleteService()
	// Init cleaning service
	go startDbCleaningService()
}

// InitDataModel Initialize table
func InitDataModel(model interface{}) {
	t := reflect.TypeOf(model)

	if !Db.HasTable(model) {
		Db.CreateTable(model)
		log.Printf("[%s] Create table", t.Name())
	}
	Db.AutoMigrate(model)

	if Db.Error != nil {
		log.Fatalln(Db.Error)
	}
	log.Printf("[%s] Init of data model done", t.Name())
}

// AutoCleanTable Add table to autoclean list
func AutoCleanTable(model interface{}) {
	tablesToAutoClean = append(tablesToAutoClean, model)
	log.Printf("[%s] Added table to auto clean", reflect.TypeOf(model).Name())
}

// Insert insert data
func Insert(data interface{}) {
	insertMutx.Lock()
	insertBuffer = append(insertBuffer, data)
	insertMutx.Unlock()
	if len(insertBuffer) >= viper.GetInt("db.bulk.size") {
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
		log.Println(tx.Error)
		tx.Rollback()
	} else {
		tx.Commit()
	}

	lastInsert = time.Now()
	// log.Println("Insert data")
}

func insertService() {
	for {
		if time.Now().Sub(lastInsert) >= viper.GetDuration("db.bulk.freq")*time.Second {
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
	if len(updateBuffer) >= viper.GetInt("db.bulk.size") {
		update()
	}
}

func update() {
	tx := Db.Begin()
	insertMutx.Lock()
	buffer := updateBuffer
	insertBuffer = make([]interface{}, 0)
	insertMutx.Unlock()

	for _, item := range buffer {
		tx.First(&item)
		tx.Save(&item)
	}

	if tx.Error != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}
	lastUpdate = time.Now()
}

func updateService() {
	for {
		if time.Now().Sub(lastUpdate) >= viper.GetDuration("db.bulk.freq")*time.Second {
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
	if len(deleteBuffer) >= viper.GetInt("db.bulk.size") {
		delete()
	}
}

func delete() {

}

func deleteService() {
	for {
		if time.Now().Sub(lastDelete) >= viper.GetDuration("db.bulk.freq")*time.Second {
			delete()
		}
		time.Sleep(serviceFreq)
	}
}

func startDbCleaningService() {
	time.Sleep(viper.GetDuration("db.cleaning.freq") * time.Second)
	for {
		cleanDb()

		time.Sleep(viper.GetDuration("db.cleaning.freq") * time.Second)
	}
}

func cleanDb() {
	for _, t := range tablesToAutoClean {
		tx := Db.Begin()

		tx.Unscoped().Delete(t, "date < ?", time.Now().Add(-1*time.Hour*24*7))

		if tx.Error != nil {
			log.Println(tx.Error)
			tx.Rollback()
		} else {
			// log.Println("Cleaning db")
			tx.Commit()
		}
	}
}
