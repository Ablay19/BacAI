package llm_processor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"ai_agent_termux/config"
	"golang.org/x/exp/slog"
)

// OllamaModel represents an ollama model
type OllamaModel struct {
	Name string `json:"name"`
	Size string `json:"size"`
}

// BatchProcessRequest represents a batch processing request
type BatchProcessRequest struct {
	Documents []string `json:"documents"`
	Model     string   `json:"model"`
	MaxTokens int      `json:"max_tokens"`
}

// BatchProcessResponse represents a batch processing response
type BatchProcessResponse struct {
	Summaries []string `json:"summaries"`
	Model     string   `json:"model"`
	TotalTime string   `json:"total_time"`
}

type EnhancedCloudLLMProcessor struct {
	config *config.Config
}

func NewEnhancedCloudLLMProcessor(cfg *config.Config) *EnhancedCloudLLMProcessor {
	return &EnhancedCloudLLMProcessor{
		config: cfg,
	}
}

// GetAvailableOllamaModels returns a list of available ollama models
func (e *EnhancedCloudLLMProcessor) GetAvailableOllamaModels() ([]OllamaModel, error) {
	// Check if ollama command is available
	_, err := exec.LookPath("ollama")
	if err != nil {
		return nil, fmt.Errorf("ollama command not found in PATH: %v", err)
	}

	// Get list of models
	cmd := exec.Command("ollama", "list", "--json")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err = cmd.Run()
	if err != nil {
		slog.Warn("Failed to list ollama models", "error", err, "output", out.String())
		return []OllamaModel{}, nil
	}

	var models []OllamaModel
	lines := strings.Split(out.String(), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		// Parse JSON output
		var model OllamaModel
		if err := json.Unmarshal([]byte(line), &model); err != nil {
			// Fallback to text parsing
			parts := strings.Fields(line)
			if len(parts) >= 1 {
				models = append(models, OllamaModel{Name: parts[0], Size: "unknown"})
			}
		} else {
			models = append(models, model)
		}
	}

	return models, nil
}

// GetBestOllamaModel selects the best available model for text processing
func (e *EnhancedCloudLLMProcessor) GetBestOllamaModel() (string, error) {
	models, err := e.GetAvailableOllamaModels()
	if err != nil {
		return "", err
	}

	// Preferred models in order of preference
	preferredModels := []string{
		"llama3.1:8b", "llama3:8b", "llama2:7b", "mistral", "phi3",
		"codellama", "vicuna", "gemma", "qwen",
	}

	// Check for preferred models first
	for _, preferred := range preferredModels {
		for _, available := range models {
			if strings.Contains(available.Name, preferred) {
				slog.Info("Selected ollama model", "model", available.Name)
				return available.Name, nil
			}
		}
	}

	// If no preferred model found, use first available
	if len(models) > 0 {
		slog.Info("Using first available ollama model", "model", models[0].Name)
		return models[0].Name, nil
	}

	// Fallback to llama2
	slog.Info("No models found, will use llama2 (will be pulled)")
	return "llama2", nil
}

// SummarizeWithOllamaEnhanced provides enhanced ollama summarization with model detection
func (e *EnhancedCloudLLMProcessor) SummarizeWithOllamaEnhanced(text string) (string, error) {
	// Check if ollama command is available
	_, err := exec.LookPath("ollama")
	if err != nil {
		return "", fmt.Errorf("ollama command not found in PATH: %v", err)
	}

	// Ensure ollama is running
	if err := e.ensureOllamaRunning(); err != nil {
		slog.Warn("Ollama service issue", "error", err)
	}

	// Get best model
	model := e.config.OllamaCloudAPIKey
	if model == "" {
		model, err = e.GetBestOllamaModel()
		if err != nil {
			return "", fmt.Errorf("failed to get ollama model: %v", err)
		}
	}

	return e.runOllamaSummarization(model, text)
}

// BatchSummarizeWithOllama processes multiple documents efficiently
func (e *EnhancedCloudLLMProcessor) BatchSummarizeWithOllama(documents []string) ([]string, error) {
	if len(documents) < 20 {
		// Use individual processing for small batches
		var summaries []string
		for _, doc := range documents {
			summary, err := e.SummarizeWithOllamaEnhanced(doc)
			if err != nil {
				slog.Error("Failed to summarize document", "error", err)
				summaries = append(summaries, "Error: "+err.Error())
			} else {
				summaries = append(summaries, summary)
			}
		}
		return summaries, nil
	}

	// For large batches (20+ files), use optimized batch processing
	slog.Info("Processing large batch", "count", len(documents), "method", "optimized")
	return e.optimizedBatchSummarize(documents)
}

// optimizedBatchSummarize handles large batches efficiently
func (e *EnhancedCloudLLMProcessor) optimizedBatchSummarize(documents []string) ([]string, error) {
	model, err := e.GetBestOllamaModel()
	if err != nil {
		return nil, err
	}

	// Process in chunks to avoid overwhelming the model
	chunkSize := 10
	var allSummaries []string

	for i := 0; i < len(documents); i += chunkSize {
		end := i + chunkSize
		if end > len(documents) {
			end = len(documents)
		}

		chunk := documents[i:end]
		chunkSummaries, err := e.processChunk(model, chunk)
		if err != nil {
			slog.Error("Failed to process chunk", "start", i, "end", end, "error", err)
			// Add error summaries for this chunk
			for j := i; j < end; j++ {
				allSummaries = append(allSummaries, "Error processing document")
			}
		} else {
			allSummaries = append(allSummaries, chunkSummaries...)
		}

		slog.Info("Processed chunk", "progress", fmt.Sprintf("%d/%d", end, len(documents)))
	}

	return allSummaries, nil
}

