package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/robaa12/product-service/cmd/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

func New() (*Database, error) {
	db := connectToDB()
	if db == nil {
		return nil, fmt.Errorf("cannot connect to database")
	}
	return &Database{
		DB: db,
	}, nil
}

func (d *Database) SetupDatabase() error {
	// Setup join table
	if err := d.DB.SetupJoinTable(&model.Sku{}, "Variants", &model.SKUVariant{}); err != nil {
		return fmt.Errorf("failed to setup join table: %w", err)
	}

	// Run migrations
	if err := d.DB.AutoMigrate(&model.Store{}, &model.Category{}, &model.Product{}, &model.Sku{}, &model.Variant{}, &model.SKUVariant{}, &model.Collection{}, &model.Review{}); err != nil {
		return fmt.Errorf("failed to run migration: %w", err)
	}

	// Create unique index
	if err := d.DB.Exec(`
			CREATE UNIQUE INDEX IF NOT EXISTS idx_products_store_id_slug
			ON products (store_id, slug)
		`).Error; err != nil {
		return fmt.Errorf("failed to create unique index: %w", err)
	}
	if err := d.DB.Exec(`
			CREATE UNIQUE INDEX IF NOT EXISTS idx_categories_store_id_slug
			ON categories (store_id, slug)
		`).Error; err != nil {
		return fmt.Errorf("failed to create unique index: %w", err)
	}
	if err := d.DB.Exec(`
			CREATE UNIQUE INDEX IF NOT EXISTS idx_collections_store_id_slug
			ON collections (store_id, slug)
		`).Error; err != nil {
		return fmt.Errorf("failed to create unique index: %w", err)
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

func connectToDB() *gorm.DB {
	dsn := os.Getenv("DSN")
	var counts int
	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not yet ready...")
			counts++
		} else {
			log.Println("Connected to Postgres!")
			return connection
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}
		log.Println("Backing off two seconds")
		time.Sleep(2 * time.Second)
		continue
	}
}
