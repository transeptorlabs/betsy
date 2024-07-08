package data

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func OpenConnection(isInMemory bool) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	if isInMemory {
		db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	} else {
		db, err = gorm.Open(sqlite.Open("betsy.db"), &gorm.Config{})
	}

	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&UserOpV7Hexify{})

	return db, nil
}
