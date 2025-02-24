package config

import (
	"embed"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"
)

type RouteConfig struct {
	Path        string   `json:"path"`
	Methods     []string `json:"methods"`
	Service     string   `json:"service"`
	Middlewares []string `json:"middlewares"`
}

type Config struct {
	Server    ServerConfig
	Services  map[string]ServiceConfig
	Routes    []RouteConfig
	Auth      AuthConfig
	RateLimit RateLimitConfig
}

type ServerConfig struct {
	Port string
	Host string
}

type ServiceConfig struct {
	URL     string
	Timeout time.Duration
}

type AuthConfig struct {
	JWTSecret string
	TokenExp  time.Duration
}

type RateLimitConfig struct {
	MaxRequests int
	Duration    time.Duration
}

//go:embed routes.json
var routesFile embed.FS

func LoadRoutesConfig() ([]RouteConfig, error) {
	// Load routes from embedded JSON config
	data, err := routesFile.ReadFile("routes.json")
	if err != nil {
		log.Println("Error loading route config:", err)
		return nil, err
	}

	var routes []RouteConfig
	if err := json.Unmarshal(data, &routes); err != nil {
		log.Println("Error parsing route config:", err)
		return nil, err
	}

	return routes, nil
}
func Load() (*Config, error) {
	routes, err := LoadRoutesConfig()
	if err != nil {
		log.Println("Failed Loading Config")
		return nil, err
	}
	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Host: getEnv("SERVER_HOST", "localhost"),
		},
		Services: map[string]ServiceConfig{ // Changed ServicesConfig to map
			"user-service": {
				URL:     getEnv("USER_SERVICE_URL", "http://localhost:3000"),
				Timeout: getDurationEnv("USER_SERVICE_TIMEOUT", 5*time.Second),
			},
			"order-service": {
				URL:     getEnv("ORDER_SERVICE_URL", "http://localhost:8082"),
				Timeout: getDurationEnv("ORDER_SERVICE_TIMEOUT", 5*time.Second),
			},
			"product-service": {
				URL:     getEnv("PRODUCT_SERVICE_URL", "http://localhost:8083"),
				Timeout: getDurationEnv("PRODUCT_SERVICE_TIMEOUT", 5*time.Second),
			},
		},

		Auth: AuthConfig{
			JWTSecret: getEnv("JWT_SECRET", "Messi is better than Ronaldo"),
			TokenExp:  getDurationEnv("JWT_EXPIRATION", 24*time.Hour),
		},
		RateLimit: RateLimitConfig{
			MaxRequests: getEnvInt("RATE_LIMIT_MAX_REQUESTS", 100),
			Duration:    getDurationEnv("RATE_LIMIT_DURATION", 1*time.Minute),
		},
		Routes: routes,
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
