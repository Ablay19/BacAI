import os
from fastapi import FastAPI, HTTPException, BackgroundTasks
from fastapi.middleware.cors import CORSMiddleware
from fastapi.middleware.gzip import GZipMiddleware
from pydantic import BaseModel, Field
from typing import List, Optional, Dict, Any
import uvicorn
import logging
from datetime import datetime

# Import our modules
from models.model_manager import ModelManager
from models.language_processor import LanguageProcessor
from utils.config import get_config
from utils.logger import setup_logger

# Configure logging
logger = setup_logger(__name__)

# Initialize FastAPI app
app = FastAPI(
    title="BACAI Model Service",
    description="AI model service for Mauritanian educational system",
    version="1.0.0",
    docs_url="/docs",
    redoc_url="/redoc"
)

# Add middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],  # Configure properly for production
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)
app.add_middleware(GZipMiddleware, minimum_size=1000)

# Global variables
model_manager: ModelManager = None
language_processor: LanguageProcessor = None
config = get_config()

# Pydantic models
class SolveRequest(BaseModel):
    exercise: str = Field(..., description="The exercise to solve")
    subject: str = Field(..., description="Subject area")
    level: str = Field(default="secondary_lycee", description="Education level")
    language: str = Field(default="en", description="Language code (ar, fr, en)")
    mode: str = Field(default="step-by-step", description="Solution mode")
    cultural_context: bool = Field(default=True, description="Include cultural context")
    request_type: str = Field(default="solve", description="Request type")

class ExplainRequest(BaseModel):
    concept: str = Field(..., description="Concept to explain")
    subject: str = Field(..., description="Subject area")
    level: str = Field(default="secondary_lycee", description="Education level")
    language: str = Field(default="en", description="Language code")
    detail_level: str = Field(default="detailed", description="Explanation detail")
    examples: bool = Field(default=True, description="Include examples")
    cultural_context: bool = Field(default=True, description="Include cultural context")
    request_type: str = Field(default="explain", description="Request type")

class ConverseRequest(BaseModel):
    message: str = Field(..., description="User message")
    conversation_id: Optional[str] = Field(None, description="Conversation ID")
    subject: str = Field(default="general", description="Subject area")
    language: str = Field(default="en", description="Language code")
    tutor_mode: str = Field(default="friendly", description="Tutoring style")
    step_by_step: bool = Field(default=True, description="Step-by-step responses")
    request_type: str = Field(default="converse", description="Request type")

class SolutionStep(BaseModel):
    explanation: str
    output: str
    step_number: int

class Solution(BaseModel):
    steps: List[SolutionStep]
    final_answer: str
    confidence: float
    language_detected: str

class ModelResponse(BaseModel):
    success: bool
    data: Optional[Dict[str, Any]] = None
    error: Optional[str] = None
    metadata: Optional[Dict[str, Any]] = None

# Startup event
@app.on_event("startup")
async def startup_event():
    global model_manager, language_processor
    
    logger.info("Starting BACAI Model Service...")
    
    try:
        # Initialize model manager
        model_manager = ModelManager(config)
        await model_manager.initialize()
        
        # Initialize language processor
        language_processor = LanguageProcessor(config)
        
        logger.info("Model service started successfully!")
        
    except Exception as e:
        logger.error(f"Failed to start model service: {e}")
        raise

# Shutdown event
@app.on_event("shutdown")
async def shutdown_event():
    global model_manager
    
    logger.info("Shutting down model service...")
    
    if model_manager:
        await model_manager.cleanup()
    
    logger.info("Model service shut down.")

# Health check endpoint
@app.get("/health")
async def health_check():
    return {
        "status": "healthy",
        "timestamp": datetime.now().isoformat(),
        "service": "bacai-model-service",
        "version": "1.0.0",
        "models_loaded": model_manager.get_loaded_models() if model_manager else []
    }

@app.get("/")
async def root():
    return {
        "name": "BACAI Model Service",
        "description": "AI model service for Mauritanian educational system",
        "version": "1.0.0",
        "endpoints": {
            "solve": "/api/process",
            "health": "/health",
            "models": "/models",
            "docs": "/docs"
        }
    }

# Main processing endpoint
@app.post("/api/process", response_model=ModelResponse)
async def process_request(request: Dict[str, Any]):
    try:
        start_time = datetime.now()
        
        # Determine request type and route to appropriate handler
        request_type = request.get("request_type", "solve")
        
        if request_type == "solve":
            response = await handle_solve_request(SolveRequest(**request))
        elif request_type == "explain":
            response = await handle_explain_request(ExplainRequest(**request))
        elif request_type == "converse":
            response = await handle_converse_request(ConverseRequest(**request))
        else:
            raise ValueError(f"Unknown request type: {request_type}")
        
        processing_time = (datetime.now() - start_time).total_seconds()
        
        return ModelResponse(
            success=True,
            data=response,
            metadata={
                "processing_time": processing_time,
                "model_used": "qwen2-8b-instruct",  # This should come from model_manager
                "timestamp": datetime.now().isoformat()
            }
        )
        
    except Exception as e:
        logger.error(f"Processing error: {e}")
        return ModelResponse(
            success=False,
            error=str(e),
            metadata={
                "timestamp": datetime.now().isoformat()
            }
        )

async def handle_solve_request(request: SolveRequest) -> Dict[str, Any]:
    """Handle problem solving requests"""
    
    # Detect language if not provided
    detected_language = language_processor.detect_language(request.exercise)
    language = request.language or detected_language
    
    # Generate solution
    solution = await model_manager.solve_exercise(
        exercise=request.exercise,
        subject=request.subject,
        level=request.level,
        language=language,
        mode=request.mode,
        cultural_context=request.cultural_context
    )
    
    return {
        "solution": solution,
        "subject": request.subject,
        "level": request.level,
        "language": language,
        "mode": request.mode
    }

async def handle_explain_request(request: ExplainRequest) -> Dict[str, Any]:
    """Handle explanation requests"""
    
    # Detect language if not provided
    detected_language = language_processor.detect_language(request.concept)
    language = request.language or detected_language
    
    # Generate explanation
    explanation = await model_manager.explain_concept(
        concept=request.concept,
        subject=request.subject,
        level=request.level,
        language=language,
        detail_level=request.detail_level,
        examples=request.examples,
        cultural_context=request.cultural_context
    )
    
    return {
        "explanation": explanation,
        "concept": request.concept,
        "subject": request.subject,
        "level": request.level,
        "language": language,
        "detail_level": request.detail_level
    }

async def handle_converse_request(request: ConverseRequest) -> Dict[str, Any]:
    """Handle conversation requests"""
    
    # Generate response
    response = await model_manager.converse(
        message=request.message,
        conversation_id=request.conversation_id,
        subject=request.subject,
        language=request.language,
        tutor_mode=request.tutor_mode,
        step_by_step=request.step_by_step
    )
    
    return {
        "response": response,
        "conversation_id": request.conversation_id,
        "subject": request.subject,
        "language": request.language
    }

# Models info endpoint
@app.get("/models")
async def get_models_info():
    """Get information about loaded models"""
    if not model_manager:
        raise HTTPException(status_code=503, detail="Model manager not initialized")
    
    return {
        "loaded_models": model_manager.get_loaded_models(),
        "available_models": model_manager.get_available_models(),
        "model_stats": model_manager.get_model_stats()
    }

if __name__ == "__main__":
    uvicorn.run(
        "main:app",
        host="0.0.0.0",
        port=int(os.getenv("PORT", 8000)),
        reload=True if os.getenv("ENVIRONMENT") == "development" else False,
        log_level="info"
    )