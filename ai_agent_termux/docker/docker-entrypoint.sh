#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging function
log() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

debug() {
    if [[ "${LOG_LEVEL}" == "debug" ]]; then
        echo -e "${BLUE}[DEBUG]${NC} $1"
    fi
}

# Function to wait for a service
wait_for_service() {
    local host=$1
    local port=$2
    local service_name=$3
    local timeout=${4:-60}
    
    log "Waiting for ${service_name} (${host}:${port})..."
    
    local count=0
    while ! nc -z "$host" "$port"; do
        if [ $count -ge $timeout ]; then
            error "Timeout waiting for ${service_name}"
            exit 1
        fi
        count=$((count + 1))
        sleep 1
    done
    
    log "${service_name} is ready!"
}

# Function to check required environment variables
check_env_vars() {
    local required_vars=("$@")
    local missing_vars=()
    
    for var in "${required_vars[@]}"; do
        if [[ -z "${!var}" ]]; then
            missing_vars+=("$var")
        fi
    done
    
    if [ ${#missing_vars[@]} -gt 0 ]; then
        error "Missing required environment variables: ${missing_vars[*]}"
        exit 1
    fi
}

# Function to setup directories
setup_directories() {
    log "Setting up directories..."
    
    mkdir -p /app/data/{uploads,processed,temp,cache}
    mkdir -p /app/logs
    mkdir -p /app/config
    
    # Set permissions
    chown -R appuser:appgroup /app
    
    log "Directories created and permissions set"
}

# Function to initialize database if needed
init_database() {
    if [[ "${POSTGRES_HOST}" ]]; then
        log "Initializing database connection..."
        wait_for_service "${POSTGRES_HOST}" "${POSTGRES_PORT:-5432}" "PostgreSQL"
    fi
}

# Function to initialize Ollama models
init_ollama() {
    if [[ "${OLLAMA_BASE_URL}" ]]; then
        log "Initializing Ollama models..."
        
        # Extract host from base URL
        local ollama_host=$(echo "${OLLAMA_BASE_URL}" | sed 's|http://||' | sed 's|https://||' | cut -d: -f1)
        local ollama_port=$(echo "${OLLAMA_BASE_URL}" | sed 's|http://||' | sed 's|https://||' | cut -d: -f2)
        
        wait_for_service "${ollama_host}" "${ollama_port:-11434}" "Ollama"
        
        # Pull default model if not exists
        log "Pulling model: ${OLLAMA_MODEL:-llava}"
        curl -s -X POST "${OLLAMA_BASE_URL}/api/pull" \
             -H "Content-Type: application/json" \
             -d "{\"name\": \"${OLLAMA_MODEL:-llava}\"}" || \
             warn "Failed to pull model ${OLLAMA_MODEL:-llava}"
        
        log "Ollama initialization completed"
    fi
}

# Function to check Google credentials
check_google_credentials() {
    if [[ "${GOOGLE_APPLICATION_CREDENTIALS}" ]]; then
        if [[ -f "${GOOGLE_APPLICATION_CREDENTIALS}" ]]; then
            log "Google Cloud credentials found at ${GOOGLE_APPLICATION_CREDENTIALS}"
        else
            warn "Google Cloud credentials file not found: ${GOOGLE_APPLICATION_CREDENTIALS}"
        fi
    else
        warn "GOOGLE_APPLICATION_CREDENTIALS not set"
    fi
}

# Function to run health checks
health_check() {
    log "Running health checks..."
    
    # Check if the service is responding
    if curl -f http://localhost:8080/health > /dev/null 2>&1; then
        log "Service health check passed"
    else
        error "Service health check failed"
        exit 1
    fi
}

# Function to handle graceful shutdown
graceful_shutdown() {
    log "Received shutdown signal, gracefully stopping..."
    
    # Send SIGTERM to the main process
    if [[ -n "$MAIN_PID" ]]; then
        kill -TERM "$MAIN_PID"
        wait "$MAIN_PID"
    fi
    
    log "Graceful shutdown completed"
    exit 0
}

# Main setup function
main_setup() {
    log "Starting Google Lens Agent setup..."
    log "Version: $(cat /app/VERSION 2>/dev/null || echo 'unknown')"
    log "Build: $(git rev-parse --short HEAD 2>/dev/null || echo 'unknown')"
    
    # Setup signal handlers
    trap graceful_shutdown SIGTERM SIGINT
    
    # Check environment
    debug "Environment variables:"
    debug "  GOOGLE_APPLICATION_CREDENTIALS=${GOOGLE_APPLICATION_CREDENTIALS:+SET}"
    debug "  SERPAPI_API_KEY=${SERPAPI_API_KEY:+SET}"
    debug "  NTFY_TOPIC=${NTFY_TOPIC:+SET}"
    debug "  OLLAMA_BASE_URL=${OLLAMA_BASE_URL}"
    debug "  OLLAMA_MODEL=${OLLAMA_MODEL}"
    debug "  LOG_LEVEL=${LOG_LEVEL}"
    
    # Run setup steps
    setup_directories
    check_google_credentials
    init_database
    init_ollama
}

# Main execution
case "${1:-serve}" in
    "serve")
        main_setup
        log "Starting Google Lens Agent in server mode..."
        exec /app/google-lens-agent serve &
        MAIN_PID=$!
        wait $MAIN_PID
        ;;
    "worker")
        main_setup
        log "Starting Google Lens Agent in worker mode..."
        exec /app/google-lens-agent worker &
        MAIN_PID=$!
        wait $MAIN_PID
        ;;
    "batch")
        main_setup
        log "Starting Google Lens Agent in batch mode..."
        exec /app/google-lens-agent batch "${@:2}"
        ;;
    "interactive")
        main_setup
        log "Starting Google Lens Agent in interactive mode..."
        exec /app/google-lens-agent interactive
        ;;
    "health")
        health_check
        ;;
    "shell")
        log "Starting shell..."
        exec /bin/bash
        ;;
    *)
        log "Usage: docker-entrypoint.sh [serve|worker|batch|interactive|health|shell]"
        log "  serve        - Start server mode (default)"
        log "  worker       - Start worker mode"
        log "  batch        - Run batch processing"
        log "  interactive  - Start interactive CLI"
        log "  health       - Run health check"
        log "  shell        - Start shell"
        exit 1
        ;;
esac