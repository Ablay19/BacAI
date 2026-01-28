package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"golang.org/x/exp/slog"
)

// OllamaClient represents a client for Ollama API
type OllamaClient struct {
	BaseURL    string
	HTTPClient *http.Client
	Model      string
}

// OllamaRequest represents a request to Ollama
type OllamaRequest struct {
	Model    string                 `json:"model"`
	Messages []OllamaMessage        `json:"messages,omitempty"`
	Prompt   string                 `json:"prompt,omitempty"`
	Images   []string               `json:"images,omitempty"`
	Stream   bool                   `json:"stream"`
	Options  map[string]interface{} `json:"options,omitempty"`
}

// OllamaMessage represents a message in chat format
type OllamaMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OllamaResponse represents a response from Ollama
type OllamaResponse struct {
	Model              string         `json:"model"`
	CreatedAt          time.Time      `json:"created_at"`
	Response           string         `json:"response,omitempty"`
	Done               bool           `json:"done"`
	TotalDuration      int64          `json:"total_duration,omitempty"`
	LoadDuration       int64          `json:"load_duration,omitempty"`
	PromptEvalCount    int            `json:"prompt_eval_count,omitempty"`
	PromptEvalDuration int64          `json:"prompt_eval_duration,omitempty"`
	EvalCount          int            `json:"eval_count,omitempty"`
	EvalDuration       int64          `json:"eval_duration,omitempty"`
	Message            *OllamaMessage `json:"message,omitempty"`
}

// OllamaModel represents an Ollama model
type OllamaModel struct {
	Name       string    `json:"name"`
	ModifiedAt time.Time `json:"modified_at"`
	Size       int64     `json:"size"`
	Digest     string    `json:"digest"`
	Details    struct {
		Format            string                   `json:"format"`
		Family            string                   `json:"family"`
		Families          []string                 `json:"families"`
		ParameterSize     string                   `json:"parameter_size"`
		QuantizationLevel string                   `json:"quantization_level"`
		Modelfile         []map[string]interface{} `json:"modelfile"`
	} `json:"details"`
}

// AIProvider represents different AI providers
type AIProvider string

const (
	ProviderOllama    AIProvider = "ollama"
	ProviderOpenAI    AIProvider = "openai"
	ProviderGoogle    AIProvider = "google"
	ProviderAnthropic AIProvider = "anthropic"
)

// AIConfig contains configuration for AI processing
type AIConfig struct {
	Provider        AIProvider    `json:"provider"`
	OllamaBaseURL   string        `json:"ollama_base_url"`
	OllamaModel     string        `json:"ollama_model"`
	OpenAIAPIKey    string        `json:"openai_api_key"`
	OpenAIModel     string        `json:"openai_model"`
	GoogleAPIKey    string        `json:"google_api_key"`
	GoogleModel     string        `json:"google_model"`
	AnthropicAPIKey string        `json:"anthropic_api_key"`
	AnthropicModel  string        `json:"anthropic_model"`
	Timeout         time.Duration `json:"timeout"`
	MaxTokens       int           `json:"max_tokens"`
	Temperature     float64       `json:"temperature"`
}

// DefaultAIConfig returns default configuration for AI processing
func DefaultAIConfig() *AIConfig {
	return &AIConfig{
		Provider:      ProviderOllama,
		OllamaBaseURL: "http://localhost:11434",
		OllamaModel:   "llava",
		Timeout:       60 * time.Second,
		MaxTokens:     4096,
		Temperature:   0.7,
	}
}

// NewOllamaClient creates a new Ollama client
func NewOllamaClient(baseURL, model string) *OllamaClient {
	return &OllamaClient{
		BaseURL:    strings.TrimSuffix(baseURL, "/"),
		HTTPClient: &http.Client{Timeout: 60 * time.Second},
		Model:      model,
	}
}

// NewAIClient creates an AI client based on provider
func NewAIClient(config *AIConfig) (*OllamaClient, error) {
	switch config.Provider {
	case ProviderOllama:
		return NewOllamaClient(config.OllamaBaseURL, config.OllamaModel), nil
	default:
		return nil, fmt.Errorf("unsupported AI provider: %s", config.Provider)
	}
}

