from abc import ABC, abstractmethod
from typing import List, Optional, Dict, Any
import asyncio
from datetime import datetime

class BaseAIModel(ABC):
    """Base class for all AI models"""
    
    def __init__(self, config: Dict[str, Any]):
        self.config = config
        self.name = config.get('name', 'unknown')
        self.type = config.get('type', 'general')
        self.languages = config.get('languages', [])
        self.specialties = config.get('specialties', [])
        self.max_tokens = config.get('max_tokens', 2048)
        self.temperature = config.get('temperature', 0.7)
        self.is_loaded = False
        
    @abstractmethod
    async def initialize(self) -> bool:
        """Initialize the model"""
        pass
    
    @abstractmethod
    async def generate(self, prompt: str, **kwargs) -> str:
        """Generate response from model"""
        pass
    
    @abstractmethod
    async def cleanup(self):
        """Cleanup resources"""
        pass
    
    def supports_language(self, language: str) -> bool:
        """Check if model supports language"""
        return language in self.languages
    
    def supports_specialty(self, specialty: str) -> bool:
        """Check if model supports specialty"""
        return specialty in self.specialties or 'general' in self.specialties
    
    def get_config(self) -> Dict[str, Any]:
        """Get model configuration"""
        return {
            'name': self.name,
            'type': self.type,
            'languages': self.languages,
            'specialties': self.specialties,
            'max_tokens': self.max_tokens,
            'temperature': self.temperature,
            'is_loaded': self.is_loaded
        }

class HuggingFaceModel(BaseAIModel):
    """Hugging Face transformer model"""
    
    def __init__(self, config: Dict[str, Any]):
        super().__init__(config)
        self.model = None
        self.tokenizer = None
        self.device = "cpu"  # Will be updated based on availability
        
    async def initialize(self) -> bool:
        """Initialize Hugging Face model"""
        try:
            # These will be imported when needed
            from transformers import AutoTokenizer, AutoModelForCausalLM
            import torch
            
            # Check for GPU availability
            if torch.cuda.is_available():
                self.device = "cuda"
                print(f"Using GPU: {torch.cuda.get_device_name()}")
            
            # Load tokenizer and model
            print(f"Loading model: {self.name}")
            self.tokenizer = AutoTokenizer.from_pretrained(
                self.name,
                trust_remote_code=True
            )
            
            self.model = AutoModelForCausalLM.from_pretrained(
                self.name,
                trust_remote_code=True,
                torch_dtype=torch.float16 if self.device == "cuda" else torch.float32,
                device_map="auto" if self.device == "cuda" else None
            )
            
            # Set padding token
            if self.tokenizer.pad_token is None:
                self.tokenizer.pad_token = self.tokenizer.eos_token
            
            self.is_loaded = True
            print(f"Model {self.name} loaded successfully!")
            return True
            
        except Exception as e:
            print(f"Failed to load model {self.name}: {e}")
            return False
    
    async def generate(self, prompt: str, **kwargs) -> str:
        """Generate response using Hugging Face model"""
        if not self.is_loaded:
            raise RuntimeError(f"Model {self.name} is not loaded")
        
        try:
            # Prepare inputs
            inputs = self.tokenizer(
                prompt,
                return_tensors="pt",
                truncation=True,
                max_length=4096
            )
            
            if self.device == "cuda":
                inputs = inputs.to("cuda")
            
            # Generate response
            with torch.no_grad():
                outputs = self.model.generate(
                    **inputs,
                    max_new_tokens=kwargs.get('max_tokens', self.max_tokens),
                    temperature=kwargs.get('temperature', self.temperature),
                    do_sample=True,
                    pad_token_id=self.tokenizer.eos_token_id,
                    eos_token_id=self.tokenizer.eos_token_id,
                )
            
            # Decode response
            response = self.tokenizer.decode(
                outputs[0],
                skip_special_tokens=True
            )
            
            # Remove the original prompt from response
            if response.startswith(prompt):
                response = response[len(prompt):].strip()
            
            return response
            
        except Exception as e:
            print(f"Generation error in model {self.name}: {e}")
            raise
    
    async def cleanup(self):
        """Cleanup model resources"""
        if self.model:
            del self.model
        if self.tokenizer:
            del self.tokenizer
        
        # Clear GPU cache if applicable
        try:
            import torch
            if torch.cuda.is_available():
                torch.cuda.empty_cache()
        except:
            pass
        
        self.is_loaded = False

class OpenAIModel(BaseAIModel):
    """OpenAI API model"""
    
    def __init__(self, config: Dict[str, Any]):
        super().__init__(config)
        self.api_key = config.get('api_key')
        self.client = None
        
    async def initialize(self) -> bool:
        """Initialize OpenAI client"""
        try:
            import openai
            self.client = openai.AsyncOpenAI(api_key=self.api_key)
            self.is_loaded = True
            print(f"OpenAI model {self.name} initialized successfully!")
            return True
        except Exception as e:
            print(f"Failed to initialize OpenAI model {self.name}: {e}")
            return False
    
    async def generate(self, prompt: str, **kwargs) -> str:
        """Generate response using OpenAI API"""
        if not self.is_loaded:
            raise RuntimeError(f"Model {self.name} is not initialized")
        
        try:
            response = await self.client.chat.completions.create(
                model=self.name,
                messages=[{"role": "user", "content": prompt}],
                max_tokens=kwargs.get('max_tokens', self.max_tokens),
                temperature=kwargs.get('temperature', self.temperature)
            )
            
            return response.choices[0].message.content
            
        except Exception as e:
            print(f"OpenAI API error: {e}")
            raise
    
    async def cleanup(self):
        """Cleanup OpenAI client"""
        self.client = None
        self.is_loaded = False