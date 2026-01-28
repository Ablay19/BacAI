package output_manager

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"ai_agent_termux/config"
	"ai_agent_termux/file_processor"
)

type ProcessedFile struct {
	FilePath       string            `json:"file_path"`
	FileName       string            `json:"file_name"`
	FileType       string            `json:"file_type"`
	FileSize       int64             `json:"file_size"`
	ModTime        time.Time         `json:"modification_time"`
	ExtractedText  string            `json:"extracted_text"`
	Embeddings     []float64         `json:"embeddings,omitempty"`
	Summary        string            `json:"summary,omitempty"`
	Categories     []string          `json:"categories,omitempty"`
	Tags           []string          `json:"tags,omitempty"`
	Metadata       map[string]string `json:"metadata,omitempty"`
	ProcessingTime time.Time         `json:"processing_time"`
	SourceLLM      string            `json:"source_llm,omitempty"`
}

type OutputManager struct {
	config *config.Config
}

func NewOutputManager(cfg *config.Config) *OutputManager {
	return &OutputManager{
		config: cfg,
	}
}

func (om *OutputManager) SaveProcessedFile(processedFile ProcessedFile) error {
	// Create output directory if it doesn't exist
	if _, err := os.Stat(om.config.OutputDir); os.IsNotExist(err) {
		err = os.MkdirAll(om.config.OutputDir, 0755)
		if err != nil {
			return fmt.Errorf("error creating output directory: %v", err)
		}
	}

	// Generate filename based on original file name and timestamp
	timestamp := time.Now().Format("20060102_150405")
	outputFileName := fmt.Sprintf("%s_%s.json",
		filepath.Base(processedFile.FilePath), timestamp)
	outputPath := filepath.Join(om.config.OutputDir, outputFileName)

	// Marshal to JSON
	jsonData, err := json.MarshalIndent(processedFile, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %v", err)
	}

	// Write to file
	err = os.WriteFile(outputPath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("error writing JSON file: %v", err)
	}

	return nil
}

func (om *OutputManager) CreateProcessedFileFromMetadata(metadata file_processor.FileMetadata, extractedText string) ProcessedFile {
	return ProcessedFile{
		FilePath:       metadata.Path,
		FileName:       metadata.Filename,
		FileType:       metadata.Type,
		FileSize:       metadata.Size,
		ModTime:        metadata.ModTime,
		ExtractedText:  extractedText,
		Metadata:       make(map[string]string),
		ProcessingTime: time.Now(),
	}
}

func (om *OutputManager) UpdateProcessedFile(processedFile *ProcessedFile, summary string, categories []string, tags []string, embeddings []float64, sourceLLM string) {
	processedFile.Summary = summary
	processedFile.Categories = categories
	processedFile.Tags = tags
	processedFile.Embeddings = embeddings
	processedFile.SourceLLM = sourceLLM
}

func (om *OutputManager) LoadProcessedFile(filePath string) (*ProcessedFile, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	var processedFile ProcessedFile
	err = json.Unmarshal(data, &processedFile)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	return &processedFile, nil
}
