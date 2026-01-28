# Google Lens AI Agent with Docker & Ollama Integration

A powerful, production-ready AI agent that integrates Google Cloud Vision API with local Ollama processing for comprehensive image analysis and document processing.

## 🚀 Features

- **Google Cloud Vision API Integration**: Real OCR, object detection, and text extraction
- **Local AI Processing**: Ollama integration with models like LLaVA for vision tasks
- **Advanced Image Preprocessing**: Resize, enhance, crop, and optimize images
- **Smart File Watching**: Monitor directories and auto-process new images
- **Interactive CLI**: Terminal user interface with progress bars
- **Batch Processing**: Process multiple images concurrently
- **Notification System**: Real-time notifications via ntfy.sh and Termux API
- **Background Task Management**: Goroutine pool with proper lifecycle management
- **Docker Deployment**: Full containerized deployment with Docker Compose
- **Monitoring**: Prometheus metrics and Grafana dashboards
- **Multi-provider Support**: OpenAI, Google Cloud, Anthropic (extensible)

## 🐳 Docker & Ollama Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Google Lens   │    │     Ollama     │    │   Redis Cache   │
│     Agent       │◄──►│   LLM Service   │◄──►│     Queue       │
│   (Go Service) │    │  (Local AI)     │    │   & Storage     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
         ┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│    PostgreSQL   │    │   Prometheus    │    │     Grafana     │
│   (Storage)     │    │  (Metrics)      │    │  (Dashboards)   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 🛠️ Quick Start with Docker

### 1. Prerequisites
```bash
# Install Docker and Docker Compose
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh
pip install docker-compose

# Clone the repository
git clone https://github.com/your-org/google-lens-agent.git
cd google-lens-agent
```

### 2. Configure Environment
```bash
# Copy and edit environment file
cp .env.example .env

# Edit with your API keys and settings
nano .env
```

### 3. Quick Deploy
```bash
# Start all services
make quick-start

# Or manually
docker-compose up -d

# Pull Ollama models
make setup-ollama
```

### 4. Access Services
- **Google Lens Agent**: http://localhost:8080
- **Grafana Dashboard**: http://localhost:3000 (admin/admin)
- **Prometheus**: http://localhost:9090
- **MinIO Console**: http://localhost:9001

## 🏗️ Development Setup

### Local Development
```bash
# Setup development environment
make setup

# Install dependencies
make deps

# Run in interactive mode
make quick-dev

# Or with Go directly
go run ./cmd/google-lens-agent serve -mode=interactive -verbose
```

### Build & Test
```bash
# Run tests
make test

# Run tests with coverage
make test-coverage

# Run benchmarks
make benchmark

# Run linter
make lint

# Run security scan
make security-scan
```

## 📱 Usage Examples

### Interactive CLI Mode
```bash
# Start interactive CLI
./google-lens-agent interactive

# Inside CLI:
>>> process image.jpg extract_text
>>> batch image1.jpg image2.jpg identify_object
>>> watch /path/to/images extract_text
>>> history
```

### Command Line Interface
```bash
# Process single image
./google-lens-agent -mode=batch -files=photo.jpg -operation=extract_text

# Batch processing
./google-lens-agent -mode=batch -files="img1.jpg,img2.png,img3.jpg" -operation=document

# Watch directory
./google-lens-agent -mode=serve -watch=/path/to/images -operation=extract_text
```

### Docker Commands
```bash
# Build and run
make docker-build
make docker-run

# Scale workers
docker-compose up --scale worker=3

# View logs
make logs

# Monitor services
make monitor
```

## 🔧 Configuration

### Environment Variables
```bash
# Google Cloud Vision
GOOGLE_APPLICATION_CREDENTIALS=/path/to/credentials.json

# Ollama Configuration
OLLAMA_BASE_URL=http://ollama:11434
OLLAMA_MODEL=llava

# Notifications
NTFY_TOPIC=google-lens-agent
SERPAPI_API_KEY=your_key_here

# Performance
MAX_CONCURRENT_TASKS=5
CACHE_SIZE=1000
```

### Supported Operations
- `extract_text` - OCR text extraction
- `identify_object` - Object detection and classification
- `solve_math` - Mathematical problem solving
- `translate_text` - Text translation
- `barcode` - Barcode and QR code scanning
- `document` - Document analysis and summarization

