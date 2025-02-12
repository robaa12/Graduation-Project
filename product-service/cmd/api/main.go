package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/robaa12/product-service/cmd/data"
	"github.com/robaa12/product-service/cmd/db"
)

// WebPort Application Port
const WebPort = "8083"

// Database connection times
var counts int

type Config struct {
	db     *db.Database
	models data.Models
}

func main() {
	log.Printf("Starting Product Service On Port %s...\n", WebPort)

	// load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Panic(err)
	}

	// Setup Database (migrations, indexes)
	database, err := db.New()
	database.SetupDatabase()

	// Set up config
	app := Config{
		db:     database,
		models: data.New(database.DB),
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
