package routes

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/robaa12/gatway-service/internal/config"
	"github.com/robaa12/gatway-service/internal/middleware/auth"
	"github.com/robaa12/gatway-service/internal/proxy"
)

type RouteManager struct {
	Router *chi.Mux
	Cfg    *config.Config
	Auth   *auth.Service
}

func NewRouter(cfg *config.Config) *RouteManager {
	rm := RouteManager{
		Router: chi.NewRouter(),
		Cfg:    cfg,
		Auth:   auth.NewAuthService(cfg),
	}
	rm.setupRouter()
	rm.coreRoutes()
	rm.registerRoutes()
	return &rm
}
func (rm *RouteManager) setupRouter() {
	// Middleware
	rm.Router.Use(middleware.Logger)
	rm.Router.Use(middleware.Recoverer)
	rm.Router.Use(middleware.RequestID)
	rm.Router.Use(middleware.RealIP)
	rm.Router.Use(middleware.ThrottleBacklog(100, 50, 60000)) // Rate limiting
	rm.Router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

}

func (rm *RouteManager) registerRoutes() {
	// API Routes
	rm.Router.Route("/", func(r chi.Router) {

		for _, route := range rm.Cfg.Routes {
			handler := rm.createRouteHandler(route)
			log.Println(route)
			// Ensure methods are provided
			if len(route.Methods) == 0 {
				log.Printf("Warning: No methods defined for route %s", route.Path)

			}
			// add methods to router
			for _, method := range route.Methods {
				r.Method(method, route.Path, handler)
			}
		}
	})

}

func (rm *RouteManager) createRouteHandler(route config.RouteConfig) http.Handler {
	service, exists := rm.Cfg.Services[route.Service]
	if !exists {
		log.Printf("Warning: Service not found for route %s", route.Path)
		return http.NotFoundHandler()
	}

	// Create base handler
	handler := proxy.NewProxyService(&service)

	middlewareMap := map[string]func(http.Handler) http.Handler{
		"auth":            rm.Auth.AuthMiddleware,
		"store-ownership": rm.Auth.StoreOwnershipMiddleware,
	}

	middlewareOrder := []string{
		"auth",
		"store-ownership",
	}

	for i := len(middlewareOrder) - 1; i >= 0; i-- {
		mwName := middlewareOrder[i]
		if contains(route.Middlewares, mwName) {
			if mw, exists := middlewareMap[mwName]; exists {
				handler = mw(handler)
			}
		}
	}

	return handler
}

func (rm *RouteManager) coreRoutes() {
	rm.Router.Post("/login", rm.Auth.Login)
	rm.Router.Post("/register", rm.Auth.Register)
	rm.Router.Get("/", rm.sayHello())
	rm.Router.Post("/refresh", rm.Auth.RefreshToken)
}

func (rm *RouteManager) sayHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	}
}

func contains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}
