package database

import (
	"github.com/deepaksinghkushwah/shop-microservices/pkg/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	_ "modernc.org/sqlite"
)

var DB *gorm.DB

func Connect() error {
	// Allow overriding the DB path via env. Useful for dist builds where each service has its own data directory.
	dbPath := config.GetEnv("DB_PATH", "data/order.db")

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		// Fallback to root (for source tree development)
		dbPath = "order.db"
		db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
		if err != nil {
			return err
		}
	}

	DB = db
	return nil
}
