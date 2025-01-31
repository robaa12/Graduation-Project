package main

import (
	"fmt"
	"log"
	"net/http"

	"order-service/cmd/database"

	"gorm.io/gorm"
)

const webPort = "8081"

type Config struct {
	db *gorm.DB
}

func main() {
	log.Println("Starting Order Service...")

	db, err := database.New()
	if err != nil {
		log.Panic(err)
	}
	db.SetupDatabase()
	app := Config{
		db: db.DB,
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
