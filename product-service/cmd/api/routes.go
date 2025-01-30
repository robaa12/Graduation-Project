package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/robaa12/product-service/cmd/api/handlers"
	"github.com/robaa12/product-service/cmd/api/middleware"
)

func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()

	// Middleware
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Initialize handlers
	productHandler := handlers.ProductHandler{DB: app.db.DB}
	skuHandler := handlers.SKUHandler{DB: app.db.DB}

	// Product routes
	mux.Route("/products", func(r chi.Router) {
		// Public endpoints
		r.Get("/{id}", productHandler.GetProduct)
		r.Get("/{id}/details", productHandler.GetProductDetails)
		r.Get("/{productID}/skus/{id}", skuHandler.GetSKU)

		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthenticateToken)

			// Product operations requiring authentication
			r.Post("/", productHandler.NewProduct)
			r.Put("/{id}", productHandler.UpdateProduct)
			r.Delete("{id}", productHandler.DeleteProduct)

			r.Route("/skus", func(r chi.Router) {
				r.Put("/{id}", skuHandler.UpdateSKU)
				r.Delete("/{id}", skuHandler.DeleteSKU)
			})
			r.Post("/{id}/skus", skuHandler.NewSKU)
		})
	})
	// Store routes
	mux.Route("/stores", func(r chi.Router) {
		r.Get("/{store_id}/products", productHandler.GetStoreProducts)
		r.Get("/{store_id}/products/{slug}", productHandler.GetProductBySlug)
	})
	return mux
}
