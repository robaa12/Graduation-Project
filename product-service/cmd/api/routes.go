package main

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/robaa12/product-service/cmd/api/handlers"
	"github.com/robaa12/product-service/cmd/api/middleware"
	"github.com/robaa12/product-service/cmd/utils"
	"github.com/robaa12/product-service/cmd/validation"
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
	collectionHandler := handlers.CollectionHandler{DB: app.db.DB, Validator: validation.NewCollectionValidator(app.db.DB)}

	// Routes under /stores/{store_id}
	mux.Route("/stores", func(r chi.Router) {
		r.Route("/{store_id}", func(r chi.Router) {
			// Public Product Routes
			r.Get("/products", productHandler.GetStoreProducts)
			r.Get("/products/slug/{slug}", productHandler.GetProductBySlug)

			r.Group(func(r chi.Router) {
				r.Use(middleware.AuthenticateToken)
				r.Use(middleware.VerifyStoreOwnership)
				r.Post("/products", productHandler.NewProduct)
			})

			// Product Detail Routes
			r.Route("/products/{product_id}", func(r chi.Router) {
				// Public endpoints
				r.Get("/", productHandler.GetProduct)
				r.Get("/details", productHandler.GetProductDetails)

				// Protected endpoints
				r.Group(func(r chi.Router) {
					r.Use(middleware.AuthenticateToken)
					r.Use(middleware.VerifyStoreOwnership)

					r.Put("/", productHandler.UpdateProduct)
					r.Delete("/", productHandler.DeleteProduct)
				})

				// SKU Routes
				r.Route("/skus", func(r chi.Router) {
					// Public endpoints
					r.Get("/{sku_id}", skuHandler.GetSKU)

					// Protected endpoints
					r.Group(func(r chi.Router) {
						r.Use(middleware.AuthenticateToken)
						r.Use(middleware.VerifyStoreOwnership)

						r.Post("/", skuHandler.NewSKU)
						r.Put("/{sku_id}", skuHandler.UpdateSKU)
						r.Delete("/{sku_id}", skuHandler.DeleteSKU)
					})
				})
			})

			// Collection Routes /stores/{store_id}/collections
			r.Route("/collections", func(r chi.Router) {
				// Public endpoints
				r.Get("/", collectionHandler.GetCollections)
				r.Get("/{collection_id}", collectionHandler.GetCollection)

				// Protected endpoints
				r.Group(func(r chi.Router) {
					r.Use(middleware.AuthenticateToken)
					r.Use(middleware.VerifyStoreOwnership)

					r.Post("/", collectionHandler.CreateCollection)
					r.Post("/{collection_id}/products", collectionHandler.AddProductToCollection)
					r.Delete("/{collection_id}/products/{product_id}", collectionHandler.RemoveProductFromCollection)
				})
			})
		})
	})

	// Token Routes
	if os.Getenv("APP_ENV") != "production" {
		mux.Post("/generate-token", func(w http.ResponseWriter, r *http.Request) {
			var req struct {
				UserID  uint   `json:"user_id"`
				StoreID uint   `json:"store_id"`
				Role    string `json:"role"`
			}

			if err := utils.ReadJSON(w, r, &req); err != nil {
				utils.ErrorJSON(w, err)
				return
			}

			token, err := utils.GenerateToken(req.UserID, req.StoreID, req.Role)
			if err != nil {
				utils.ErrorJSON(w, err)
				return
			}

			utils.WriteJSON(w, http.StatusOK, map[string]string{
				"token": token,
			})
		})
	}

	return mux
}
