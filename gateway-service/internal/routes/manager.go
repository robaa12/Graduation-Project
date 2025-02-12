package routes

import (
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
	Auth   *auth.AuthService
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

	// Convert middlewares slice to map for O(1) lookup
	middlewareMap := make(map[string]bool)
	for _, m := range route.Middlewares {
		middlewareMap[m] = true
	}

	// Create base handler
	handler := proxy.NewProxyService(&service)

	// Apply middlewares in fixed order
	if middlewareMap["auth"] {
		handler = rm.Auth.AuthMiddleware(handler)
	}
	if middlewareMap["store-ownership"] {
		handler = rm.Auth.StoreOwnershipMiddleware(handler)
	}

	return handler
}

func (rm *RouteManager) coreRoutes() {
	rm.Router.Post("/login", rm.Auth.Login)
	rm.Router.Post("/register", rm.Auth.Register)
}
