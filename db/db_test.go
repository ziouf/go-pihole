package db

import (
	"testing"
	"time"
)

type TestTable struct {
	Date time.Time
}

func TestDB(t *testing.T) {
	t.Run("Database init test", testInitDb)
	t.Run("Database init datamodel test", testInitDatamodel)
	t.Run("Database cleaning test", testDbCleaning)
}

func testInitDb(t *testing.T) {

	if Db != nil {
		t.Errorf("Db must be nil befor init")
	}

	InitDB()

	if Db == nil {
		t.Errorf("Failed to Init db")
	}

}

func testInitDatamodel(t *testing.T) {
	InitDB()
	InitDataModel(&TestTable{})

	if !Db.HasTable(&TestTable{}) {
		t.Errorf("Failed to create table")
	}

}

func testDbCleaning(t *testing.T) {
	var count = 0

	// Init DB
	InitDB()
	InitDataModel(TestTable{})
	AutoCleanTable(TestTable{})

	for i := time.Duration(14); i >= 0; i-- {
		Db.Create(TestTable{Date: time.Now().Add(-1 * time.Hour * 24 * i)})
	}

	Db.Model(&TestTable{}).Count(&count)

	if count != 15 {
		t.Errorf("Error while test db init [count : %d]", count)
	}

	// Run service
	cleanDb()

	Db.Model(&TestTable{}).Count(&count)

	if count != 7 {
		t.Errorf("Failed to clean db [count : %d]", count)
	}

	Db.DropTable(TestTable{})

}
