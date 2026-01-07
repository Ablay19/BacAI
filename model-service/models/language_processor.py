import re
from typing import Dict, Optional
from ..utils.logger import get_logger

logger = get_logger(__name__)

class LanguageProcessor:
    """Handles language detection and processing"""
    
    def __init__(self, config):
        self.config = config
        
    def detect_language(self, text: str) -> str:
        """Detect the primary language of text"""
        
        # Simple pattern-based language detection
        arabic_pattern = re.compile(r'[\u0600-\u06FF]')
        french_pattern = re.compile(r'[àâäçéèêëïîôöùûüÿœæ]', re.IGNORECASE)
        
        # Count characters for each language
        arabic_chars = len(arabic_pattern.findall(text))
        french_chars = len(french_pattern.findall(text))
        
        # Remove whitespace for percentage calculation
        clean_text = re.sub(r'\s', '', text)
        total_chars = len(clean_text)
        
        if total_chars == 0:
            return 'en'
        
        # Calculate percentages
        arabic_percentage = arabic_chars / total_chars
        french_percentage = french_chars / total_chars
        
        # Determine language based on thresholds
        if arabic_percentage > 0.3:
            return 'ar'
        elif french_percentage > 0.1:
            return 'fr'
        else:
            return 'en'
    
    def extract_subject(self, text: str, language: str) -> str:
        """Extract subject from text using keyword matching"""
        
        subject_keywords = {
            'ar': {
                'mathematics': [
                    'رياضيات', 'معادلة', 'حل', 'رقم', 'عدد', 'مساحة', 'حجم', 
                    'زاوية', 'مثلث', 'دائرة', 'مربع', 'جبر', 'هندسة'
                ],
                'sciences': [
                    'الفيزياء', 'الكيمياء', 'الأحياء', 'تجربة', 'مادة', 
                    'طاقة', 'قوة', 'حرارة', 'كهرباء', 'ذرة'
                ],
                'arabic': [
                    'قواعد', 'نحو', 'صرف', 'بلاغة', 'أدب', 'شعر', 'قصة',
                    'لغة عربية', 'قواعد النحو', 'صرف', 'بلاغ'
                ],
                'french': [
                    'لغة فرنسية', 'قواعد فرنسية', 'فرنسي', 'قواعد', '词汇'
                ],
                'english': [
                    'لغة إنجليزية', 'إنجليزي', 'english', 'grammar'
                ],
                'islamic_studies': [
                    'قرآن', 'حديث', 'شريعة', 'فقه', 'سنة', 'صلاة', 
                    'صوم', 'زكاة', 'حج', 'إسلام'
                ]
            },
            'fr': {
                'mathematics': [
                    'mathématiques', 'équation', 'résoudre', 'nombre', 'surface', 
                    'volume', 'angle', 'triangle', 'cercle', 'carré', 'algèbre'
                ],
                'sciences': [
                    'physique', 'chimie', 'biologie', 'expérience', 'matière',
                    'énergie', 'force', 'chaleur', 'électricité', 'atome'
                ],
                'arabic': [
                    'arabe', 'grammaire arabe', 'langue arabe', 'alphapet'
                ],
                'french': [
                    'français', 'grammaire française', 'conjugaison', 'littérature'
                ],
                'english': [
                    'anglais', 'grammar anglaise', 'vocabulary'
                ],
                'islamic_studies': [
                    'coran', 'islam', 'charia', 'droit islamique', 'fiqh'
                ]
            },
            'en': {
                'mathematics': [
                    'mathematics', 'equation', 'solve', 'number', 'area', 'volume',
                    'angle', 'triangle', 'circle', 'square', 'algebra', 'geometry'
                ],
                'sciences': [
                    'physics', 'chemistry', 'biology', 'experiment', 'matter',
                    'energy', 'force', 'heat', 'electricity', 'atom', 'science'
                ],
                'arabic': [
                    'arabic', 'arabic language', 'arabic grammar'
                ],
                'french': [
                    'french', 'french language', 'french grammar', 'french literature'
                ],
                'english': [
                    'english', 'english language', 'english grammar', 'literature'
                ],
                'islamic_studies': [
                    'quran', 'islam', 'sharia', 'fiqh', 'hadith', 'sunnah'
                ]
            }
        }
        
        # Get keywords for the detected language
        keywords = subject_keywords.get(language, subject_keywords['en'])
        
        # Convert text to lowercase for matching
        text_lower = text.lower()
        
        # Count matches for each subject
        scores = {}
        for subject, words in keywords.items():
            score = 0
            for word in words:
                # Use word boundaries for better matching
                pattern = r'\b' + re.escape(word.lower()) + r'\b'
                matches = re.findall(pattern, text_lower)
                score += len(matches)
            scores[subject] = score
        
        # Find subject with highest score
        best_subject = 'general'
        max_score = 0
        
        for subject, score in scores.items():
            if score > max_score:
                max_score = score
                best_subject = subject
        
        # Return 'general' if no strong matches
        if max_score == 0:
            return 'general'
        
        return best_subject
    
    def get_subject_display_name(self, subject: str, language: str) -> str:
        """Get display name for subject in given language"""
        
        display_names = {
            'mathematics': {
                'ar': 'الرياضيات',
                'fr': 'Mathématiques',
                'en': 'Mathematics'
            },
            'sciences': {
                'ar': 'العلوم',
                'fr': 'Sciences',
                'en': 'Sciences'
            },
            'arabic': {
                'ar': 'اللغة العربية',
                'fr': 'Arabe',
                'en': 'Arabic'
            },
            'french': {
                'ar': 'اللغة الفرنسية',
                'fr': 'Français',
                'en': 'French'
            },
            'english': {
                'ar': 'اللغة الإنجليزية',
                'fr': 'Anglais',
                'en': 'English'
            },
            'islamic_studies': {
                'ar': 'الدراسات الإسلامية',
                'fr': 'Études islamiques',
                'en': 'Islamic Studies'
            },
            'general': {
                'ar': 'عام',
                'fr': 'Général',
                'en': 'General'
            }
        }
        
        return display_names.get(subject, {}).get(language, subject)
    
    def validate_curriculum_alignment(
        self, 
        subject: str, 
        level: str, 
        language: str
    ) -> bool:
        """Validate subject/level/language combination against Mauritanian curriculum"""
        
        valid_combinations = {
            'secondary_basic': {
                'subjects': [
                    'mathematics', 'sciences', 'arabic', 'french', 
                    'english', 'islamic_studies'
                ],
                'languages': ['ar', 'fr', 'en']
            },
            'secondary_lycee': {
                'subjects': [
                    'mathematics', 'sciences', 'arabic', 'french', 
                    'english', 'islamic_studies'
                ],
                'languages': ['ar', 'fr', 'en']
            },
            'university': {
                'subjects': [
                    'mathematics', 'sciences', 'arabic', 'french', 
                    'english', 'islamic_studies'
                ],
                'languages': ['ar', 'fr', 'en']
            }
        }
        
        level_config = valid_combinations.get(level)
        if not level_config:
            return False
        
        return (subject in level_config['subjects'] and 
                language in level_config['languages'])
    
    def get_level_display_name(self, level: str, language: str) -> str:
        """Get display name for education level"""
        
        display_names = {
            'secondary_basic': {
                'ar': 'التعليم الثانوي الأساسي',
                'fr': 'Secondaire fondamental',
                'en': 'Secondary Basic'
            },
            'secondary_lycee': {
                'ar': 'التعليم الثانوي الثانوي',
                'fr': 'Secondaire lycée',
                'en': 'Secondary Lycée'
            },
            'university': {
                'ar': 'التعليم الجامعي',
                'fr': 'Université',
                'en': 'University'
            }
        }
        
        return display_names.get(level, {}).get(language, level)
    
    def preprocess_text(self, text: str, language: str) -> str:
        """Preprocess text based on language-specific rules"""
        
        # Remove extra whitespace
        text = re.sub(r'\s+', ' ', text.strip())
        
        # Language-specific preprocessing
        if language == 'ar':
            # Arabic diacritics handling (optional)
            text = re.sub(r'[\u064B-\u065F]', '', text)  # Remove diacritics
            
        elif language == 'fr':
            # French apostrophe and accent normalization
            text = text.replace(''', "'").replace(''', "'")
            
        elif language == 'en':
            # English-specific cleaning
            text = re.sub(r'[^a-zA-Z0-9\s.,!?;:()-]', ' ', text)
        
        return text
    
    def get_cultural_context_keywords(self, subject: str, language: str) -> list:
        """Get culturally relevant keywords for Mauritanian context"""
        
        cultural_keywords = {
            'ar': {
                'mathematics': ['نواكشوط', 'موريتانيا', 'أفريقيا'],
                'sciences': ['صحراء', 'الساحل', 'موريتانيا', 'بيئة محلية'],
                'arabic': ['موريتانيا', 'ثقافة', 'شعر موريتاني', 'أدب موريتاني'],
                'islamic_studies': ['موريتانيا', 'محظرة', 'علماء موريتانيا', 'تراث']
            },
            'fr': {
                'mathematics': ['Nouakchott', 'Mauritanie', 'Afrique'],
                'sciences': ['désert', 'Sahel', 'Mauritanie', 'environnement local'],
                'arabic': ['Mauritanie', 'culture', 'poésie mauritanienne'],
                'islamic_studies': ['Mauritanie', 'mahadra', 'savants mauritaniens']
            },
            'en': {
                'mathematics': ['Nouakchott', 'Mauritania', 'Africa'],
                'sciences': ['desert', 'Sahel', 'Mauritania', 'local environment'],
                'arabic': ['Mauritania', 'culture', 'Mauritanian poetry'],
                'islamic_studies': ['Mauritania', 'mahadra', 'Mauritanian scholars']
            }
        }
        
        return cultural_keywords.get(language, {}).get(subject, [])