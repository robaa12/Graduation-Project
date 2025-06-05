# Makefile for Graduation Project
# This Makefile helps with common development tasks

.PHONY: help build up down logs ps clean clean-all clean-volumes

# Colors for prettier output
GREEN := $(shell tput setaf 2)
YELLOW := $(shell tput setaf 3)
BLUE := $(shell tput setaf 4)
RESET := $(shell tput sgr0)

# Default target
.DEFAULT_GOAL := help

# Executables
DOCKER_COMPOSE = docker compose

# Options
DOCKER_COMPOSE_FILE = docker-compose.yaml

help: ## Show this help
	@echo "$(GREEN)Graduation Project Makefile$(RESET)"
	@echo "$(YELLOW)Usage:$(RESET)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(YELLOW)%-15s$(RESET) %s\n", $$1, $$2}'

build: ## Build all Docker images
	@echo "$(BLUE)Building Docker images...$(RESET)"
	@$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) build

build-no-cache: ## Build all Docker images without using cache
	@echo "$(BLUE)Building Docker images without cache...$(RESET)"
	@$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) build --no-cache

up: ## Start all services
	@echo "$(BLUE)Starting services...$(RESET)"
	@$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) up -d
	@echo "$(GREEN)Services started$(RESET)"
	@echo "$(YELLOW)Gateway service is running at: http://localhost:8030$(RESET)"

up-build: build up ## Build and start all services

up-logs: up ## Start all services and follow logs
	@$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) logs -f

down: ## Stop all services
	@echo "$(BLUE)Stopping services...$(RESET)"
	@$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) down
	@echo "$(GREEN)Services stopped$(RESET)"

logs: ## Show logs from all containers
	@$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) logs

logs-follow: ## Follow logs from all containers
	@$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) logs -f

logs-service: ## Show logs from specific service (usage: make logs-service SERVICE=gateway-service)
	@$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) logs $(SERVICE)

ps: ## List running containers
	@$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) ps

restart: ## Restart all services
	@echo "$(BLUE)Restarting services...$(RESET)"
	@$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) restart
	@echo "$(GREEN)Services restarted$(RESET)"

restart-service: ## Restart specific service (usage: make restart-service SERVICE=gateway-service)
	@echo "$(BLUE)Restarting $(SERVICE)...$(RESET)"
	@$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) restart $(SERVICE)
	@echo "$(GREEN)$(SERVICE) restarted$(RESET)"

clean: ## Stop and remove containers
	@echo "$(BLUE)Stopping and removing containers...$(RESET)"
	@$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) down
	@echo "$(GREEN)Containers removed$(RESET)"

clean-all: ## Stop and remove containers, networks, and images
	@echo "$(BLUE)Stopping and removing containers, networks, and images...$(RESET)"
	@$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) down --rmi all
	@echo "$(GREEN)Containers, networks, and images removed$(RESET)"

clean-volumes: ## Stop and remove containers and volumes
	@echo "$(BLUE)Stopping and removing containers and volumes...$(RESET)"
	@$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) down -v
	@echo "$(GREEN)Containers and volumes removed$(RESET)"

# Service-specific commands
gateway: ## Build and start gateway service
	@echo "$(BLUE)Building and starting gateway service...$(RESET)"
	@$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) up -d --build gateway-service
	@echo "$(GREEN)Gateway service started at: http://localhost:8030$(RESET)"

user: ## Build and start user service
	@echo "$(BLUE)Building and starting user service...$(RESET)"
	@$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) up -d --build user-service
	@echo "$(GREEN)User service started at: http://localhost:3000$(RESET)"

product: ## Build and start product service
	@echo "$(BLUE)Building and starting product service...$(RESET)"
	@$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) up -d --build product-service
	@echo "$(GREEN)Product service started at: http://localhost:8083$(RESET)"

order: ## Build and start order service
	@echo "$(BLUE)Building and starting order service...$(RESET)"
	@$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) up -d --build order-service
	@echo "$(GREEN)Order service started at: http://localhost:8082$(RESET)"

db: ## Start just the database
	@echo "$(BLUE)Starting PostgreSQL database...$(RESET)"
	@$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) up -d postgres
	@echo "$(GREEN)Database started at: postgresql://postgres:postgres@localhost:5432$(RESET)"

redis: ## Start just Redis
	@echo "$(BLUE)Starting Redis...$(RESET)"
	@$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) up -d redis
	@echo "$(GREEN)Redis started at: localhost:6379$(RESET)"