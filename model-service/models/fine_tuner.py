import torch
from transformers import (
    AutoTokenizer, AutoModelForCausalLM, 
    TrainingArguments, Trainer,
    DataCollatorForLanguageModeling
)
from peft import LoraConfig, get_peft_model, TaskType
from datasets import Dataset, load_dataset
import json
import logging
from typing import Dict, Any, List
import os
from pathlib import Path

from utils.config import Config
from utils.logger import get_logger

logger = get_logger(__name__)

class FineTuner:
    """Handles fine-tuning of AI models on Mauritanian educational data"""
    
    def __init__(self, config: Config):
        self.config = config
        self.model = None
        self.tokenizer = None
        self.trainer = None
        
    async def load_model_for_finetuning(
        self, 
        model_name: str,
        dataset_path: str
    ) -> bool:
        """Load model and prepare for fine-tuning"""
        
        try:
            logger.info(f"Loading model {model_name} for fine-tuning")
            
            # Load tokenizer
            self.tokenizer = AutoTokenizer.from_pretrained(
                model_name,
                trust_remote_code=True,
                padding_side='right'
            )
            
            # Set padding token if not present
            if self.tokenizer.pad_token is None:
                self.tokenizer.pad_token = self.tokenizer.eos_token
            
            # Load model
            self.model = AutoModelForCausalLM.from_pretrained(
                model_name,
                trust_remote_code=True,
                torch_dtype=torch.float16,
                device_map="auto"
            )
            
            # Prepare model for training
            self.model = self._prepare_model_for_training(self.model)
            
            # Load and prepare dataset
            dataset = await self._load_dataset(dataset_path)
            processed_dataset = self._prepare_dataset(dataset)
            
            # Setup trainer
            self.trainer = self._setup_trainer(processed_dataset)
            
            logger.info(f"Model {model_name} loaded successfully for fine-tuning")
            return True
            
        except Exception as e:
            logger.error(f"Failed to load model {model_name}: {e}")
            return False
    
    def _prepare_model_for_training(self, model):
        """Prepare model for LoRA fine-tuning"""
        
        # LoRA configuration for efficient fine-tuning
        lora_config = LoraConfig(
            task_type=TaskType.CAUSAL_LM,
            inference_mode=False,
            r=8,  # rank
            lora_alpha=32,
            lora_dropout=0.1,
            target_modules=["q_proj", "v_proj", "k_proj", "o_proj"],
            bias="none",
        )
        
        # Apply LoRA
        model = get_peft_model(model, lora_config)
        model.print_trainable_parameters()
        
        return model
    
    async def _load_dataset(self, dataset_path: str) -> Dataset:
        """Load training dataset"""
        
        # Support multiple formats
        if dataset_path.endswith('.json'):
            dataset = load_dataset('json', data_files=dataset_path)
        elif dataset_path.endswith('.jsonl'):
            dataset = load_dataset('json', data_files=dataset_path)
        else:
            raise ValueError(f"Unsupported dataset format: {dataset_path}")
        
        logger.info(f"Loaded dataset: {dataset}")
        return dataset['train']  # Use train split
    
    def _prepare_dataset(self, dataset: Dataset) -> Dataset:
        """Prepare dataset for training"""
        
        def format_examples(examples):
            """Format examples for instruction fine-tuning"""
            
            formatted_texts = []
            
            for exercise, solution, subject, level, language in zip(
                examples['exercise'],
                examples['solution'], 
                examples['subject'],
                examples['level'],
                examples['language']
            ):
                # Create instruction format
                if language == 'ar':
                    instruction = f'''حل المسألة التالية في مادة {subject} للمستوى {level} خطوة بخطوة:

المسألة: {exercise}

الحل: {solution}'''
                elif language == 'fr':
                    instruction = f'''Résolvez l'exercice suivant en {subject} pour le niveau {level} étape par étape:

Exercice: {exercise}

Solution: {solution}'''
                else:
                    instruction = f'''Solve the following exercise in {subject} for {level} level step by step:

Exercise: {exercise}

Solution: {solution}'''
                
                formatted_texts.append(instruction)
            
            return {
                'text': formatted_texts
            }
        
        # Tokenize dataset
        def tokenize_function(examples):
            return self.tokenizer(
                examples['text'],
                truncation=True,
                padding='max_length',
                max_length=2048,
                return_tensors='pt'
            )
        
        # Apply formatting and tokenization
        dataset = dataset.map(format_examples, batched=True)
        dataset = dataset.map(tokenize_function, batched=True)
        
        # Remove unused columns
        dataset = dataset.remove_columns([
            'exercise', 'solution', 'subject', 'level', 'language', 'text'
        ])
        
        logger.info(f"Prepared dataset with {len(dataset)} examples")
        return dataset
    
    def _setup_trainer(self, dataset: Dataset) -> Trainer:
        """Setup training configuration and trainer"""
        
        # Training arguments
        training_args = TrainingArguments(
            output_dir="./models/fine_tuned",
            overwrite_output_dir=True,
            num_train_epochs=3,
            per_device_train_batch_size=2,
            gradient_accumulation_steps=4,
            warmup_steps=100,
            max_steps=1000,
            learning_rate=2e-4,
            fp16=True,  # Use mixed precision
            logging_steps=10,
            save_steps=100,
            eval_steps=100,
            evaluation_strategy="steps",
            load_best_model_at_end=True,
            metric_for_best_model="eval_loss",
            greater_is_better=False,
            report_to="none",  # Disable wandb/mlflow for now
        )
        
        # Data collator
        data_collator = DataCollatorForLanguageModeling(
            tokenizer=self.tokenizer,
            mlm=False  # We're doing causal LM, not masked LM
        )
        
        # Create trainer
        trainer = Trainer(
            model=self.model,
            args=training_args,
            train_dataset=dataset,
            eval_dataset=dataset.select(range(min(100, len(dataset)))),  # Small eval set
            data_collator=data_collator,
        )
        
        return trainer
    
    async def train_model(self, save_path: str = "./models/fine_tuned") -> bool:
        """Fine-tune the model"""
        
        try:
            logger.info("Starting model fine-tuning")
            
            # Train the model
            train_result = self.trainer.train()
            
            # Save the model
            self.trainer.save_model(save_path)
            self.tokenizer.save_pretrained(save_path)
            
            # Save training logs
            log_history = self.trainer.state.log_history
            with open(f"{save_path}/training_log.json", "w") as f:
                json.dump(log_history, f, indent=2)
            
            logger.info(f"Fine-tuning completed. Model saved to {save_path}")
            return True
            
        except Exception as e:
            logger.error(f"Fine-tuning failed: {e}")
            return False
    
    async def evaluate_model(self, test_dataset_path: str) -> Dict[str, float]:
        """Evaluate fine-tuned model"""
        
        try:
            logger.info("Evaluating fine-tuned model")
            
            # Load test dataset
            test_dataset = await self._load_dataset(test_dataset_path)
            processed_dataset = self._prepare_dataset(test_dataset)
            
            # Evaluate
            eval_result = self.trainer.evaluate()
            
            metrics = {
                "eval_loss": eval_result.get("eval_loss", 0.0),
                "perplexity": torch.exp(torch.tensor(eval_result.get("eval_loss", 0.0))).item(),
                "eval_samples": eval_result.get("eval_samples", 0)
            }
            
            logger.info(f"Evaluation results: {metrics}")
            return metrics
            
        except Exception as e:
            logger.error(f"Evaluation failed: {e}")
            return {}
    
    def cleanup(self):
        """Clean up resources"""
        if self.model:
            del self.model
        if self.tokenizer:
            del self.tokenizer
        if self.trainer:
            del self.trainer
        
        # Clear GPU cache
        if torch.cuda.is_available():
            torch.cuda.empty_cache()

