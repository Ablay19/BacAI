# `/pkg`

This directory contains library code that can be used by external applications. Packages in this directory should be designed with the intention that other projects might import and use them.

## Packages

### gcpvision
Package for interacting with Google Cloud Vision API to perform OCR and text detection in images.

Features:
- Text detection from image files
- Google Cloud authentication handling
- Simple API for document text extraction

### serpapi
Package for interacting with the SerpAPI service to perform web searches and retrieve organic search results.

Features:
- Search engine integration with SerpAPI
- Structured response parsing
- Configurable HTTP client with timeouts
- Support for custom search parameters