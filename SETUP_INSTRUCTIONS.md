# A.I.D.A. Setup Instructions
## For Existing aichat and Ollama Installations

## Prerequisites
Make sure you have the following tools installed on your system:
1. **Go** (1.20 or later)
2. **Python** (3.8 or later) with required packages
3. **aichat** - Already installed as mentioned
4. **Ollama** - Already installed as mentioned

## Installing Required Python Packages
```bash
# Install Python dependencies
pip install pytesseract Pillow PyPDF2 sentence-transformers faiss-cpu

# Install Tesseract OCR engine and language packs
pkg install tesseract tesseract-ocr-eng  # Adjust for your system/package manager
```

## Setting Up the AI Agent

### 1. Clone or Copy the Project
```bash
# If you're copying from this environment:
cp -r /home/abdoullahelvogani/organizer/ai_agent_termux ~/ai_agent_termux

# Or clone from repository (when available):
# git clone <repository-url>
```

### 2. Build the AI Agent
```bash
cd ~/ai_agent_termux
go build -o ai_agent
```

### 3. Configure Environment Variables (Optional)
If you want to use cloud APIs, set these environment variables:
```bash
# For Google Gemini
export GEMINI_API_KEY="your_gemini_api_key"

# For Anthropic Claude (using ClawDBot API key field)
export CLAWDBOT_API_KEY="your_claude_api_key"

# For Ollama model selection (optional)
export OLLAMA_CLOUD_API_KEY="llama2"  # or "mistral", "phi", etc.
```

## Using Your Existing Tools

### aichat Integration
The AI Agent will automatically detect and use your existing aichat installation:
- No additional configuration needed
- Used as first fallback for LLM tasks

### Ollama Integration
The AI Agent will automatically detect and use your existing Ollama installation:
- Pull models if needed: `ollama pull llama2`
- Used as second fallback for LLM tasks

## Testing Your Setup
```bash
# Test preprocessing
./ai_agent preprocess /path/to/document.pdf

# Test summarization (will use aichat/ollama)
./ai_agent summarize /path/to/document.txt

# Scan and process all documents
./ai_agent scan

# Enable debug logging
DEBUG=true ./ai_agent scan
```

## Logging
The AI Agent uses structured logging with slog:
- Default logs to JSON format in `~/ai_agent.log`
- Enable debug mode with `DEBUG=true` environment variable
- Logs include source file, line numbers, and timestamps

## Troubleshooting

### If aichat is not found
1. Ensure aichat is in your PATH:
   ```bash
   which aichat
   ```
2. If not found, add to PATH or install:
   ```bash
   # Installation example (adjust for your system):
   go install github.com/aichat/aichat@latest
   ```

### If Ollama is not found
1. Ensure Ollama is in your PATH:
   ```bash
   which ollama
   ```
2. If not found, install Ollama:
   ```bash
   # Visit https://ollama.ai for installation instructions
   ```

### If no LLM providers work
1. Check internet connectivity
2. Verify API keys are set correctly
3. Check logs for specific error messages:
   ```bash
   tail -f ~/ai_agent.log
   ```

## Customization Options

### Configure Scan Directories
Modify `config/config.go` to change default scan directories:
```go
ScanDirs: []string{
    filepath.Join(homeDir, "Documents"),
    filepath.Join(homeDir, "Downloads"),
    // Add your preferred directories
},
```

### Adjust Summary Length
Modify `config/config.go` to change default summary length:
```go
MinSummaryLength: 50, // Number of words
```

### Change Default Language for OCR
Modify `config/config.go` to change OCR language:
```go
TesseractLang: "eng", // Change to "spa", "fra", etc.
```

## Example Usage Workflow

1. **Set up directory structure**
   ```bash
   mkdir -p ~/documents/research
   cp my_paper.pdf ~/documents/research/
   ```

2. **Process documents**
   ```bash
   ./ai_agent scan
   ```

3. **View results**
   ```bash
   # Results are saved in ~/processed_data/
   ls ~/processed_data/
   cat ~/processed_data/*.json | jq '.'
   ```

4. **Get summaries**
   ```bash
   ./ai_agent summarize ~/documents/research/my_paper.pdf
   ```

The AI Agent will automatically use your existing aichat and Ollama installations, providing seamless integration with your current setup while offering enhanced capabilities through its autonomous processing pipeline.