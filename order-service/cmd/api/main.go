package main

import (
	"fmt"
	"log"
	"net/http"
	"order-service/cmd/data"
	"order-service/cmd/database"
)

const webPort = "8081"

type Config struct {
	db     *database.Database
	models data.Models
}

func main() {
	log.Println("Starting Order Service...")

	db, err := database.New()
	if err != nil {
		log.Panic(err)
	}
	db.SetupDatabase()
	app := Config{
		db:     db,
		models: data.New(),
	}
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Panic(err)

	}

}
