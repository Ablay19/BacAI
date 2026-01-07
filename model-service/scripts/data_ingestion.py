#!/usr/bin/env python3
"""
Data ingestion pipeline for BACAI Mauritanian educational content
"""

import json
import asyncio
import logging
from pathlib import Path
from typing import Dict, List, Any, Optional
from datetime import datetime
import hashlib

from utils.config import get_config
from utils.logger import get_logger
from models.language_processor import LanguageProcessor

logger = get_logger(__name__)

class DataIngestionPipeline:
    """Pipeline for processing and ingesting educational content"""
    
    def __init__(self, config):
        self.config = config
        self.language_processor = LanguageProcessor(config)
        self.processed_count = 0
        self.failed_count = 0
        self.duplicates_count = 0
        
    async def process_file(
        self, 
        file_path: str, 
        output_path: str,
        validate_curriculum: bool = True
    ) -> Dict[str, Any]:
        """Process a single file of educational content"""
        
        logger.info(f"Processing file: {file_path}")
        
        try:
            # Load raw data
            raw_data = await self._load_file(file_path)
            
            # Process each item
            processed_items = []
            for item in raw_data:
                try:
                    processed_item = await self._process_item(item, validate_curriculum)
                    if processed_item:
                        processed_items.append(processed_item)
                except Exception as e:
                    logger.warning(f"Failed to process item: {e}")
                    self.failed_count += 1
            
            # Remove duplicates
            unique_items = await self._remove_duplicates(processed_items)
            self.duplicates_count = len(processed_items) - len(unique_items)
            
            # Save processed data
            await self._save_processed_data(unique_items, output_path)
            
            # Generate report
            report = self._generate_processing_report(file_path, unique_items)
            
            self.processed_count = len(unique_items)
            logger.info(f"Successfully processed {len(unique_items)} items from {file_path}")
            
            return report
            
        except Exception as e:
            logger.error(f"Failed to process file {file_path}: {e}")
            raise
    
    async def _load_file(self, file_path: str) -> List[Dict[str, Any]]:
        """Load data from file based on format"""
        
        file_path = Path(file_path)
        
        if file_path.suffix.lower() == '.json':
            with open(file_path, 'r', encoding='utf-8') as f:
                return json.load(f)
                
        elif file_path.suffix.lower() == '.jsonl':
            data = []
            with open(file_path, 'r', encoding='utf-8') as f:
                for line in f:
                    if line.strip():
                        data.append(json.loads(line))
            return data
            
        else:
            raise ValueError(f"Unsupported file format: {file_path.suffix}")
    
    async def _process_item(
        self, 
        item: Dict[str, Any], 
        validate_curriculum: bool
    ) -> Optional[Dict[str, Any]]:
        """Process a single educational item"""
        
        # Required fields validation
        if not all(key in item for key in ['exercise', 'subject']):
            logger.warning(f"Item missing required fields: {item}")
            return None
        
        # Auto-detect language if not provided
        if 'language' not in item:
            item['language'] = self.language_processor.detect_language(
                item.get('exercise', '') + item.get('concept', '')
            )
        
        # Auto-detect subject if not provided
        if 'subject' not in item or item['subject'] == 'general':
            text = item.get('exercise', '') + item.get('concept', '') + item.get('message', '')
            item['subject'] = self.language_processor.extract_subject(
                text, item['language']
            )
        
        # Set defaults
        item.setdefault('level', 'secondary_lycee')
        item.setdefault('difficulty', 'medium')
        item.setdefault('tags', [])
        
        # Validate curriculum alignment
        if validate_curriculum:
            if not self.language_processor.validate_curriculum_alignment(
                item['subject'], item['level'], item['language']
            ):
                logger.warning(f"Invalid curriculum combination: {item}")
                return None
        
        # Clean and preprocess text
        item['exercise'] = self.language_processor.preprocess_text(
            item.get('exercise', ''), item['language']
        )
        
        item['concept'] = self.language_processor.preprocess_text(
            item.get('concept', ''), item['language']
        )
        
        # Add metadata
        item['processed_at'] = datetime.now().isoformat()
        item['content_hash'] = self._generate_content_hash(item)
        
        # Quality scoring
        item['quality_score'] = self._calculate_quality_score(item)
        
        return item
    
    async def _remove_duplicates(
        self, 
        items: List[Dict[str, Any]]
    ) -> List[Dict[str, Any]]:
        """Remove duplicate items based on content hash"""
        
        seen_hashes = set()
        unique_items = []
        
        for item in items:
            content_hash = item.get('content_hash', '')
            if content_hash not in seen_hashes:
                seen_hashes.add(content_hash)
                unique_items.append(item)
        
        return unique_items
    
    async def _save_processed_data(
        self, 
        items: List[Dict[str, Any]], 
        output_path: str
    ):
        """Save processed data to output file"""
        
        output_path = Path(output_path)
        output_path.parent.mkdir(parents=True, exist_ok=True)
        
        # Save as JSONL for streaming processing
        with open(output_path, 'w', encoding='utf-8') as f:
            for item in items:
                f.write(json.dumps(item, ensure_ascii=False) + '\n')
        
        # Also save as JSON for easy inspection
        json_path = output_path.with_suffix('.json')
        with open(json_path, 'w', encoding='utf-8') as f:
            json.dump(items, f, ensure_ascii=False, indent=2)
        
        logger.info(f"Saved {len(items)} items to {output_path}")
    
    def _generate_content_hash(self, item: Dict[str, Any]) -> str:
        """Generate hash of item content for deduplication"""
        
        content = (
            item.get('exercise', '') + 
            item.get('concept', '') + 
            item.get('message', '')
        )
        
        return hashlib.sha256(content.encode('utf-8')).hexdigest()
    
    def _calculate_quality_score(self, item: Dict[str, Any]) -> float:
        """Calculate quality score for processed item"""
        
        score = 0.0
        
        # Content length (longer is generally better)
        content_length = len(
            item.get('exercise', '') + 
            item.get('concept', '') + 
            item.get('solution', '')
        )
        score += min(content_length / 1000, 0.3) * 0.2  # Max 0.2
        
        # Has solution
        if item.get('solution') and len(item['solution']) > 10:
            score += 0.2
        
        # Has tags
        if item.get('tags') and len(item['tags']) > 0:
            score += 0.1
        
        # Valid curriculum combination
        if self.language_processor.validate_curriculum_alignment(
            item.get('subject', ''), 
            item.get('level', ''), 
            item.get('language', '')
        ):
            score += 0.3
        
        # Language quality indicators
        language = item.get('language', '')
        if language in ['ar', 'fr', 'en']:
            score += 0.2
        
        return min(score, 1.0)
    
    def _generate_processing_report(
        self, 
        input_path: str, 
        processed_items: List[Dict[str, Any]]
    ) -> Dict[str, Any]:
        """Generate processing report"""
        
        # Statistics by subject
        subject_counts = {}
        language_counts = {}
        level_counts = {}
        quality_scores = []
        
        for item in processed_items:
            subject = item.get('subject', 'unknown')
            language = item.get('language', 'unknown')
            level = item.get('level', 'unknown')
            quality = item.get('quality_score', 0.0)
            
            subject_counts[subject] = subject_counts.get(subject, 0) + 1
            language_counts[language] = language_counts.get(language, 0) + 1
            level_counts[level] = level_counts.get(level, 0) + 1
            quality_scores.append(quality)
        
        return {
            'input_file': input_path,
            'processed_at': datetime.now().isoformat(),
            'statistics': {
                'total_processed': self.processed_count,
                'failed': self.failed_count,
                'duplicates_removed': self.duplicates_count,
                'success_rate': self.processed_count / (self.processed_count + self.failed_count) if (self.processed_count + self.failed_count) > 0 else 0.0,
                'average_quality_score': sum(quality_scores) / len(quality_scores) if quality_scores else 0.0
            },
            'distribution': {
                'by_subject': subject_counts,
                'by_language': language_counts,
                'by_level': level_counts
            }
        }

