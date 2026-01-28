package llm_processor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"

	"ai_agent_termux/config"
	"ai_agent_termux/utils"
	"golang.org/x/exp/slog"
)

type CloudLLMProcessor struct {
	config *config.Config
}

func NewCloudLLMProcessor(cfg *config.Config) *CloudLLMProcessor {
	return &CloudLLMProcessor{
		config: cfg,
	}
}

// GeminiRequest represents the request structure for Gemini API
type GeminiRequest struct {
	Contents []Content `json:"contents"`
}

type Content struct {
	Parts []Part `json:"parts"`
}

type Part struct {
	Text string `json:"text"`
}

// GeminiResponse represents the response structure from Gemini API
type GeminiResponse struct {
	Candidates []Candidate `json:"candidates"`
}

type Candidate struct {
	Content Content `json:"content"`
}

func (cloud *CloudLLMProcessor) SummarizeWithGemini(text string) (string, error) {
	if cloud.config.GeminiAPIKey == "" {
		return "", fmt.Errorf("Gemini API key not configured")
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent?key=%s", cloud.config.GeminiAPIKey)

	prompt := fmt.Sprintf("Summarize the following text in %d words or less:\n\n%s",
		cloud.config.MinSummaryLength, text)

	reqBody := GeminiRequest{
		Contents: []Content{
			{
				Parts: []Part{
					{
						Text: prompt,
					},
				},
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var geminiResp GeminiResponse
	err = json.Unmarshal(body, &geminiResp)
	if err != nil {
		return "", fmt.Errorf("error unmarshaling response: %v", err)
	}

	if len(geminiResp.Candidates) > 0 && len(geminiResp.Candidates[0].Content.Parts) > 0 {
		return geminiResp.Candidates[0].Content.Parts[0].Text, nil
	}

	return "", fmt.Errorf("no summary generated")
}

func (cloud *CloudLLMProcessor) SummarizeWithOllama(text string) (string, error) {
	// Check if ollama command is available
	_, err := exec.LookPath("ollama")
	if err != nil {
		return "", fmt.Errorf("ollama command not found in PATH: %v", err)
	}

	// Check if ollama is running
	_, err = utils.ExecuteCommand("ollama", "list")
	if err != nil {
		// Try to start ollama service
		slog.Info("Attempting to start ollama service")
		_, err = utils.ExecuteCommand("ollama", "serve")
		if err != nil {
			slog.Warn("Failed to start ollama service", "error", err)
		}
	}

	// Use default model if not specified
	model := "llama2"
	if cloud.config.OllamaCloudAPIKey != "" {
		model = cloud.config.OllamaCloudAPIKey // Using API key field to store model name
	}

	prompt := fmt.Sprintf("Summarize the following text in %d words or less:\n\n%s",
		cloud.config.MinSummaryLength, text)

	// Run ollama with the prompt
	cmd := exec.Command("ollama", "run", model, prompt)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("ollama command failed: %v, output: %s", err, out.String())
	}

	return strings.TrimSpace(out.String()), nil
}

func (cloud *CloudLLMProcessor) SummarizeWithAIChat(text string) (string, error) {
	// Check if aichat command is available
	_, err := exec.LookPath("aichat")
	if err != nil {
		return "", fmt.Errorf("aichat command not found in PATH: %v", err)
	}

	prompt := fmt.Sprintf("Summarize the following text in %d words or less:\n\n%s",
		cloud.config.MinSummaryLength, text)

	// Run aichat with the prompt
	cmd := exec.Command("aichat", "--text", prompt)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("aichat command failed: %v, output: %s", err, out.String())
	}

	return strings.TrimSpace(out.String()), nil
}

func (cloud *CloudLLMProcessor) SummarizeWithClaude(text string) (string, error) {
	if cloud.config.ClawDBotAPIKey == "" {
		return "", fmt.Errorf("Claude API key not configured")
	}

	url := "https://api.anthropic.com/v1/messages"

	prompt := fmt.Sprintf("Summarize the following text in %d words or less:\n\n%s",
		cloud.config.MinSummaryLength, text)

	reqBody := map[string]interface{}{
		"model": "claude-3-opus-20240229",
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"max_tokens": 1024,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", cloud.config.ClawDBotAPIKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse Claude response (simplified)
	var claudeResp map[string]interface{}
	err = json.Unmarshal(body, &claudeResp)
	if err != nil {
		return "", fmt.Errorf("error unmarshaling response: %v", err)
	}

	if content, ok := claudeResp["content"].([]interface{}); ok && len(content) > 0 {
		if textBlock, ok := content[0].(map[string]interface{}); ok {
			if textContent, ok := textBlock["text"].(string); ok {
				return textContent, nil
			}
		}
	}

	return "", fmt.Errorf("no summary generated")
}

// Fallback mechanism that tries different cloud providers
func (cloud *CloudLLMProcessor) SummarizeWithFallback(text string) (string, error) {
	slog.Info("Attempting cloud LLM fallback options")

	// Try AIChat first since it's already installed
	slog.Info("Trying AIChat for summarization")
	summary, err := cloud.SummarizeWithAIChat(text)
	if err == nil {
		slog.Info("AIChat summarization successful")
		return summary, nil
	}
	slog.Warn("AIChat summarization failed", "error", err)

	// Try Ollama next
	slog.Info("Trying Ollama for summarization")
	summary, err = cloud.SummarizeWithOllama(text)
	if err == nil {
		slog.Info("Ollama summarization successful")
		return summary, nil
	}
	slog.Warn("Ollama summarization failed", "error", err)

	// Try Claude (ClaudBot API key)
	if cloud.config.ClawDBotAPIKey != "" {
		slog.Info("Trying Claude for summarization")
		summary, err = cloud.SummarizeWithClaude(text)
		if err == nil {
			slog.Info("Claude summarization successful")
			return summary, nil
		}
		slog.Warn("Claude summarization failed", "error", err)
	}

	// Try Gemini as last resort
	if cloud.config.GeminiAPIKey != "" {
		slog.Info("Trying Gemini for summarization")
		summary, err = cloud.SummarizeWithGemini(text)
		if err == nil {
			slog.Info("Gemini summarization successful")
			return summary, nil
		}
		slog.Warn("Gemini summarization failed", "error", err)
	}

	return "", fmt.Errorf("all cloud LLM providers failed")
}
