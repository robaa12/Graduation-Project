package database

import (
	"fmt"
	"log"
	"order-service/cmd/model"
	"os"
	"time"

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
	// Setup join table
	if err := db.DB.SetupJoinTable(&model.Store{}, "Customers", &model.StoreCustomer{}); err != nil {
		return fmt.Errorf("failed to setup join table: %w", err)
	}

	err := db.DB.AutoMigrate(&model.Order{}, &model.OrderItem{}, &model.Customer{}, &model.Store{}, &model.StoreCustomer{}, &model.OrderStatusHistory{})
	if err != nil {
		return fmt.Errorf("failed to run migration: %w", err)
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
