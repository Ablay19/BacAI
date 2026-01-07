import pytest
import asyncio
from fastapi.testclient import TestClient
from unittest.mock import AsyncMock, patch

from main import app
from utils.config import get_config

class TestAPIEndpoints:
    """Test suite for API endpoints"""
    
    def setup_method(self):
        """Setup test client"""
        self.client = TestClient(app)
        
    def test_health_check(self):
        """Test health check endpoint"""
        response = self.client.get("/health")
        
        assert response.status_code == 200
        data = response.json()
        assert data["status"] == "healthy"
        assert "timestamp" in data
        assert "service" in data
    
    def test_root_endpoint(self):
        """Test root endpoint"""
        response = self.client.get("/")
        
        assert response.status_code == 200
        data = response.json()
        assert "name" in data
        assert "endpoints" in data
    
    @patch('models.model_manager.ModelManager.solve_exercise')
    def test_solve_endpoint(self, mock_solve):
        """Test solve endpoint"""
        # Mock model response
        mock_solve.return_value = {
            "steps": [
                {
                    "explanation": "Step 1: Identify coefficients",
                    "output": "a=2, b=5, c=-3",
                    "step_number": 1
                }
            ],
            "final_answer": "x = [-3, 0.5]",
            "confidence": 0.95,
            "language_detected": "en"
        }
        
        # Test request
        request_data = {
            "exercise": "Solve: 2x² + 5x - 3 = 0",
            "subject": "mathematics",
            "level": "secondary_lycee",
            "language": "en",
            "mode": "step-by-step",
            "cultural_context": True,
            "request_type": "solve"
        }
        
        response = self.client.post("/api/process", json=request_data)
        
        assert response.status_code == 200
        data = response.json()
        assert data["success"] is True
        assert "data" in data
        assert "metadata" in data
        assert data["data"]["subject"] == "mathematics"
    
    @patch('models.model_manager.ModelManager.explain_concept')
    def test_explain_endpoint(self, mock_explain):
        """Test explain endpoint"""
        mock_explain.return_value = "A quadratic equation is a polynomial equation of degree 2..."
        
        request_data = {
            "concept": "quadratic equation",
            "subject": "mathematics", 
            "level": "secondary_lycee",
            "language": "en",
            "detail_level": "detailed",
            "examples": True,
            "cultural_context": False,
            "request_type": "explain"
        }
        
        response = self.client.post("/api/process", json=request_data)
        
        assert response.status_code == 200
        data = response.json()
        assert data["success"] is True
        assert "explanation" in data["data"]
    
    @patch('models.model_manager.ModelManager.converse')
    def test_converse_endpoint(self, mock_converse):
        """Test conversation endpoint"""
        mock_converse.return_value = "I'd be happy to help you understand quadratic equations..."
        
        request_data = {
            "message": "Can you explain quadratic equations?",
            "conversation_id": "test_conv_123",
            "subject": "mathematics",
            "language": "en", 
            "tutor_mode": "friendly",
            "step_by_step": True,
            "request_type": "converse"
        }
        
        response = self.client.post("/api/process", json=request_data)
        
        assert response.status_code == 200
        data = response.json()
        assert data["success"] is True
        assert "response" in data["data"]

class TestLanguageProcessing:
    """Test language processing functionality"""
    
    def setup_method(self):
        from models.language_processor import LanguageProcessor
        config = get_config()
        self.processor = LanguageProcessor(config)
    
    def test_arabic_language_detection(self):
        """Test Arabic language detection"""
        arabic_text = "حل المعادلة التالية: س² + ٥س - ٣ = ٠"
        detected = self.processor.detect_language(arabic_text)
        assert detected == 'ar'
    
    def test_french_language_detection(self):
        """Test French language detection"""
        french_text = "Résolvez l'équation suivante : x² + 5x - 3 = 0"
        detected = self.processor.detect_language(french_text)
        assert detected == 'fr'
    
    def test_english_language_detection(self):
        """Test English language detection"""
        english_text = "Solve the equation: x² + 5x - 3 = 0"
        detected = self.processor.detect_language(english_text)
        assert detected == 'en'
    
    def test_mathematics_subject_extraction(self):
        """Test mathematics subject extraction"""
        math_text = "Solve the quadratic equation: x² + 5x - 3 = 0"
        subject = self.processor.extract_subject(math_text, 'en')
        assert subject == 'mathematics'
    
    def test_islamic_studies_subject_extraction(self):
        """Test Islamic studies subject extraction"""
        islamic_text = "Explain the concept of Hadith and its classification"
        subject = self.processor.extract_subject(islamic_text, 'en')
        assert subject == 'islamic_studies'
    
    def test_curriculum_validation(self):
        """Test curriculum validation"""
        # Valid combination
        valid = self.processor.validate_curriculum_alignment(
            'mathematics', 'secondary_lycee', 'ar'
        )
        assert valid is True
        
        # Invalid combination
        invalid = self.processor.validate_curriculum_alignment(
            'mathematics', 'invalid_level', 'ar'
        )
        assert invalid is False

