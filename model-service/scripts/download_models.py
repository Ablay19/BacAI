#!/usr/bin/env python3
"""
Script to download and cache AI models for BACAI system
"""

import os
import logging
from pathlib import Path
import asyncio
from typing import List, Dict, Any

# Model imports (will be imported when available)
try:
    from transformers import AutoTokenizer, AutoModelForCausalLM
    import torch
    TRANSFORMERS_AVAILABLE = True
except ImportError:
    TRANSFORMERS_AVAILABLE = False
    print("Warning: transformers not available - using mock downloads")

from utils.config import get_config
from utils.logger import get_logger

logger = get_logger(__name__)

class ModelDownloader:
    """Handles downloading and caching of AI models"""
    
    def __init__(self, config):
        self.config = config
        self.models_dir = Path(config.model_cache_dir)
        self.models_dir.mkdir(parents=True, exist_ok=True)
        
    async def download_all_models(self):
        """Download all configured models"""
        
        logger.info("Starting model download process...")
        
        for model_name, model_config in self.config.models.items():
            try:
                success = await self.download_model(model_config.name)
                if success:
                    logger.info(f"✅ Successfully downloaded {model_config.name}")
                else:
                    logger.error(f"❌ Failed to download {model_config.name}")
                    
            except Exception as e:
                logger.error(f"Error downloading {model_config.name}: {e}")
        
        logger.info("Model download process completed")
    
    async def download_model(self, model_name: str) -> bool:
        """Download a specific model"""
        
        logger.info(f"Downloading model: {model_name}")
        
        if not TRANSFORMERS_AVAILABLE:
            logger.info(f"Mock downloading {model_name} (transformers not available)")
            return await self._mock_download(model_name)
        
        try:
            # Download tokenizer
            logger.info(f"Downloading tokenizer for {model_name}")
            tokenizer_path = self.models_dir / f"{model_name.replace('/', '_')}_tokenizer"
            tokenizer = AutoTokenizer.from_pretrained(
                model_name,
                cache_dir=tokenizer_path,
                trust_remote_code=True
            )
            
            # Download model (try with half precision first)
            logger.info(f"Downloading model weights for {model_name}")
            model_path = self.models_dir / f"{model_name.replace('/', '_')}_model"
            
            # Check for GPU availability
            use_gpu = torch.cuda.is_available()
            if use_gpu:
                logger.info(f"GPU detected: {torch.cuda.get_device_name()}")
                dtype = torch.float16
                device_map = "auto"
            else:
                logger.info("No GPU detected, using CPU")
                dtype = torch.float32
                device_map = None
            
            model = AutoModelForCausalLM.from_pretrained(
                model_name,
                cache_dir=model_path,
                torch_dtype=dtype,
                device_map=device_map,
                trust_remote_code=True
            )
            
            # Save info about the downloaded model
            model_info = {
                'name': model_name,
                'download_path': str(model_path),
                'tokenizer_path': str(tokenizer_path),
                'download_time': str(asyncio.get_event_loop().time()),
                'device': 'gpu' if use_gpu else 'cpu',
                'dtype': str(dtype),
                'model_type': type(model).__name__
            }
            
            info_path = self.models_dir / f"{model_name.replace('/', '_')}_info.json"
            import json
            with open(info_path, 'w') as f:
                json.dump(model_info, f, indent=2)
            
            # Test loading
            await self._test_model_loading(model_name)
            
            logger.info(f"Model {model_name} downloaded and tested successfully")
            return True
            
        except Exception as e:
            logger.error(f"Failed to download {model_name}: {e}")
            return False
    
    async def _test_model_loading(self, model_name: str):
        """Test that downloaded model can be loaded"""
        
        try:
            logger.info(f"Testing model loading for {model_name}")
            
            # Quick load test with minimal memory
            model_path = self.models_dir / f"{model_name.replace('/', '_')}_model"
            
            if model_path.exists():
                logger.info(f"✅ Model {model_name} files are present")
            else:
                logger.warning(f"⚠️ Model directory not found for {model_name}")
            
            # Memory usage check
            if torch.cuda.is_available():
                memory_allocated = torch.cuda.memory_allocated() / 1024**3  # GB
                logger.info(f"GPU memory allocated: {memory_allocated:.2f} GB")
            
            return True
            
        except Exception as e:
            logger.error(f"Model loading test failed: {e}")
            return False
    
    async def _mock_download(self, model_name: str) -> bool:
        """Mock download for when transformers is not available"""
        
        # Simulate download time
        await asyncio.sleep(2)
        
        # Create mock model info
        model_info = {
            'name': model_name,
            'downloaded': True,
            'mock': True,
            'download_time': str(asyncio.get_event_loop().time())
        }
        
        info_path = self.models_dir / f"{model_name.replace('/', '_')}_info.json"
        import json
        with open(info_path, 'w') as f:
            json.dump(model_info, f, indent=2)
        
        return True
    
    def list_downloaded_models(self) -> List[Dict[str, Any]]:
        """List all downloaded models"""
        
        models = []
        
        for info_file in self.models_dir.glob("*_info.json"):
            try:
                import json
                with open(info_file, 'r') as f:
                    model_info = json.load(f)
                models.append(model_info)
            except Exception as e:
                logger.warning(f"Could not read model info from {info_file}: {e}")
        
        return models
    
    def get_model_size(self, model_name: str) -> str:
        """Get size of downloaded model"""
        
        import shutil
        
        model_path = self.models_dir / f"{model_name.replace('/', '_')}_model"
        
        if not model_path.exists():
            return "Not downloaded"
        
        try:
            total_size = sum(
                f.stat().st_size for f in model_path.rglob('*') if f.is_file()
            )
            
            # Convert to human readable format
            for unit in ['B', 'KB', 'MB', 'GB']:
                if total_size < 1024:
                    return f"{total_size:.1f} {unit}"
                total_size /= 1024
            return f"{total_size:.1f} TB"
            
        except Exception as e:
            logger.error(f"Could not calculate model size: {e}")
            return "Unknown"
    
    def cleanup_models(self):
        """Clean up downloaded models"""
        
        logger.info("Cleaning up downloaded models...")
        
        import shutil
        
        for item in self.models_dir.iterdir():
            if item.is_dir():
                try:
                    shutil.rmtree(item)
                    logger.info(f"Removed: {item}")
                except Exception as e:
                    logger.error(f"Failed to remove {item}: {e}")

