package api

import (
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/robaa12/gatway-service/internal/auth"
	"github.com/robaa12/gatway-service/internal/proxy"
)

func Routes(authService *auth.AuthService, proxyService *proxy.ProxyService) *chi.Mux {
	router := chi.NewRouter()

	// Middleware
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Public Routes
	router.Group(func(r chi.Router) {
		r.Post("/login", authService.Login)
		r.Post("/user/register", proxyService.UserServiceProxy().ServeHTTP)
		r.Get("/store/{store_id}", proxyService.UserServiceProxy().ServeHTTP)
		r.Route("/store/{store_id}/products", func(r chi.Router) {
			r.Get("/", proxyService.ProductServiceProxy().ServeHTTP)
			r.Get("/{product_id}", proxyService.ProductServiceProxy().ServeHTTP)
			r.Get("/products/slug/{slug}", proxyService.ProductServiceProxy().ServeHTTP)
			r.Get("{produt_id}/details", proxyService.ProductServiceProxy().ServeHTTP)
			r.Get("{product_id}/sku/{sku_id}", proxyService.ProductServiceProxy().ServeHTTP)

		})
		r.Route("/store/{store_id}/collection", func(r chi.Router) {
			r.Get("/", proxyService.ProductServiceProxy().ServeHTTP)
			r.Get("/{collection_id}", proxyService.ProductServiceProxy().ServeHTTP)
		})
	})

	// Protected Routes
	router.Group(func(r chi.Router) {
		r.Use(authService.AuthMiddleware)
		r.Use(middleware.ThrottleBacklog(100, 50, 60000)) // Rate limiting

		// User Service Routes
		r.Mount("/store", proxyService.UserServiceProxy())
		r.Mount("/user", proxyService.UserServiceProxy())
		// Product Service Routes
		r.Mount("/stores/{store_id}", proxyService.ProductServiceProxy())

		// Order Service Routes
		r.Mount("/order", proxyService.OrderServiceProxy())

	})
	return router
}
