# Makefile for BACAI System

.PHONY: help install dev build test deploy clean docker-build docker-run

# Default target
help:
	@echo "BACAI - Mauritanian AI Educational System"
	@echo ""
	@echo "Available commands:"
	@echo "  install      Install all dependencies"
	@echo "  dev          Start development servers"
	@echo "  build        Build all services"
	@echo "  test         Run tests"
	@echo "  deploy       Deploy to production"
	@echo "  clean        Clean build artifacts"
	@echo "  docker-build Build Docker images"
	@echo "  docker-run   Run services with Docker"

# Install dependencies
install:
	@echo "📦 Installing dependencies..."
	npm install
	cd frontend && npm install
	cd api && npm install
	cd model-service && pip install -r requirements.txt
	@echo "✅ Dependencies installed"

# Development
dev:
	@echo "🚀 Starting development servers..."
	npm run dev

# Build
build:
	@echo "🏗️ Building all services..."
	npm run build
	@echo "✅ Build completed"

# Test
test:
	@echo "🧪 Running tests..."
	npm run test

# Deploy to free cloud services
deploy:
	@echo "☁️ Deploying to free cloud services..."
	cd frontend && npx vercel --prod
	cd api && npx wrangler deploy
	cd model-service && railway up
	@echo "✅ Deployment completed"

# Clean
clean:
	@echo "🧹 Cleaning build artifacts..."
	rm -rf node_modules
	rm -rf frontend/node_modules
	rm -rf api/node_modules
	rm -rf frontend/dist
	rm -rf api/dist
	find . -name "*.pyc" -delete
	find . -name "__pycache__" -type d -exec rm -rf {} +
	@echo "✅ Clean completed"

# Docker
docker-build:
	@echo "🐳 Building Docker images..."
	docker-compose build
	@echo "✅ Docker images built"

docker-run:
	@echo "🐳 Running services with Docker..."
	docker-compose up -d
	@echo "✅ Services running"

# Setup for first time
setup: install
	@echo "⚙️ Creating .env files..."
	@if [ ! -f .env ]; then cp .env.example .env; echo "✅ Created .env file - please configure your API keys"; fi
	@echo "🎉 Setup complete!"

# Local development helpers
dev-frontend:
	cd frontend && npm run dev

dev-api:
	cd api && npm run dev

dev-model-service:
	cd model-service && uvicorn main:app --reload --host 0.0.0.0 --port 8000

# Database setup (if using local database)
db-setup:
	@echo "🗄️ Setting up database..."
	cd model-service && python -m alembic upgrade head
	@echo "✅ Database setup complete"

# Model download (for local development)
download-models:
	@echo "📥 Downloading models (this may take a while)..."
	cd model-service && python scripts/download_models.py
	@echo "✅ Models downloaded"