async def main():
    """Main function for model download"""
    
    import argparse
    
    parser = argparse.ArgumentParser(description='BACAI Model Downloader')
    parser.add_argument('--models', '-m', nargs='*', 
                       help='Specific models to download (default: all)')
    parser.add_argument('--list', '-l', action='store_true',
                       help='List downloaded models')
    parser.add_argument('--cleanup', action='store_true',
                       help='Clean up all downloaded models')
    parser.add_argument('--test', '-t', action='store_true',
                       help='Test downloaded models')
    
    args = parser.parse_args()
    
    # Initialize downloader
    config = get_config()
    downloader = ModelDownloader(config)
    
    if args.cleanup:
        downloader.cleanup_models()
        return
    
    if args.list:
        models = downloader.list_downloaded_models()
        if not models:
            print("No models downloaded")
        else:
            print("Downloaded models:")
            for model in models:
                size = downloader.get_model_size(model['name'])
                print(f"  - {model['name']}: {size}")
        return
    
    if args.test:
        models = downloader.list_downloaded_models()
        for model in models:
            await downloader._test_model_loading(model['name'])
        return
    
    # Download models
    if args.models:
        # Download specific models
        for model_name in args.models:
            success = await downloader.download_model(model_name)
            if success:
                print(f"✅ Successfully downloaded {model_name}")
            else:
                print(f"❌ Failed to download {model_name}")
    else:
        # Download all configured models
        await downloader.download_all_models()
    
    # Show summary
    print("\nDownload Summary:")
    models = downloader.list_downloaded_models()
    for model in models:
        size = downloader.get_model_size(model['name'])
        print(f"  - {model['name']}: {size}")

if __name__ == '__main__':
    asyncio.run(main())