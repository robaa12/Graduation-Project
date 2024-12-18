package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/robaa12/product-service/cmd/api/handlers"
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

	productHandler := handlers.ProductHandler{DB: app.db}
	skuHandler := handlers.SKUHandler{DB: app.db}

	mux.Use(middleware.Heartbeat("/ping"))
	mux.Post("/products", productHandler.NewProduct)
	mux.Get("/products/{id}", productHandler.GetProduct)
	mux.Put("/products/{id}", productHandler.UpdateProduct)
	mux.Delete("/products/{id}", productHandler.DeleteProduct)
	mux.Get("/stores/{store_id}/products", productHandler.GetStoreProducts)
	mux.Get("/products/{id}/details", productHandler.GetProductDetails)
	mux.Put("/products/skus/{id}", skuHandler.UpdateSKU)
	mux.Get("/products/{productID}/skus/{id}", skuHandler.GetSKU)
	mux.Delete("/products/skus/{id}", skuHandler.DeleteSKU)
	mux.Post("/products/{id}/skus", skuHandler.NewSKU)

	return mux
}
