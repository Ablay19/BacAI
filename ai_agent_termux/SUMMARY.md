# A.I.D.A - Android Intelligent Document Analyst
## Project Summary

## Overview
A.I.D.A. (Android Intelligent Document Analyst) is an autonomous, hybrid AI agent designed to operate fully within a Termux environment on Android devices. It can discover, scan, extract, analyze, categorize, summarize, index, and semantically search through various file types including images, PDFs, and text documents.

## Key Features Implemented

1. **File Discovery & Classification**
   - Recursive scanning of specified directories
   - Automatic classification of files by type (image, PDF, text, other)

2. **Multi-format Preprocessing**
   - OCR processing for images using Tesseract
   - PDF text extraction using PyPDF2
   - Direct text file reading

3. **Semantic Analysis**
   - Embedding generation using sentence-transformers
   - FAISS vector database integration for semantic search

4. **Intelligent Processing**
   - Local LLM processing with LLaMA.cpp
   - Cloud LLM fallback (Gemini, Ollama Cloud, ClawDBot, AIChat)
   - Automated summarization, categorization, and tagging

5. **Structured Output**
   - Comprehensive JSON output with extracted text, metadata, and processing results
   - Timestamped output files for audit trail

## Technical Architecture

### Core Components
- **Go Orchestrator**: Central control unit managing workflow and component coordination
- **Python Preprocessors**: Specialized scripts for OCR, PDF parsing, and embedding generation
- **LLM Processors**: Both local (LLaMA.cpp) and cloud-based (Gemini, Ollama, etc.) processing
- **Vector Database**: FAISS integration for semantic search capabilities

### Data Flow
1. File discovery in configured directories
2. File type classification
3. Format-specific preprocessing
4. Semantic embedding generation
5. LLM-based analysis (summarization, categorization, tagging)
6. Structured JSON output storage

## Technologies Used
- **Go**: Main orchestration language
- **Python**: Preprocessing and specialized tasks
- **Tesseract OCR**: Image text extraction
- **PyPDF2**: PDF text extraction
- **Sentence Transformers**: Semantic embeddings
- **FAISS**: Vector database for semantic search
- **LLaMA.cpp**: Local LLM processing
- **Various Cloud LLM APIs**: Gemini, Ollama Cloud, ClawDBot, AIChat

## Current Capabilities Demonstrated
- Successfully processes text files, extracting content and generating JSON output
- Handles multiple files in a single scan operation
- Integrates local LLM processing (with cloud fallback)
- Produces structured, timestamped output files

## Future Enhancements
- Termux-Tasker integration for scheduled processing
- Termux-Widget for home screen actions
- Shizuko dashboard integration
- Advanced image captioning with BLIP-2
- Enhanced categorization and tagging
- Improved error handling and logging
- Batch processing optimizations
- Multimodal processing for complex document types

## Setup and Usage
The agent is ready to use on any Termux-enabled Android device with the required dependencies installed. Users can simply run:
```
./ai_agent scan
```
to process all files in the configured directories.

## Conclusion
A.I.D.A. represents a comprehensive solution for autonomous document analysis on mobile devices, combining the power of local processing with cloud-based enhancements to deliver a seamless user experience.