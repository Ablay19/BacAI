# AGENTS.md - Development Guidelines for BACAI

## 🚀 Build, Test & Lint Commands

### Root Level Commands

```bash
# Development - Start all services
npm run dev                    # Start frontend, API, and model-service

# Build all services
npm run build                   # Build frontend, API, and model-service

# Testing
npm test                        # Run Jest tests across workspaces
npm run test:watch             # Run tests in watch mode

# Linting & Formatting
npm run lint                   # ESLint for all TS/JS files
npm run lint:fix               # Auto-fix linting issues
npm run format                 # Format with Prettier

# Setup
npm run setup                  # Install all dependencies (Node + Python)
```

### Individual Service Commands

#### Frontend (React + Vite + TypeScript)

```bash
cd frontend
npm run dev                    # Start development server (localhost:5173)
npm run build                  # Build for production
npm run preview                # Preview production build
npm run lint                   # ESLint frontend code
npm run lint:fix              # Auto-fix frontend issues
```

#### API (Hono + Cloudflare Workers + TypeScript)

```bash
cd api
npm run dev                    # Start Wrangler dev server
npm run build                  # Compile TypeScript
npm run deploy                 # Deploy to Cloudflare Workers
npm run test                   # Run Jest tests
npm run lint                   # ESLint API code
```

#### Model Service (FastAPI + Python)

```bash
cd model-service
uvicorn main:app --reload      # Start development server (localhost:8000)
pytest                         # Run all tests
pytest tests/test_api.py       # Run single test file
pytest -k "test_solve"         # Run specific test by name
pytest -v                      # Verbose test output
pytest --cov                   # Run with coverage
black .                        # Format Python code
ruff check .                   # Lint Python code
ruff check . --fix             # Auto-fix Python linting
```

### Makefile Commands

```bash
make help                      # Show all available commands
make install                   # Install all dependencies
make dev                       # Start development servers
make build                     # Build all services
make test                      # Run tests
make deploy                    # Deploy to production
make clean                     # Clean build artifacts
make docker-build              # Build Docker images
make docker-run                # Run with Docker
make setup                     # First-time setup
```

## 📝 Code Style Guidelines

### TypeScript/JavaScript (Frontend & API)

#### Import Organization

```typescript
// 1. External libraries (react, next, etc.)
import React from "react";
import { useState } from "react";
import axios from "axios";

// 2. Internal utilities (alphabetical)
import { callModelService } from "../utils/model-client";
import { detectLanguage } from "../utils/language-utils";

// 3. Types and interfaces
import type { Subject, Language } from "../types";

// 4. Component imports (if any)
import { ComponentName } from "../components/ComponentName";
```

#### TypeScript Interfaces

```typescript
// Use interface for objects, type for unions/primitives
interface SubjectCardProps {
  subject: Subject;
  onSelect: (subject: Subject) => void;
  index: number;
}

type Language = "ar" | "fr" | "en";
type EducationLevel = "secondary_basic" | "secondary_lycee" | "university";

// Prefer readonly for immutable data
interface ApiResponse<T> {
  readonly data: T;
  readonly success: boolean;
  readonly message?: string;
}
```

#### Component Structure

```typescript
// React functional component with props interface
interface ComponentProps {
  title: string
  onSubmit: (data: FormData) => void
}

export default function Component({ title, onSubmit }: ComponentProps) {
  const [state, setState] = useState<string>('')

  // Event handlers
  const handleSubmit = useCallback((data: FormData) => {
    onSubmit(data)
  }, [onSubmit])

  // Early returns for edge cases
  if (!title) return null

  return (
    <div className="container">
      {/* JSX content */}
    </div>
  )
}
```

#### Error Handling

```typescript
// API endpoints - consistent error responses
try {
  const result = await callModelService(payload);
  return c.json({ success: true, data: result });
} catch (error) {
  console.error("Operation failed:", error);
  return c.json(
    {
      success: false,
      error: "Operation failed",
      code: "OPERATION_ERROR",
      details: error.message,
    },
    500,
  );
}

// Async operations with proper typing
const fetchSubject = async (id: string): Promise<Subject> => {
  try {
    const response = await api.get<Subject>(`/subjects/${id}`);
    return response.data;
  } catch (error) {
    throw new Error(`Failed to fetch subject: ${error.message}`);
  }
};
```

### Python (Model Service)

#### Import Organization

```python
# 1. Standard library
import asyncio
import os
from typing import Dict, List, Optional

# 2. Third-party libraries
import pytest
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from transformers import AutoTokenizer

# 3. Local imports
from models.base_models import BaseModel as CustomBaseModel
from utils.config import get_config
from utils.helpers import process_text
```

#### Type Hints & Data Models