class MauritanianCurriculumProcessor(DataIngestionPipeline):
    """Specialized processor for Mauritanian curriculum data"""
    
    def __init__(self, config):
        super().__init__(config)
        
        # Mauritanian curriculum standards
        self.curriculum_standards = {
            'subjects': [
                'mathematics', 'arabic', 'french', 'english', 
                'sciences', 'islamic_studies'
            ],
            'levels': [
                'secondary_basic', 'secondary_lycee', 'university'
            ],
            'certificates': {
                'secondary_basic': 'BEPC',
                'secondary_lycee': 'Baccalaureate',
                'university': 'Degree'
            }
        }
    
    async def _process_item(
        self, 
        item: Dict[str, Any], 
        validate_curriculum: bool
    ) -> Optional[Dict[str, Any]]:
        """Process item with Mauritanian curriculum specifics"""
        
        # Standard processing
        processed_item = await super()._process_item(item, validate_curriculum)
        
        if not processed_item:
            return None
        
        # Add Mauritanian context
        processed_item['curriculum'] = {
            'country': 'mauritania',
            'system': 'baccalaureate',
            'certificate': self.curriculum_standards['certificates'].get(
                processed_item.get('level'), 'Unknown'
            )
        }
        
        # Cultural relevance scoring
        processed_item['cultural_relevance'] = self._calculate_cultural_relevance(
            processed_item
        )
        
        return processed_item
    
    def _calculate_cultural_relevance(self, item: Dict[str, Any]) -> float:
        """Calculate cultural relevance score for Mauritanian context"""
        
        score = 0.0
        content = (
            item.get('exercise', '') + 
            item.get('concept', '') + 
            item.get('solution', '')
        ).lower()
        
        # Mauritanian cultural keywords
        mauritanian_keywords = [
            'mauritania', 'mauritanie', 'موريتانيا', 
            'nouakchott', 'نواكشوط',
            'sahara', 'صحراء',
            'mahdara', 'محظرة',
            'arabic', 'اللغة العربية',
            'islamic', 'إسلامي'
        ]
        
        # Count cultural references
        cultural_count = sum(1 for keyword in mauritanian_keywords if keyword in content)
        score += min(cultural_count / 3, 0.3) * 0.3  # Max 0.3
        
        # Language relevance (Arabic is primary in Mauritania)
        if item.get('language') == 'ar':
            score += 0.3
        elif item.get('language') == 'fr':
            score += 0.2
        
        # Subject relevance to curriculum
        subject = item.get('subject', '')
        if subject in self.curriculum_standards['subjects']:
            score += 0.4
        
        return min(score, 1.0)

