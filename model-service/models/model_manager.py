from typing import Dict, List, Any, Optional
import asyncio
from datetime import datetime
import json

from .base_models import BaseAIModel, HuggingFaceModel, OpenAIModel
from ..utils.config import Config, get_config
from ..utils.logger import get_logger

logger = get_logger(__name__)

class ModelManager:
    """Manages multiple AI models and routes requests appropriately"""
    
    def __init__(self, config: Config):
        self.config = config
        self.models: Dict[str, BaseAIModel] = {}
        self.primary_model = None
        self.fallback_model = None
        self.specialized_models = {}
        
    async def initialize(self):
        """Initialize all configured models"""
        logger.info("Initializing model manager...")
        
        try:
            # Initialize primary model
            if 'primary' in self.config.models:
                primary_config = self.config.models['primary']
                self.primary_model = await self._create_model(primary_config)
                if self.primary_model:
                    self.models['primary'] = self.primary_model
                    logger.info(f"Primary model {primary_config.name} initialized")
            
            # Initialize fallback model
            if 'fallback' in self.config.models:
                fallback_config = self.config.models['fallback']
                self.fallback_model = await self._create_model(fallback_config)
                if self.fallback_model:
                    self.models['fallback'] = self.fallback_model
                    logger.info(f"Fallback model {fallback_config.name} initialized")
            
            # Initialize specialized models
            for model_name, model_config in self.config.models.items():
                if model_name not in ['primary', 'fallback']:
                    model = await self._create_model(model_config)
                    if model:
                        self.models[model_name] = model
                        self.specialized_models[model_name] = model
                        logger.info(f"Specialized model {model_config.name} initialized")
            
            if not self.models:
                raise RuntimeError("No models could be initialized")
            
            logger.info(f"Model manager initialized with {len(self.models)} models")
            
        except Exception as e:
            logger.error(f"Failed to initialize model manager: {e}")
            raise
    
    async def _create_model(self, model_config) -> Optional[BaseAIModel]:
        """Create a model instance based on configuration"""
        try:
            model_type = model_config.type.lower()
            
            if model_config.name.startswith('gpt-') or model_config.name.startswith('text-'):
                # OpenAI model
                config_dict = {
                    'name': model_config.name,
                    'type': model_config.type,
                    'api_key': self.config.openai_api_key,
                    **model_config.__dict__
                }
                model = OpenAIModel(config_dict)
            else:
                # Hugging Face model
                config_dict = {
                    'name': model_config.name,
                    'type': model_config.type,
                    **model_config.__dict__
                }
                model = HuggingFaceModel(config_dict)
            
            success = await model.initialize()
            return model if success else None
            
        except Exception as e:
            logger.error(f"Failed to create model {model_config.name}: {e}")
            return None
    
    def _select_model(self, request_data: Dict[str, Any]) -> BaseAIModel:
        """Select the best model for a given request"""
        subject = request_data.get('subject', 'general')
        language = request_data.get('language', 'en')
        
        # Try to find specialized model first
        for model_name, model in self.specialized_models.items():
            if (model.supports_language(language) and 
                model.supports_specialty(subject)):
                logger.debug(f"Selected specialized model {model_name} for {subject}/{language}")
                return model
        
        # Fall back to primary model
        if self.primary_model and self.primary_model.supports_language(language):
            logger.debug(f"Selected primary model for {subject}/{language}")
            return self.primary_model
        
        # Use fallback model
        if self.fallback_model:
            logger.debug(f"Selected fallback model for {subject}/{language}")
            return self.fallback_model
        
        # Return any available model
        if self.models:
            return list(self.models.values())[0]
        
        raise RuntimeError("No models available")
    
    async def solve_exercise(
        self,
        exercise: str,
        subject: str,
        level: str,
        language: str,
        mode: str,
        cultural_context: bool
    ) -> Dict[str, Any]:
        """Solve an educational exercise"""
        
        # Select appropriate model
        model = self._select_model({
            'subject': subject,
            'language': language
        })
        
        # Create prompt
        prompt = self._create_solve_prompt(
            exercise, subject, level, language, mode, cultural_context
        )
        
        # Generate solution
        response = await model.generate(
            prompt,
            max_tokens=model.max_tokens,
            temperature=model.temperature
        )
        
        # Parse response into structured format
        solution = self._parse_solution_response(response, language)
        
        return solution
    
    async def explain_concept(
        self,
        concept: str,
        subject: str,
        level: str,
        language: str,
        detail_level: str,
        examples: bool,
        cultural_context: bool
    ) -> str:
        """Explain a concept"""
        
        # Select appropriate model
        model = self._select_model({
            'subject': subject,
            'language': language
        })
        
        # Create prompt
        prompt = self._create_explain_prompt(
            concept, subject, level, language, detail_level, examples, cultural_context
        )
        
        # Generate explanation
        response = await model.generate(
            prompt,
            max_tokens=model.max_tokens,
            temperature=model.temperature
        )
        
        return response
    
    async def converse(
        self,
        message: str,
        conversation_id: Optional[str],
        subject: str,
        language: str,
        tutor_mode: str,
        step_by_step: bool
    ) -> str:
        """Handle conversational interaction"""
        
        # Select appropriate model
        model = self._select_model({
            'subject': subject,
            'language': language
        })
        
        # Create prompt
        prompt = self._create_converse_prompt(
            message, subject, language, tutor_mode, step_by_step
        )
        
        # Generate response
        response = await model.generate(
            prompt,
            max_tokens=model.max_tokens,
            temperature=model.temperature
        )
        
        return response
    
    def _create_solve_prompt(
        self,
        exercise: str,
        subject: str,
        level: str,
        language: str,
        mode: str,
        cultural_context: bool
    ) -> str:
        """Create a prompt for solving exercises"""
        
        # Language-specific templates
        templates = {
            'ar': f'''
أنت معلم ذكاء اصطناعي متخصص في مواضيع {subject} للمستوى {level}.
يرجى حل المسألة التالية خطوة بخطوة:

المسألة: {exercise}

قدم حلاً مفصلاً مع شرح لكل خطوة.
{'أدرج سياقاً ثقافياً مناسباً للموريتانيا إن أمكن.' if cultural_context else ''}
        ''',
            'fr': f'''
Vous êtes un professeur d'IA spécialisé en {subject} pour le niveau {level}.
Veuillez résoudre l'exercice suivant étape par étape :

Exercice : {exercise}

Fournissez une solution détaillée avec des explications pour chaque étape.
{'Incluez un contexte culturel mauritanien approprié si possible.' if cultural_context else ''}
        ''',
            'en': f'''
You are an AI teacher specializing in {subject} for {level} level.
Please solve the following exercise step by step:

Exercise: {exercise}

Provide a detailed solution with explanations for each step.
{'Include appropriate Mauritanian cultural context if possible.' if cultural_context else ''}
        '''
        }
        
        return templates.get(language, templates['en'])
    
    def _create_explain_prompt(
        self,
        concept: str,
        subject: str,
        level: str,
        language: str,
        detail_level: str,
        examples: bool,
        cultural_context: bool
    ) -> str:
        """Create a prompt for explaining concepts"""
        
        # Language-specific templates
        templates = {
            'ar': f'''
اشرح مفهوم "{concept}" في مادة {subject} لمستوى {level}.

مستوى التفصيل: {detail_level}
{'أضف أمثلة توضيحية.' if examples else ''}
{'أدرج سياقاً ثقافياً موريتانياً.' if cultural_context else ''}
        ''',
            'fr': f'''
Expliquez le concept "{concept}" en {subject} pour le niveau {level}.

Niveau de détail : {detail_level}
{'Ajoutez des exemples illustratifs.' if examples else ''}
{'Incluez un contexte culturel mauritanien.' if cultural_context else ''}
        ''',
            'en': f'''
Explain the concept "{concept}" in {subject} for {level} level.

Detail level: {detail_level}
{'Include illustrative examples.' if examples else ''}
{'Include Mauritanian cultural context.' if cultural_context else ''}
        '''
        }
        
        return templates.get(language, templates['en'])
    
    def _create_converse_prompt(
        self,
        message: str,
        subject: str,
        language: str,
        tutor_mode: str,
        step_by_step: bool
    ) -> str:
        """Create a prompt for conversation"""
        
        # Mode descriptions
        mode_descriptions = {
            'friendly': 'friendly and encouraging',
            'formal': 'formal and professional',
            'encouraging': 'motivational and supportive'
        }
        
        templates = {
            'ar': f'''
أنت مدرس ذكاء اصطناعي {mode_descriptions.get(tutor_mode, "ودود")} لمادة {subject}.
رسالة الطالب: {message}

{'قدم إجابات خطوة بخطوة.' if step_by_step else ''}
        ''',
            'fr': f'''
Vous êtes un tuteur IA {mode_descriptions.get(tutor_mode, "amical")} en {subject}.
Message de l'étudiant : {message}

{'Fournissez des réponses étape par étape.' if step_by_step else ''}
        ''',
            'en': f'''
You are an {mode_descriptions.get(tutor_mode, "friendly")} AI tutor for {subject}.
Student message: {message}

{'Provide step-by-step answers.' if step_by_step else ''}
        '''
        }
        
        return templates.get(language, templates['en'])
    
    def _parse_solution_response(self, response: str, language: str) -> Dict[str, Any]:
        """Parse model response into structured solution format"""
        
        # Simple parsing - in production, use more sophisticated parsing
        steps = []
        
        # Split response into steps (simple heuristic)
        lines = response.split('\n')
        current_step = {"explanation": "", "output": "", "step_number": 0}
        
        step_number = 1
        for line in lines:
            line = line.strip()
            if not line:
                continue
                
            # Look for step indicators
            if any(keyword in line.lower() for keyword in ['step', 'خطوة', 'étape']):
                if current_step["explanation"] or current_step["output"]:
                    steps.append(current_step)
                    step_number += 1
                current_step = {
                    "explanation": line,
                    "output": "",
                    "step_number": step_number
                }
            else:
                current_step["output"] += line + "\n"
        
        # Add the last step
        if current_step["explanation"] or current_step["output"]:
            steps.append(current_step)
        
        # Extract final answer (simplified)
        final_answer = response.split('\n')[-1].strip() if response else ""
        
        return {
            "steps": steps,
            "final_answer": final_answer,
            "confidence": 0.85,  # Mock confidence
            "language_detected": language
        }
    
    def get_loaded_models(self) -> List[str]:
        """Get list of loaded model names"""
        return [name for name, model in self.models.items() if model.is_loaded]
    
    def get_available_models(self) -> List[str]:
        """Get list of all configured models"""
        return list(self.config.models.keys())
    
    def get_model_stats(self) -> Dict[str, Any]:
        """Get statistics about loaded models"""
        stats = {}
        for name, model in self.models.items():
            stats[name] = model.get_config()
        return stats
    
    async def cleanup(self):
        """Cleanup all models"""
        logger.info("Cleaning up model manager...")
        
        for name, model in self.models.items():
            try:
                await model.cleanup()
                logger.info(f"Cleaned up model {name}")
            except Exception as e:
                logger.error(f"Error cleaning up model {name}: {e}")
        
        self.models.clear()
        self.primary_model = None
        self.fallback_model = None
        self.specialized_models.clear()