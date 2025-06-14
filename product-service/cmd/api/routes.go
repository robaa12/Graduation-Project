package main

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"                 // ✅ Corrected import for Chi middleware
	"github.com/robaa12/product-service/cmd/api/handlers" // ✅ Alias for your custom middleware
	"github.com/robaa12/product-service/cmd/database"
	"github.com/robaa12/product-service/cmd/repository"
	"github.com/robaa12/product-service/cmd/service"
	"github.com/robaa12/product-service/cmd/utils"
)

func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()

	// Middleware
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)

	// Order Handler
	OrderHandler := handlers.OrderHandler{DB: app.db.DB}

	// Product Handler
	productRepository := repository.NewProductRepository(*app.db)
	reviewRepository := repository.NewReviewRepository(*app.db)
	storeRepository := repository.NewStoreRepository(*app.db)
	reviewService := service.NewReviewService(reviewRepository)
	storeService := service.NewStoreService(storeRepository)

	// Dependancy Injection To access review service in product service
	productService := service.NewProductService(productRepository, reviewService)

	productHandler := handlers.ProductHandler{
		ProductService: *productService,
	}
	storeHandler := handlers.NewStoreHandler(storeService)
	reviewHandler := handlers.NewReviewHandler(reviewService)
	skuHandler := setupSKUHandler(app.db)

	mux.Post("/verify-order", OrderHandler.VerifyOrderItems)
	mux.Post("/update-inventory", OrderHandler.UpdateInventory)

	// Routes under /stores/{store_id}
	mux.Route("/stores", func(r chi.Router) {
		r.Post("/", storeHandler.CreateStore)
		r.Get("/slug/{store_slug}/products", productHandler.GetProductsByStoreSlug)

		r.Route("/{store_id}", func(r chi.Router) {
			// Public Product Routes
			r.Delete("/", storeHandler.DeleteStore)
			r.Get("/products", productHandler.GetStoreProducts)
			r.Get("/products/slug/{slug}", productHandler.GetProductBySlug)
			r.Post("/skus/info", skuHandler.GetSKUs)

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
				r.Get("/reviews", reviewHandler.GetProductReviews)
				r.Get("/reviews/{review_id}", reviewHandler.GetReview)
				r.Get("/reviews/statistics", reviewHandler.GetReviewStatistics)
				r.Post("/reviews", reviewHandler.CreateReview)
			})
			// Collection Routes /stores/{store_id}/collections
			r.Route("/collections", app.collection)
			// Categories Routes /stores/{store_id}/categories
			r.Route("/categories", app.category)
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
				_ = utils.ErrorJSON(w, err)
				return
			}

			token, err := utils.GenerateToken(req.UserID, req.StoreID, req.Role)
			if err != nil {
				_ = utils.ErrorJSON(w, err)
				return
			}

			_ = utils.WriteJSON(w, http.StatusOK, map[string]string{
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
func setupCollectionHandler(db *database.Database) *handlers.CollectionHandler {
	collectionRepo := repository.NewCollectionRepository(db)
	collectionService := service.NewCollectionService(collectionRepo)
	collectionHandler := handlers.NewCollectionHandler(collectionService)

	return collectionHandler
}
func (app *Config) collection(r chi.Router) {
	collectionHandler := setupCollectionHandler(app.db)
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

}
func setupCategoryHandler(db *database.Database) *handlers.CategoryHandler {
	categoryRepo := repository.NewCategoryRepository(db)
	categoryService := service.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	return categoryHandler
}
func (app *Config) category(r chi.Router) {
	categoryHandler := setupCategoryHandler(app.db)
	// Public endpoints
	r.Get("/", categoryHandler.GetCategories)
	r.Post("/", categoryHandler.CreateCategory)
	r.Get("/slug/{category_slug}", categoryHandler.GetCategoryBySlug)
	r.Route("/{category_id}", func(r chi.Router) {
		r.Get("/", categoryHandler.GetCategoryByID)
		r.Post("/", categoryHandler.UpdateCategory)
		r.Delete("/", categoryHandler.DeleteCategory)
	})

}