// processChunk processes a smaller chunk of documents
func (e *EnhancedCloudLLMProcessor) processChunk(model string, documents []string) ([]string, error) {
	var summaries []string

	for i, doc := range documents {
		summary, err := e.runOllamaSummarization(model, doc)
		if err != nil {
			slog.Error("Failed to summarize document in chunk", "index", i, "error", err)
			summaries = append(summaries, "Error: "+err.Error())
		} else {
			summaries = append(summaries, summary)
		}
	}

	return summaries, nil
}

// runOllamaSummarization performs the actual ollama summarization
func (e *EnhancedCloudLLMProcessor) runOllamaSummarization(model, text string) (string, error) {
	// Truncate text if too long for ollama
	maxLength := 8000 // Conservative limit for ollama
	if len(text) > maxLength {
		text = text[:maxLength] + "..."
	}

	prompt := fmt.Sprintf("Provide a concise summary of this document in %d words or less:\n\n%s",
		e.config.MinSummaryLength, text)

	// Run ollama with optimized parameters
	cmd := exec.Command("ollama", "run", model, prompt)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("ollama command failed: %v, output: %s", err, out.String())
	}

	return strings.TrimSpace(out.String()), nil
}

// ensureOllamaRunning ensures ollama service is running
func (e *EnhancedCloudLLMProcessor) ensureOllamaRunning() error {
	// Check if ollama is responding
	cmd := exec.Command("ollama", "list")
	if err := cmd.Run(); err == nil {
		return nil // ollama is running
	}

	// Try to start ollama service
	slog.Info("Attempting to start ollama service")
	cmd = exec.Command("ollama", "serve")
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start ollama service: %v", err)
	}

	// Give it a moment to start, but don't wait indefinitely
	go cmd.Wait()
	return nil
}

// SummarizeWithAIChatEnhanced provides enhanced aichat support
func (e *EnhancedCloudLLMProcessor) SummarizeWithAIChatEnhanced(text string) (string, error) {
	// Check if aichat command is available
	_, err := exec.LookPath("aichat")
	if err != nil {
		return "", fmt.Errorf("aichat command not found in PATH: %v", err)
	}

	// Get available models from aichat if possible
	model := e.getPreferredAIChatModel()

	// Truncate text if too long
	maxLength := 10000
	if len(text) > maxLength {
		text = text[:maxLength] + "..."
	}

	prompt := fmt.Sprintf("Provide a concise summary of this document in %d words or less:\n\n%s",
		e.config.MinSummaryLength, text)

	// Run aichat with the best parameters
	var cmd *exec.Cmd
	if model != "" {
		cmd = exec.Command("aichat", "--model", model, "--text", prompt)
	} else {
		cmd = exec.Command("aichat", "--text", prompt)
	}

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("aichat command failed: %v, output: %s", err, out.String())
	}

	return strings.TrimSpace(out.String()), nil
}

// getPreferredAIChatModel tries to determine the best model for aichat
func (e *EnhancedCloudLLMProcessor) getPreferredAIChatModel() string {
	// Try to get available models
	cmd := exec.Command("aichat", "list")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return "" // Use default model
	}

	output := out.String()

	// Look for preferred models
	preferredModels := []string{
		"gpt-4", "gpt-3.5-turbo", "claude-3", "gemini-pro",
	}

	for _, preferred := range preferredModels {
		if strings.Contains(strings.ToLower(output), preferred) {
			return preferred
		}
	}

	return "" // Use default model
}

// BatchSummarizeWithAIChat processes multiple documents with aichat
func (e *EnhancedCloudLLMProcessor) BatchSummarizeWithAIChat(documents []string) ([]string, error) {
	var summaries []string

	for i, doc := range documents {
		summary, err := e.SummarizeWithAIChatEnhanced(doc)
		if err != nil {
			slog.Error("Failed to summarize document with aichat", "index", i, "error", err)
			summaries = append(summaries, "Error: "+err.Error())
		} else {
			summaries = append(summaries, summary)
		}

		// Add small delay to avoid rate limiting
		if i > 0 && i%10 == 0 {
			slog.Info("AIChat batch progress", "processed", i, "total", len(documents))
		}
	}

	return summaries, nil
}

// GetAvailableProviders returns all available LLM providers
func (e *EnhancedCloudLLMProcessor) GetAvailableProviders() map[string]bool {
	providers := make(map[string]bool)

	// Check ollama
	if _, err := exec.LookPath("ollama"); err == nil {
		providers["ollama"] = true
	}

	// Check aichat
	if _, err := exec.LookPath("aichat"); err == nil {
		providers["aichat"] = true
	}

	// Check cloud APIs
	if e.config.GeminiAPIKey != "" {
		providers["gemini"] = true
	}

	if e.config.ClawDBotAPIKey != "" {
		providers["claude"] = true
	}

	return providers
}
