package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server    ServerConfig
	Services  ServicesConfig
	Auth      AuthConfig
	RateLimit RateLimitConfig
}

type ServerConfig struct {
	Port string
	Host string
}

type ServicesConfig struct {
	UserService    ServiceConfig
	OrderService   ServiceConfig
	ProductService ServiceConfig
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

func Load() (*Config, error) {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
		},
		Services: ServicesConfig{
			UserService: ServiceConfig{
				URL:     getEnv("USER_SERVICE_URL", "http://localhost:3000"),
				Timeout: getDurationEnv("USER_SERVICE_TIMEOUT", 5*time.Second),
			},
			OrderService: ServiceConfig{
				URL:     getEnv("ORDER_SERVICE_URL", "http://localhost:8082"),
				Timeout: getDurationEnv("ORDER_SERVICE_TIMEOUT", 5*time.Second),
			},
			ProductService: ServiceConfig{
				URL:     getEnv("PRODUCT_SERVICE_URL", "http://localhost:8083"),
				Timeout: getDurationEnv("PRODUCT_SERVICE_TIMEOUT", 5*time.Second),
			},
		},
		Auth: AuthConfig{
			JWTSecret: getEnv("JWT_SECRET", "Messi is better than Ronaldo"),
			TokenExp:  getDurationEnv("TOKEN_EXP", 24*time.Hour),
		},
		RateLimit: RateLimitConfig{
			MaxRequests: getEnvInt("RATE_LIMIT_MAX_REQUESTS", 100),
			Duration:    getDurationEnv("RATE_LIMIT_DURATION", 1*time.Minute),
		},
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
