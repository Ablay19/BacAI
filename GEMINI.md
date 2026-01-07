# Gemini Workspace Context: BACAI - Mauritanian AI Educational System

This document provides essential context for the BACAI project, an AI-powered tutoring system for Mauritanian education.

## Project Overview

BACAI is a multi-component system designed to provide AI-driven educational support across various subjects like Mathematics, Arabic, and Sciences. It is tailored to the Mauritanian curriculum.

The architecture consists of three main services:
1.  **Frontend**: A React/Vite-based web application providing the user interface. It runs on `localhost:3000`.
2.  **API**: A Cloudflare Worker using the Hono framework that acts as a router and lightweight backend. It runs on `localhost:8787` and forwards requests to the model service.
3.  **Model Service**: A Python/FastAPI service that performs heavy AI model inference. It runs on `localhost:8000`.

The system also uses Redis for caching and Nginx as a reverse proxy in the Docker setup.

### Technology Stack

*   **Frontend**: React, TypeScript, Vite, Tailwind CSS, Zustand
*   **API**: Cloudflare Workers, Hono, TypeScript
*   **Model Service**: Python, FastAPI, PyTorch, Transformers, LangChain
*   **Orchestration**: Docker Compose, Makefile
*   **Deployment**: Vercel (Frontend), Cloudflare (API), Railway (Model Service)

## Building and Running

This project uses a combination of `make` and `docker-compose` for orchestration.

### Initial Setup

1.  **Install Dependencies**: Run the main installation command. This will install npm packages for the frontend and API, and Python packages for the model service.
    ```bash
    make install
    ```
2.  **Configure Environment**: Copy the example `.env` file. You will need to add your own API keys (e.g., `HUGGINGFACE_API_KEY`).
    ```bash
    cp .env.example .env
    ```

### Local Development (Docker)

The recommended way to run the full system locally is with Docker Compose.

1.  **Build Images**:
    ```bash
    make docker-build
    ```
2.  **Run Services**:
    ```bash
    make docker-run
    ```

The services will be available at:
- Frontend: `http://localhost:3000` (or `http://localhost:80` via Nginx)
- API: `http://localhost:8787`
- Model Service: `http://localhost:8000`

### Local Development (Manual)

You can also run each service manually in separate terminals.

- **Start Frontend**:
  ```bash
  make dev-frontend
  ```
- **Start API Service**:
  ```bash
  make dev-api
  ```
- **Start Model Service**:
  ```bash
  make dev-model-service
  ```

### Other Useful Commands

- **Run Tests**:
  ```bash
  make test
  ```
- **Build for Production**:
  ```bash
  make build
  ```
- **Clean Artifacts**:
  ```bash
  make clean
  ```
- **Deploy**: The `Makefile` contains targets for deploying each service to its respective platform (Vercel, Cloudflare, Railway).
  ```bash
  make deploy
  ```

## Development Conventions

*   **Monorepo Structure**: The project is organized as a monorepo with distinct services in the `frontend/`, `api/`, and `model-service/` directories.
*   **Configuration**: Central configuration is managed in `config/curriculum.yaml`. Service-specific configurations are found in files like `frontend/vite.config.ts` and `api/wrangler.toml`.
*   **API-Driven**: The frontend communicates with the model service via the API service, which acts as a gateway. The API contract is defined by the Pydantic models in `model-service/main.py` and the routes in `api/src/routes/`.
*   **Python Style**: The model service uses `black` for code formatting and `ruff` for linting.
*   **State Management**: The frontend uses Zustand for global state management (e.g., theme, language).
*   **Environment Variables**: Each service is configured via environment variables, as defined in `docker-compose.yml`.
