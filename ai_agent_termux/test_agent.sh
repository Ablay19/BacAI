#!/bin/bash

# Test script for AI Agent
echo "Testing AI Agent for Termux"

# Create a test directory
mkdir -p ~/test_docs

# Create a sample text file for testing
cat > ~/test_docs/sample.txt << EOF
This is a sample document for testing the AI Agent.
The agent should be able to process this document,
extract its content, generate embeddings, and
provide a summary using either local or cloud LLMs.
EOF

echo "Created sample test document"

# Run the AI agent scan command
echo "Running AI agent scan..."
go run main.go scan

echo "Test completed. Check the processed_data directory for output."