// Generate generates text from prompt
func (c *OllamaClient) Generate(ctx context.Context, prompt string) (*OllamaResponse, error) {
	request := OllamaRequest{
		Model:  c.Model,
		Prompt: prompt,
		Stream: false,
	}

	return c.doRequest(ctx, "/api/generate", request)
}

// Chat generates response from messages
func (c *OllamaClient) Chat(ctx context.Context, messages []OllamaMessage) (*OllamaResponse, error) {
	request := OllamaRequest{
		Model:    c.Model,
		Messages: messages,
		Stream:   false,
	}

	return c.doRequest(ctx, "/api/chat", request)
}

// AnalyzeImage analyzes an image with text prompt
func (c *OllamaClient) AnalyzeImage(ctx context.Context, imageBase64 string, prompt string) (*OllamaResponse, error) {
	request := OllamaRequest{
		Model:  c.Model,
		Prompt: prompt,
		Images: []string{imageBase64},
		Stream: false,
		Options: map[string]interface{}{
			"temperature": 0.1,
		},
	}

	return c.doRequest(ctx, "/api/generate", request)
}

// AnalyzeImages analyzes multiple images with text prompt
func (c *OllamaClient) AnalyzeImages(ctx context.Context, images []string, prompt string) (*OllamaResponse, error) {
	request := OllamaRequest{
		Model:  c.Model,
		Prompt: prompt,
		Images: images,
		Stream: false,
		Options: map[string]interface{}{
			"temperature": 0.1,
		},
	}

	return c.doRequest(ctx, "/api/generate", request)
}

