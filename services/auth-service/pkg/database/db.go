package database

import (
	"fmt"
	"log"

	"github.com/deepaksinghkushwah/shop-microservices/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() error {
	host := config.GetEnv("DB_HOST", "localhost")
	port := config.GetEnv("DB_PORT", "5432")
	user := config.GetEnv("DB_USER", "postgres")
	password := config.GetEnv("DB_PASSWORD", "password")
	dbname := config.GetEnv("DB_NAME", "auth")
	sslmode := config.GetEnv("DB_SSLMODE", "disable")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, dbname, port, sslmode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	// Optional: connection pool tuning
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)

	DB = db
	log.Println("Connected to PostgreSQL")
	return nil
}
