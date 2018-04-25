package db

import (
	"cm-cloud.fr/go_pihole/db/model"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

func init() {
	if db, err := gorm.Open("sqlite3", viper.GetString("db_file")); err != nil {
		panic("Fialed to connect database")
	} else {
		defer db.Close()

		migrateSchema(db)
	}
}

func migrateSchema(db *gorm.DB) {

	db.AutoMigrate(&model.DnsmasqLog{})

}