// ListModels lists available models
func (c *OllamaClient) ListModels(ctx context.Context) ([]OllamaModel, error) {
	url := fmt.Sprintf("%s/api/tags", c.BaseURL)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ollama API returned status %d", resp.StatusCode)
	}

	var response struct {
		Models []OllamaModel `json:"models"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Models, nil
}

// PullModel pulls a model
func (c *OllamaClient) PullModel(ctx context.Context, modelName string, progressCallback func(float64, string)) error {
	url := fmt.Sprintf("%s/api/pull", c.BaseURL)

	request := map[string]string{
		"name": modelName,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Handle streaming response
	decoder := json.NewDecoder(resp.Body)
	for {
		var response struct {
			Status    string `json:"status"`
			Digest    string `json:"digest"`
			Total     int64  `json:"total"`
			Completed int64  `json:"completed"`
		}

		if err := decoder.Decode(&response); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		if progressCallback != nil && response.Total > 0 {
			progress := float64(response.Completed) / float64(response.Total) * 100
			progressCallback(progress, response.Status)
		}

		slog.Debug("Model pull progress", "model", modelName, "status", response.Status)
	}

	return nil
}

// DeleteModel deletes a model
func (c *OllamaClient) DeleteModel(ctx context.Context, modelName string) error {
	url := fmt.Sprintf("%s/api/delete", c.BaseURL)

	request := map[string]string{
		"name": modelName,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "DELETE", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete model, status: %d", resp.StatusCode)
	}

	slog.Info("Model deleted successfully", "model", modelName)
	return nil
}

// HealthCheck checks if Ollama server is healthy
func (c *OllamaClient) HealthCheck(ctx context.Context) error {
	url := fmt.Sprintf("%s/api/tags", c.BaseURL)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ollama server unhealthy, status: %d", resp.StatusCode)
	}

	return nil
}

// doRequest performs an HTTP request to Ollama API
func (c *OllamaClient) doRequest(ctx context.Context, endpoint string, request interface{}) (*OllamaResponse, error) {
	url := fmt.Sprintf("%s%s", c.BaseURL, endpoint)

	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ollama API returned status %d: %s", resp.StatusCode, string(body))
	}

	var response OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	slog.Debug("Ollama request completed",
		"endpoint", endpoint,
		"model", response.Model,
		"done", response.Done)

	return &response, nil
}

// ImageProcessor provides image analysis using Ollama
type ImageProcessor struct {
	ollamaClient *OllamaClient
	visionAPI    *OllamaClient // Could be different model for vision
}

// NewImageProcessor creates a new image processor
func NewImageProcessor(config *AIConfig) (*ImageProcessor, error) {
	ollamaClient, err := NewAIClient(config)
	if err != nil {
		return nil, err
	}

	// Use vision model if specified, otherwise use default model
	visionModel := config.OllamaModel
	if strings.Contains(strings.ToLower(config.OllamaModel), "llava") {
		visionModel = config.OllamaModel
	}

	visionClient := NewOllamaClient(config.OllamaBaseURL, visionModel)

	return &ImageProcessor{
		ollamaClient: ollamaClient,
		visionAPI:    visionClient,
	}, nil
}

// AnalyzeDocument analyzes a document image
func (ip *ImageProcessor) AnalyzeDocument(ctx context.Context, imageBase64 string) (string, error) {
	prompt := `Analyze this document image and provide:
1. A summary of the content
2. Key information extracted
3. Any important dates, names, or numbers
4. Document type (invoice, contract, letter, etc.)

Please be detailed and accurate in your analysis.`

	response, err := ip.visionAPI.AnalyzeImage(ctx, imageBase64, prompt)
	if err != nil {
		return "", fmt.Errorf("failed to analyze document: %v", err)
	}

	return response.Response, nil
}

// ExtractText extracts and enhances text from image
func (ip *ImageProcessor) ExtractText(ctx context.Context, imageBase64 string) (string, error) {
	prompt := `Extract all text from this image accurately. 
- Preserve the original formatting as much as possible
- If there are multiple columns or sections, organize them clearly
- Correct any obvious OCR errors
- Provide the complete text content

Return only the extracted text without commentary.`

	response, err := ip.visionAPI.AnalyzeImage(ctx, imageBase64, prompt)
	if err != nil {
		return "", fmt.Errorf("failed to extract text: %v", err)
	}

	return response.Response, nil
}

// IdentifyObjects identifies objects in the image
func (ip *ImageProcessor) IdentifyObjects(ctx context.Context, imageBase64 string) (string, error) {
	prompt := `Identify all objects in this image. For each object, provide:
1. Object name
2. Approximate location in the image
3. Confidence level
4. Any relevant details or attributes

Be comprehensive and accurate in your analysis.`

	response, err := ip.visionAPI.AnalyzeImage(ctx, imageBase64, prompt)
	if err != nil {
		return "", fmt.Errorf("failed to identify objects: %v", err)
	}

	return response.Response, nil
}

// SolveMath solves math problems from image
func (ip *ImageProcessor) SolveMath(ctx context.Context, imageBase64 string) (string, error) {
	prompt := `Solve the mathematical problem shown in this image. Provide:
1. The complete problem statement
2. Step-by-step solution
3. Final answer clearly stated
4. Alternative methods if applicable

Show all your work and explain each step.`

	response, err := ip.visionAPI.AnalyzeImage(ctx, imageBase64, prompt)
	if err != nil {
		return "", fmt.Errorf("failed to solve math problem: %v", err)
	}

	return response.Response, nil
}

// TranslateText translates text from image
func (ip *ImageProcessor) TranslateText(ctx context.Context, imageBase64 string, targetLanguage string) (string, error) {
	prompt := fmt.Sprintf(`Extract and translate all text from this image to %s. Provide:
1. Original text (transliterated if needed)
2. Translation
3. Any important context or notes

Ensure the translation is accurate and natural-sounding.`, targetLanguage)

	response, err := ip.visionAPI.AnalyzeImage(ctx, imageBase64, prompt)
	if err != nil {
		return "", fmt.Errorf("failed to translate text: %v", err)
	}

	return response.Response, nil
}

// EnhanceOCR enhances OCR results using AI
func (ip *ImageProcessor) EnhanceOCR(ctx context.Context, ocrText string) (string, error) {
	prompt := `Please review and enhance this OCR result. Your task is to:
1. Correct any OCR errors or misrecognitions
2. Improve formatting and structure
3. Fix any obvious spelling or grammar mistakes
4. Preserve the original meaning and content

OCR Result:
"""%s"""

Return the enhanced text only, without commentary.`

	response, err := ip.ollamaClient.Generate(ctx, fmt.Sprintf(prompt, ocrText))
	if err != nil {
		return "", fmt.Errorf("failed to enhance OCR: %v", err)
	}

	return response.Response, nil
}
