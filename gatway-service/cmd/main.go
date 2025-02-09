package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/robaa12/gatway-service/api"
	"github.com/robaa12/gatway-service/config"
	"github.com/robaa12/gatway-service/internal/auth"
	"github.com/robaa12/gatway-service/internal/proxy"
)

type Application struct {
	config     *config.Config
	auth       *auth.AuthService
	proxy      *proxy.ProxyService
	httpServer *http.Server
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// App instance
	app := &Application{
		config: cfg,
	}

	// Setup services
	err = app.setupServices()
	if err != nil {
		log.Fatalf("Failed to setup services: %v", err)
	}

	// Start server
	go func() {
		app.startServer()
	}()

	// Wait for shutdown signal
	app.waitForShutdown()
}

func (app *Application) setupServices() error {
	app.auth = auth.NewAuthService(app.config)
	app.proxy = proxy.NewProxyService(app.config)

	routes := api.Routes(app.auth, app.proxy)

	app.httpServer = &http.Server{
		Addr:         fmt.Sprintf("%s:%s", app.config.Server.Host, app.config.Server.Port),
		Handler:      routes,
		IdleTimeout:  15 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	return nil
}

func (app *Application) startServer() {
	log.Printf("Starting server on %s:%s\n", app.config.Server.Host, app.config.Server.Port)
	err := app.httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func (app *Application) waitForShutdown() {
	quit := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-quit
		log.Printf("Received shutdown signal: %v\n", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		log.Println("Shutting down server...")
		if err := app.httpServer.Shutdown(ctx); err != nil {
			log.Printf("Server forced to shutdown: %v\n", err)
		}
		done <- true
	}()
	<-done
	log.Println("Server exited properly")
}
