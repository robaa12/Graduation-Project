package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Use(middleware.Heartbeat("/ping"))
	mux.Post("/products", app.NewProduct)
	mux.Get("/products/{id}", app.GetProduct)
	mux.Put("/products/{id}", app.UpdateProduct)
	mux.Delete("/products/{id}", app.DeleteProduct)
	mux.Get("/stores/{store_id}/products", app.GetStoreProducts)
	mux.Get("products/{id}/details", app.GetProductDetails)

	return mux
}
