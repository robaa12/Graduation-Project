package database

import (
	"fmt"
	"log"
	"order-service/cmd/data"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

func New() (*Database, error) {
	db := connectToDb()
	if db == nil {
		return nil, fmt.Errorf("can't connect to postgres")
	}
	return &Database{DB: db}, nil

}
func (db *Database) SetupDatabase() error {
	err := db.DB.AutoMigrate(&data.Customer{}, &data.Order{}, &data.OrderItem{})
	if err != nil {
		return fmt.Errorf("failed to run migration: %w", err)
	}
	err = db.DB.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS idx_customers_email 
	ON customers (email)`).Error
	if err != nil {
		return fmt.Errorf("failed to create unique index %w", err)
	}
	return nil

}

func openDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func connectToDb() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		log.Panic(err)
	}
	dsn := os.Getenv("DSN")
	var count int
	for {
		db, err := openDB(dsn)

		if err != nil {
			log.Println("Postgres not yet ready...")
			count++
		} else {
			log.Println("Connected to Postgres!")
			return db
		}
		if count > 9 {
			log.Println(err)
			return nil
		}

		log.Println("Backing off two seconds")
		time.Sleep(2 * time.Second)

	}

}
