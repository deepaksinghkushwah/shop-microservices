package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	_ "modernc.org/sqlite"
)

var DB *gorm.DB

func Connect() error {
	// Try data directory first (for dist deployment)
	dbPath := "data/auth-service/auth.db"
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		// Fallback to services directory (for source tree development)
		dbPath = "services/auth-service/auth.db"
		db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
		if err != nil {
			return err
		}
	}

	DB = db
	return nil
}
