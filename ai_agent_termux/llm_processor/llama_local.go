package llm_processor

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"ai_agent_termux/config"
)

type LocalLLMProcessor struct {
	config *config.Config
}

func NewLocalLLMProcessor(cfg *config.Config) *LocalLLMProcessor {
	return &LocalLLMProcessor{
		config: cfg,
	}
}

func (llm *LocalLLMProcessor) Summarize(text string) (string, error) {
	// Check if llama.cpp binary exists
	llamaCppPath := "/data/data/com.termux/files/home/llama.cpp/main"
	if _, err := os.Stat(llamaCppPath); os.IsNotExist(err) {
		return "", fmt.Errorf("llama.cpp binary not found at %s", llamaCppPath)
	}

	// Create prompt for summarization
	prompt := fmt.Sprintf("Summarize the following text in %d words or less:\n\n%s",
		llm.config.MinSummaryLength, text)

	// Execute llama.cpp with prompt
	cmd := exec.Command(
		llamaCppPath,
		"--model", llm.config.LLaMACppModelPath,
		"--prompt", prompt,
		"--n_predict", "128",
		"--temp", "0.7",
	)

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error executing llama.cpp: %v", err)
	}

	return strings.TrimSpace(string(output)), nil
}

func (llm *LocalLLMProcessor) Categorize(text string) (string, error) {
	// Check if llama.cpp binary exists
	llamaCppPath := "/data/data/com.termux/files/home/llama.cpp/main"
	if _, err := os.Stat(llamaCppPath); os.IsNotExist(err) {
		return "", fmt.Errorf("llama.cpp binary not found at %s", llamaCppPath)
	}

	// Create prompt for categorization
	prompt := fmt.Sprintf("Categorize the following text into one of these categories: business, personal, technical, creative, educational, other. Respond with only the category name.\n\n%s", text)

	// Execute llama.cpp with prompt
	cmd := exec.Command(
		llamaCppPath,
		"--model", llm.config.LLaMACppModelPath,
		"--prompt", prompt,
		"--n_predict", "10",
		"--temp", "0.5",
	)

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error executing llama.cpp: %v", err)
	}

	return strings.TrimSpace(string(output)), nil
}

func (llm *LocalLLMProcessor) GenerateTags(text string) ([]string, error) {
	// Check if llama.cpp binary exists
	llamaCppPath := "/data/data/com.termux/files/home/llama.cpp/main"
	if _, err := os.Stat(llamaCppPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("llama.cpp binary not found at %s", llamaCppPath)
	}

	// Create prompt for tag generation
	prompt := fmt.Sprintf("Generate 3-5 relevant tags for the following text. Respond with only the tags separated by commas.\n\n%s", text)

	// Execute llama.cpp with prompt
	cmd := exec.Command(
		llamaCppPath,
		"--model", llm.config.LLaMACppModelPath,
		"--prompt", prompt,
		"--n_predict", "32",
		"--temp", "0.7",
	)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error executing llama.cpp: %v", err)
	}

	// Split tags by comma and trim whitespace
	tags := strings.Split(strings.TrimSpace(string(output)), ",")
	for i, tag := range tags {
		tags[i] = strings.TrimSpace(tag)
	}

	return tags, nil
}
