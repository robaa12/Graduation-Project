package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/robaa12/product-service/cmd/database"
	"github.com/robaa12/product-service/cmd/model"
)

// WebPort Application Port
const WebPort = "8083"

// Database connection times

type Config struct {
	db     *database.Database
	models model.Models
}

func main() {
	log.Printf("Starting Product Service On Port %s...\n", WebPort)

	// Setup Database (migrations, indexes)
	DB, err := database.New()
	if err != nil {
		log.Panic(err)
	}
	err = DB.SetupDatabase()
	if err != nil {
		return
	}

	// Set up config
	app := Config{
		db:     DB,
		models: model.New(DB.DB),
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
