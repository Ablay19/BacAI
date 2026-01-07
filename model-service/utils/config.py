import os
from typing import Dict, Any, Optional
from dataclasses import dataclass
import yaml
from pathlib import Path

@dataclass
class ModelConfig:
    name: str
    type: str
    languages: list[str]
    specialties: list[str]
    max_tokens: int
    temperature: float

@dataclass
class DatabaseConfig:
    url: str
    pool_size: int = 10
    max_overflow: int = 20

@dataclass
class RedisConfig:
    url: str
    ttl: int = 3600
    max_connections: int = 10

@dataclass
class Config:
    # Model configurations
    models: Dict[str, ModelConfig]
    
    # Database configuration
    database: DatabaseConfig
    
    # Redis configuration
    redis: RedisConfig
    
    # Service configuration
    service_name: str = "bacai-model-service"
    environment: str = "development"
    port: int = 8000
    log_level: str = "INFO"
    
    # AI service configuration
    huggingface_token: Optional[str] = None
    openai_api_key: Optional[str] = None
    anthropic_api_key: Optional[str] = None
    
    # Model serving configuration
    max_concurrent_requests: int = 10
    request_timeout: int = 30
    model_cache_dir: str = "./models/pretrained"
    
    # Curriculum configuration
    curriculum_config_path: str = "./config/curriculum.yaml"

def get_config() -> Config:
    """Load configuration from environment and config files"""
    
    # Load from environment variables
    config_path = os.getenv("CONFIG_PATH", "./config/curriculum.yaml")
    
    # Load curriculum config
    curriculum_data = {}
    if os.path.exists(config_path):
        with open(config_path, 'r', encoding='utf-8') as f:
            curriculum_data = yaml.safe_load(f)
    
    # Model configurations
    models_config = {}
    if 'models' in curriculum_data:
        for model_name, model_data in curriculum_data['models'].items():
            models_config[model_name] = ModelConfig(
                name=model_data['name'],
                type=model_data['type'],
                languages=model_data['languages'],
                specialties=model_data['specialties'],
                max_tokens=model_data.get('max_tokens', 2048),
                temperature=model_data.get('temperature', 0.7)
            )
    
    # Database configuration
    database_url = os.getenv("DATABASE_URL", "sqlite:///./bacai.db")
    if not database_url.startswith("sqlite"):
        # PostgreSQL or other database
        database_config = DatabaseConfig(
            url=database_url,
            pool_size=int(os.getenv("DB_POOL_SIZE", "10")),
            max_overflow=int(os.getenv("DB_MAX_OVERFLOW", "20"))
        )
    else:
        # SQLite
        database_config = DatabaseConfig(url=database_url)
    
    # Redis configuration
    redis_url = os.getenv("REDIS_URL", "redis://localhost:6379/0")
    redis_config = RedisConfig(
        url=redis_url,
        ttl=int(os.getenv("REDIS_TTL", "3600")),
        max_connections=int(os.getenv("REDIS_MAX_CONNECTIONS", "10"))
    )
    
    return Config(
        models=models_config,
        database=database_config,
        redis=redis_config,
        service_name=os.getenv("SERVICE_NAME", "bacai-model-service"),
        environment=os.getenv("ENVIRONMENT", "development"),
        port=int(os.getenv("PORT", "8000")),
        log_level=os.getenv("LOG_LEVEL", "INFO"),
        huggingface_token=os.getenv("HUGGINGFACE_API_KEY"),
        openai_api_key=os.getenv("OPENAI_API_KEY"),
        anthropic_api_key=os.getenv("ANTHROPIC_API_KEY"),
        max_concurrent_requests=int(os.getenv("MAX_CONCURRENT_REQUESTS", "10")),
        request_timeout=int(os.getenv("REQUEST_TIMEOUT", "30")),
        model_cache_dir=os.getenv("MODEL_CACHE_DIR", "./models/pretrained"),
        curriculum_config_path=config_path
    )

def get_env_variable(key: str, default: Optional[str] = None, required: bool = False) -> Optional[str]:
    """Get environment variable with optional default and required validation"""
    value = os.getenv(key, default)
    
    if required and not value:
        raise ValueError(f"Required environment variable {key} is not set")
    
    return value