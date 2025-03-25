package main

import (
	"fmt"
	"log"
	"net/http"

	"order-service/cmd/database"

	"gorm.io/gorm"
)

const webPort = "8084"

type Config struct {
	db *gorm.DB
}

func main() {
	log.Printf("Starting Order Service On Port %s...\n", webPort)

	db, err := database.New()
	if err != nil {
		log.Panic(err)
	}
	err = db.SetupDatabase()
	if err != nil {
		return
	}
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