async def main():
    """Main function for data ingestion"""
    
    import argparse
    
    parser = argparse.ArgumentParser(description='BACAI Data Ingestion Pipeline')
    parser.add_argument('--input', '-i', required=True, help='Input file or directory')
    parser.add_argument('--output', '-o', required=True, help='Output directory')
    parser.add_argument('--curriculum', '-c', action='store_true', 
                       help='Validate against Mauritanian curriculum')
    parser.add_argument('--recursive', '-r', action='store_true',
                       help='Process directory recursively')
    
    args = parser.parse_args()
    
    # Initialize processor
    config = get_config()
    
    if args.curriculum:
        processor = MauritanianCurriculumProcessor(config)
    else:
        processor = DataIngestionPipeline(config)
    
    # Process input
    input_path = Path(args.input)
    output_path = Path(args.output)
    
    if input_path.is_file():
        # Process single file
        output_file = output_path / f"processed_{input_path.name}"
        report = await processor.process_file(str(input_path), str(output_file))
        print(json.dumps(report, indent=2))
        
    elif input_path.is_dir():
        # Process directory
        if args.recursive:
            files = list(input_path.rglob('*'))
        else:
            files = list(input_path.glob('*'))
        
        # Filter supported files
        supported_files = [f for f in files if f.suffix.lower() in ['.json', '.jsonl']]
        
        if not supported_files:
            logger.error(f"No supported files found in {input_path}")
            return
        
        logger.info(f"Found {len(supported_files)} files to process")
        
        reports = []
        for file_path in supported_files:
            output_file = output_path / f"processed_{file_path.name}"
            try:
                report = await processor.process_file(str(file_path), str(output_file))
                reports.append(report)
            except Exception as e:
                logger.error(f"Failed to process {file_path}: {e}")
        
        # Generate summary report
        summary_report = {
            'summary': {
                'files_processed': len(reports),
                'total_items': sum(r['statistics']['total_processed'] for r in reports),
                'total_failed': sum(r['statistics']['failed'] for r in reports),
                'total_duplicates': sum(r['statistics']['duplicates_removed'] for r in reports)
            },
            'reports': reports
        }
        
        summary_file = output_path / 'summary_report.json'
        with open(summary_file, 'w', encoding='utf-8') as f:
            json.dump(summary_report, f, ensure_ascii=False, indent=2)
        
        print(f"Processing complete. Summary saved to {summary_file}")
        
    else:
        logger.error(f"Input path does not exist: {input_path}")

if __name__ == '__main__':
    asyncio.run(main())