## 🐳 Docker Services

### Core Services
- **google-lens-agent**: Main application container
- **ollama**: Local LLM service for AI processing
- **redis**: Caching and task queue
- **nginx**: Reverse proxy and load balancer

### Monitoring
- **prometheus**: Metrics collection
- **grafana**: Visualization dashboards

### Storage (Optional)
- **postgres**: Persistent data storage
- **minio**: Object storage for images

## 📊 Monitoring & Observability

### Prometheus Metrics
- Request rate and latency
- Error rates and types
- Resource utilization
- Cache hit/miss ratios
- Task queue depth

### Grafana Dashboards
- System performance
- API monitoring
- Resource usage
- Error tracking

### Health Checks
```bash
# Service health
curl http://localhost:8080/health

# Ollama health
curl http://localhost:11434/api/tags

# Redis health
redis-cli ping
```

## 🔄 CI/CD Pipeline

### Automated Deployments
```bash
# CI pipeline
make ci

# Release deployment
make release

# Environment-specific deployments
make deploy-staging
make deploy-production
```

### Multi-Platform Builds
```bash
# Build for all platforms
make build-all

# Docker multi-platform build
make docker-build-all
```

## 🚀 Scaling & Performance

### Horizontal Scaling
```yaml
# docker-compose.yml
services:
  worker:
    replicas: 5
    deploy:
      resources:
        limits:
          memory: 2G
          cpus: '2'
```

### Performance Tuning
```bash
# Environment variables
MAX_CONCURRENT_TASKS=10
CACHE_SIZE=5000
PROCESSING_TIMEOUT=60m
PREPROCESSING_ENABLED=true
```

## 🔐 Security

### Security Features
- API key management through environment variables
- Secure communication between services
- Rate limiting and request validation
- Security scanning in CI pipeline

### Security Scans
```bash
# Run security analysis
make security-scan

# Vulnerability scanning
docker run --rm -v "$PWD":/app clair-scanner:latest
```

## 🐛 Troubleshooting

### Common Issues

**Service won't start**
```bash
# Check logs
docker-compose logs google-lens-agent

# Check environment
docker-compose exec google-lens-agent env | grep -E "(GOOGLE|OLLAMA|NTFY)"
```

**Ollama not responding**
```bash
# Restart Ollama
docker-compose restart ollama

# Check model
curl http://localhost:11434/api/tags
```

**High memory usage**
```bash
# Limit resources
docker-compose up --scale worker=1

# Monitor resources
docker stats
```

### Debug Mode
```bash
# Enable debug logging
LOG_LEVEL=debug make docker-compose-up

# Interactive debug container
docker run -it --entrypoint /bin/bash google-lens-agent:latest
```

## 📚 API Documentation

### REST Endpoints
```bash
# Health check
GET /health

# Process image
POST /api/v1/process
{
  "image_path": "/path/to/image.jpg",
  "operation": "extract_text",
  "options": {
    "preprocess": true,
    "enhance": true
  }
}

# Batch process
POST /api/v1/batch
{
  "images": ["/path/to/img1.jpg", "/path/to/img2.jpg"],
  "operation": "identify_object"
}
```

### Webhook Events
```json
{
  "event": "image_processed",
  "data": {
    "image_path": "/path/to/image.jpg",
    "result": "Extracted text content...",
    "timestamp": "2024-01-01T12:00:00Z"
  }
}
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests and linting
5. Submit a pull request

### Development Workflow
```bash
# Setup development environment
make setup

# Make changes
git checkout -b feature/new-feature

# Test changes
make test
make lint

# Submit
git push origin feature/new-feature
gh pr create
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Google Cloud Vision API
- Ollama Project
- Docker & Docker Compose
- Prometheus & Grafana
- Go community packages

## 📞 Support

- **Issues**: [GitHub Issues](https://github.com/your-org/google-lens-agent/issues)
- **Discussions**: [GitHub Discussions](https://github.com/your-org/google-lens-agent/discussions)
- **Documentation**: [Wiki](https://github.com/your-org/google-lens-agent/wiki)

---

**Built with Go, Docker, and ❤️ for AI-powered image processing**