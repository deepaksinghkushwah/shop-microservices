package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() error {
	// Try data directory first (for dist deployment)
	dbPath := "data/catalog-service/catalog.db"
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		// Fallback to root (for source tree development)
		dbPath = "catalog.db"
		db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
		if err != nil {
			return err
		}
	}

	DB = db
	return nil
}
