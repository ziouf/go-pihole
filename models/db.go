package models

import (
	"log"

	"github.com/jinzhu/gorm"
	// SQlite3 dialect for gORM
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/spf13/viper"
)

// Db Database
var Db *gorm.DB

// InitDB Initialize database
func InitDB() {
	var err error

	if Db, err = gorm.Open("sqlite3", viper.GetString("db_file")); err != nil {
		log.Fatalln(err)
		// panic("Fialed to connect database")
	}

	// Auto migrate schema
	Db.AutoMigrate(&DnsmasqLog{})

	if !Db.HasTable(&DnsmasqLog{}) {
		Db.CreateTable(&DnsmasqLog{})
	}

}