```python
from typing import Dict, List, Optional, Union
from pydantic import BaseModel, Field

class ExerciseRequest(BaseModel):
    """Request model for exercise processing"""
    exercise: str = Field(..., min_length=1, description="Exercise text")
    subject: str = Field(..., regex="^(mathematics|arabic|french)$")
    level: Optional[str] = Field(None, regex="^(secondary_basic|secondary_lycee)$")
    language: str = Field(default="en", regex="^(ar|fr|en)$")

class ModelResponse(BaseModel):
    """Response model from AI models"""
    solution: str
    confidence: float
    processing_time: int
    model_used: str
```

#### Async/Await Patterns

```python
import asyncio
from typing import List

async def process_multiple_exercises(exercises: List[str]) -> List[dict]:
    """Process multiple exercises concurrently"""
    tasks = [process_single_exercise(ex) for ex in exercises]
    results = await asyncio.gather(*tasks, return_exceptions=True)

    # Handle exceptions
    processed_results = []
    for i, result in enumerate(results):
        if isinstance(result, Exception):
            processed_results.append({
                "exercise": exercises[i],
                "error": str(result),
                "success": False
            })
        else:
            processed_results.append(result)

    return processed_results
```

#### Error Handling

```python
from fastapi import HTTPException
import logging

logger = logging.getLogger(__name__)

class ModelServiceError(Exception):
    """Custom exception for model service errors"""
    pass

async def call_model(payload: dict) -> dict:
    """Call AI model with proper error handling"""
    try:
        response = await model_api.generate(payload)
        return response.json()
    except ConnectionError as e:
        logger.error(f"Model API connection failed: {e}")
        raise HTTPException(
            status_code=503,
            detail="Model service unavailable"
        )
    except Exception as e:
        logger.error(f"Unexpected error in model call: {e}")
        raise ModelServiceError(f"Model processing failed: {str(e)}")
```

### Testing Guidelines

#### Test Naming & Structure

```typescript
// API endpoint tests
describe("Solve API Endpoint", () => {
  beforeEach(() => {
    // Setup before each test
  });

  it("should solve mathematics exercise correctly", async () => {
    const requestData = { exercise: "2+2=4", subject: "mathematics" };
    const response = await request(app)
      .post("/api/solve")
      .send(requestData)
      .expect(200);

    expect(response.body.success).toBe(true);
    expect(response.body.data.solution).toBeDefined();
  });

  it("should return 400 for invalid input", async () => {
    await request(app)
      .post("/api/solve")
      .send({ exercise: "" }) // Invalid empty exercise
      .expect(400);
  });
});
```

```python
# Python test structure
class TestModelManager:
    """Test suite for ModelManager functionality"""

    @pytest.fixture
    def mock_config(self):
        """Setup mock configuration"""
        return create_test_config()

    @pytest.fixture
    def model_manager(self, mock_config):
        """Setup model manager instance"""
        return ModelManager(mock_config)

    @pytest.mark.asyncio
    async def test_solve_mathematics_exercise(self, model_manager):
        """Test solving mathematics exercises"""
        request_data = {
            "exercise": "Solve: x² + 5x - 3 = 0",
            "subject": "mathematics",
            "language": "en"
        }

        result = await model_manager.solve_exercise(request_data)

        assert result["success"] is True
        assert "solution" in result["data"]
        assert result["data"]["subject"] == "mathematics"
```

## 🔧 Development Workflow

### Before Committing

1. **Run linting**: `npm run lint` (or `ruff check .` for Python)
2. **Run tests**: `npm test` (or `pytest` for Python)
3. **Check types**: TypeScript automatically checks on build
4. **Format code**: `npm run format` (or `black .` for Python)

### Git Workflow

```bash
# Feature development
git checkout -b feature/new-subject-type
npm run dev  # Test your changes
npm test && npm run lint  # Verify quality
git add .
git commit -m "feat: add physics subject support"
git push origin feature/new-subject-type
```

### Single Test Execution

```bash
# Node.js/TypeScript
npm test -- --testNamePattern="solve endpoint"
npm test api/tests/solve.test.ts

# Python
pytest tests/test_api.py::TestAPIEndpoints::test_solve_endpoint
pytest -k "solve_endpoint" -v
```

## 🎯 Project-Specific Guidelines

### Multilingual Support

- Always test with Arabic (`ar`), French (`fr`), and English (`en`)
- Use proper RTL/LTR styling in frontend
- Validate language codes with regex patterns

### Mauritanian Context

- Include cultural context examples in AI responses
- Support local curriculum levels: `secondary_basic`, `secondary_lycee`, `university`
- Use appropriate educational terminology for each region

### Performance Considerations

- API responses should be under 2 seconds
- Use async/await for all I/O operations
- Implement proper caching for repeated requests
- Monitor model response times and implement timeouts

### Security Guidelines

- Validate all input with Zod schemas (Node.js) or Pydantic (Python)
- Sanitize user inputs before processing
- Never log sensitive information (API keys, personal data)
- Use HTTPS for all API communications