class CurriculumTrainer:
    """Specialized trainer for curriculum-aligned fine-tuning"""
    
    def __init__(self, config: Config):
        self.config = config
        self.fine_tuner = FineTuner(config)
        
    async def train_on_curriculum(
        self, 
        subject: str,
        level: str,
        curriculum_data_path: str
    ) -> bool:
        """Train model on specific curriculum data"""
        
        logger.info(f"Training {subject}/{level} model on curriculum data")
        
        # Load curriculum-specific data
        curriculum_data = self._load_curriculum_data(curriculum_data_path)
        
        # Filter data for subject and level
        filtered_data = [
            item for item in curriculum_data 
            if item.get('subject') == subject and item.get('level') == level
        ]
        
        if not filtered_data:
            logger.warning(f"No data found for {subject}/{level}")
            return False
        
        # Create temporary dataset
        temp_dataset_path = f"/tmp/{subject}_{level}_dataset.json"
        with open(temp_dataset_path, 'w', encoding='utf-8') as f:
            json.dump(filtered_data, f, ensure_ascii=False, indent=2)
        
        # Train model
        model_name = self.config.models.get('primary', {}).name or 'Qwen/Qwen2-8B-Instruct'
        success = await self.fine_tuner.load_model_for_finetuning(
            model_name, temp_dataset_path
        )
        
        if success:
            save_path = f"./models/fine_tuned/{subject}_{level}"
            success = await self.fine_tuner.train_model(save_path)
        
        # Cleanup
        os.remove(temp_dataset_path)
        self.fine_tuner.cleanup()
        
        return success
    
    def _load_curriculum_data(self, data_path: str) -> List[Dict[str, Any]]:
        """Load curriculum data from file"""
        
        try:
            with open(data_path, 'r', encoding='utf-8') as f:
                if data_path.endswith('.json'):
                    return json.load(f)
                elif data_path.endswith('.jsonl'):
                    data = []
                    for line in f:
                        data.append(json.loads(line))
                    return data
        except Exception as e:
            logger.error(f"Failed to load curriculum data: {e}")
            return []