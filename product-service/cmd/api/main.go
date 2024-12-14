package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/robaa12/product-service/cmd/data"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// WebPort Application Port
const WebPort = "8080"

// Database connection times
var counts int

type Config struct {
	db     *gorm.DB
	models data.Models
}

func main() {
	log.Println("Starting Product Service")

	// load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Panic(err)
	}

	// Connect to Database
	db := connectToDB()
	if db == nil {
		log.Panic("Can't connect to Postgres!")
	}

	// Set up config
	app := Config{
		db:     db,
		models: data.New(db),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", WebPort),
		Handler: app.routes(),
	}
	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
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
