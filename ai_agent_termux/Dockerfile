# Google Lens AI Agent - Docker Deployment
FROM golang:1.21-alpine AS builder

# Install required build dependencies
RUN apk add --no-cache \
    git \
    make \
    gcc \
    musl-dev \
    libgcc \
    linux-headers \
    pkgconfig \
    imagemagick-dev \
    pkgconfig

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o google-lens-agent ./cmd/google-lens-agent

# Production stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    imagemagick \
    tesseract-ocr \
    tesseract-ocr-data-eng \
    tesseract-ocr-data-fra \
    tesseract-ocr-data-deu \
    tesseract-ocr-data-spa \
    wget \
    curl \
    bash \
    python3 \
    py3-pip

# Install Python dependencies for additional AI capabilities
RUN pip3 install \
    torch \
    transformers \
    pillow \
    opencv-python

# Create app user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Create necessary directories
RUN mkdir -p /app/data /app/logs /app/config /app/temp && \
    chown -R appuser:appgroup /app

WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/google-lens-agent .
RUN chown appuser:appgroup /app/google-lens-agent

# Copy configuration files
COPY --chown=appuser:appgroup config/ ./config/
COPY --chown=appuser:appgroup docker/docker-entrypoint.sh /usr/local/bin/

# Make scripts executable
RUN chmod +x /usr/local/bin/docker-entrypoint.sh

# Switch to app user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# Set environment variables
ENV GOOGLE_APPLICATION_CREDENTIALS=/app/config/gcp-credentials.json
ENV SERPAPI_API_KEY=""
ENV NTFY_TOPIC=""
ENV OLLAMA_BASE_URL=http://ollama:11434
ENV TERMUX_API=false

# Entry point
ENTRYPOINT ["/usr/local/bin/docker-entrypoint.sh"]
CMD ["serve"]