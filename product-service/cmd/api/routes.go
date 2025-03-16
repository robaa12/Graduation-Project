package main

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware" // ✅ Corrected import for Chi middleware
	"github.com/go-chi/cors"
	"github.com/robaa12/product-service/cmd/api/handlers" // ✅ Alias for your custom middleware
	"github.com/robaa12/product-service/cmd/database"
	"github.com/robaa12/product-service/cmd/repository"
	"github.com/robaa12/product-service/cmd/service"
	"github.com/robaa12/product-service/cmd/utils"
	"github.com/robaa12/product-service/cmd/validation"
)

func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()

	// Middleware
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Initialize handlers
	OrderHandler := handlers.OrderHandler{DB: app.db.DB}
	productHandler := handlers.ProductHandler{DB: app.db.DB}

	collectionHandler := handlers.CollectionHandler{DB: app.db.DB, Validator: validation.NewCollectionValidator(app.db.DB)}

	mux.Post("/verify-order", OrderHandler.VerifyOrderItems)
	mux.Post("/update-inventory", OrderHandler.UpdateInventory)
	// Routes under /stores/{store_id}
	mux.Route("/stores", func(r chi.Router) {
		r.Route("/{store_id}", func(r chi.Router) {
			// Public Product Routes
			r.Get("/products", productHandler.GetStoreProducts)
			r.Get("/products/slug/{slug}", productHandler.GetProductBySlug)

			r.Group(func(r chi.Router) {
				//r.Use(customMiddleware.AuthenticateToken)
				//r.Use(customMiddleware.VerifyStoreOwnership)
				r.Post("/products", productHandler.NewProduct)
			})

			// Product Detail Routes
			r.Route("/products/{product_id}", func(r chi.Router) {
				// Public endpoints
				r.Get("/", productHandler.GetProduct)
				r.Get("/details", productHandler.GetProductDetails)

				// Protected endpoints
				r.Group(func(r chi.Router) {
					//r.Use(customMiddleware.AuthenticateToken)
					//r.Use(customMiddleware.VerifyStoreOwnership)

					r.Put("/", productHandler.UpdateProduct)
					r.Delete("/", productHandler.DeleteProduct)
				})

				// SKU Routes
				r.Route("/skus", app.sku)
			})

			// Collection Routes /stores/{store_id}/collections
			r.Route("/collections", func(r chi.Router) {
				// Public endpoints
				r.Get("/", collectionHandler.GetCollections)
				r.Get("/{collection_id}", collectionHandler.GetCollection)

				// Protected endpoints
				r.Group(func(r chi.Router) {
					//r.Use(customMiddleware.AuthenticateToken)
					//r.Use(customMiddleware.VerifyStoreOwnership)

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
func setupSKUHandler(db *database.Database) *handlers.SKUHandler {
	skuRepo := repository.NewSkuRepository(db)
	skuService := service.NewSKUService(skuRepo)
	skuHandler := handlers.NewSKUHandler(skuService)

	return skuHandler
}
func (app *Config) sku(r chi.Router) {
	skuHandler := setupSKUHandler(app.db)
	// Public endpoints
	r.Get("/{sku_id}", skuHandler.GetSKU)

	// Protected endpoints
	r.Group(func(r chi.Router) {
		//	r.Use(customMiddleware.AuthenticateToken)
		//	r.Use(customMiddleware.VerifyStoreOwnership)

		r.Post("/", skuHandler.NewSKU)
		r.Put("/{sku_id}", skuHandler.UpdateSKU)
		r.Delete("/{sku_id}", skuHandler.DeleteSKU)
	})
}
