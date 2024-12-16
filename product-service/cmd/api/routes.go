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
	mux.Get("/products/{id}/details", app.GetProductDetails)
	mux.Put("/products/skus/{id}", app.UpdateSKU)
	mux.Get("/products/skus/{id}", app.GetSKU)
	mux.Delete("/products/skus/{id}", app.DeleteSKU)
	mux.Post("/products/{id}/skus", app.NewSKU)

	return mux
}

// To do endpoints
// 1. Get product details by ID (Done)
// 2. Get sku details by ID (Done)
// 2. Delete SKU by ID (Done)
// 3. Edit SKU by ID  (Done)
// 4. Create SKU for a product by ID