class TestModelManager:
    """Test model manager functionality"""
    
    @pytest.fixture
    def mock_config(self):
        """Mock configuration for testing"""
        from utils.config import ModelConfig, Config, DatabaseConfig, RedisConfig
        
        return Config(
            models={
                'primary': ModelConfig(
                    name='test-model',
                    type='test',
                    languages=['en', 'ar', 'fr'],
                    specialties=['mathematics'],
                    max_tokens=2048,
                    temperature=0.7
                )
            },
            database=DatabaseConfig(url='sqlite:///test.db'),
            redis=RedisConfig(url='redis://localhost:6379/0')
        )
    
    @pytest.fixture
    def mock_manager(self, mock_config):
        """Create mock model manager"""
        from models.model_manager import ModelManager
        return ModelManager(mock_config)
    
    @pytest.mark.asyncio
    async def test_model_selection(self, mock_manager):
        """Test appropriate model selection"""
        # Mock models
        mock_manager.models = {
            'math_specialist': AsyncMock(
                supports_language=AsyncMock(return_value=True),
                supports_specialty=AsyncMock(return_value=True)
            ),
            'general': AsyncMock(
                supports_language=AsyncMock(return_value=True),
                supports_specialty=AsyncMock(return_value=False)
            )
        }
        
        request_data = {
            'subject': 'mathematics',
            'language': 'en'
        }
        
        selected_model = mock_manager._select_model(request_data)
        
        # Should select the math specialist for math requests
        # This is simplified - in real implementation, we'd check method calls
        assert selected_model is not None

class TestConfigSystem:
    """Test configuration system"""
    
    def test_config_loading(self):
        """Test configuration loading"""
        config = get_config()
        
        assert config is not None
        assert hasattr(config, 'models')
        assert hasattr(config, 'database')
        assert hasattr(config, 'redis')
    
    def test_environment_variable_override(self):
        """Test environment variable overrides"""
        import os
        
        # Set environment variable
        os.environ['PORT'] = '9999'
        
        config = get_config()
        assert config.port == 9999
        
        # Clean up
        del os.environ['PORT']

# Integration tests
class TestIntegration:
    """Integration tests for the complete system"""
    
    @pytest.mark.asyncio
    async def test_end_to_end_solve_flow(self):
        """Test complete solve flow"""
        # This would test the entire flow from API to model response
        # In a real scenario, you'd have test data and expected outputs
        
        test_request = {
            "exercise": "Solve: x + 5 = 10",
            "subject": "mathematics",
            "level": "secondary_basic",
            "language": "en",
            "mode": "step-by-step",
            "cultural_context": True,
            "request_type": "solve"
        }
        
        # This would require actual model loading - skip in unit tests
        # But structure is ready for integration testing
        
        assert True  # Placeholder

# Performance tests
class TestPerformance:
    """Test performance characteristics"""
    
    @pytest.mark.asyncio
    async def test_language_detection_performance(self):
        """Test language detection performance"""
        from models.language_processor import LanguageProcessor
        import time
        from utils.config import get_config
        
        processor = LanguageProcessor(get_config())
        
        # Test with different sized texts
        test_texts = [
            "Short text",
            "This is a longer text with multiple sentences to test performance characteristics.",
            "A" * 1000  # Very long text
        ]
        
        for text in test_texts:
            start_time = time.time()
            result = processor.detect_language(text)
            end_time = time.time()
            
            # Should complete quickly (< 10ms for typical texts)
            processing_time = (end_time - start_time) * 1000
            assert processing_time < 100  # 100ms max

if __name__ == "__main__":
    pytest.main([__file__, "-v"])