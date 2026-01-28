package cmd

import (
	"ai_agent_termux/config"
	"ai_agent_termux/embedding_generator"
	"ai_agent_termux/file_processor"
	"ai_agent_termux/llm_processor"
	"ai_agent_termux/output_manager"
	"ai_agent_termux/preprocessor"

	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
)

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan [dirs...]",
	Short: "Scan and process documents in specified directories",
	Long: `Scan and process documents in specified directories or configured default directories.
	If no directories are provided, scans configured default directories.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.LoadConfig()

		// Use provided directories or default ones
		dirs := args
		if len(dirs) == 0 {
			dirs = cfg.ScanDirs
		}

		slog.Info("Starting file discovery", "directories", dirs)

		// Discover files
		files, err := file_processor.DiscoverFiles(dirs)
		if err != nil {
			slog.Error("Error discovering files", "error", err)
			return
		}

		slog.Info("Files discovered", "count", len(files))

		// Process each file
		for _, file := range files {
			err := processFile(file, cfg)
			if err != nil {
				slog.Warn("Error processing file, continuing with others", "file", file.Path, "error", err)
			}
		}

		slog.Info("File processing completed")
	},
}

func processFile(file file_processor.FileMetadata, cfg *config.Config) error {
	slog.Info("Processing file", "path", file.Path, "type", file.Type)

	// Initialize components
	outputManager := output_manager.NewOutputManager(cfg)

	// Preprocess file
	preprocessedContent, err := preprocessor.PreprocessFile(file, cfg)
	if err != nil {
		slog.Error("Error preprocessing file", "file", file.Path, "error", err)
		return err
	}

	// Create processed file record
	processedFile := outputManager.CreateProcessedFileFromMetadata(file, preprocessedContent.Text)

	// Generate embeddings
	embeddings, err := embedding_generator.GenerateEmbeddingsForFile(preprocessedContent.Text, cfg)
	if err != nil {
		slog.Warn("Error generating embeddings for file", "file", file.Path, "error", err)
	}

	// Process with LLMs
	var summary string
	var categories []string
	var tags []string
	sourceLLM := "local"

	// Try local LLM first
	localLLM := llm_processor.NewLocalLLMProcessor(cfg)
	summary, err = localLLM.Summarize(preprocessedContent.Text)
	if err != nil {
		slog.Warn("Local LLM summarization failed, trying cloud providers", "file", file.Path, "error", err)

		// Fallback to cloud LLMs
		cloudLLM := llm_processor.NewCloudLLMProcessor(cfg)
		summary, err = cloudLLM.SummarizeWithFallback(preprocessedContent.Text)
		if err != nil {
			slog.Error("Cloud LLM summarization failed", "file", file.Path, "error", err)
			summary = "Summary generation failed"
		} else {
			sourceLLM = "cloud"
		}
	}

	// Update processed file with LLM results
	outputManager.UpdateProcessedFile(&processedFile, summary, categories, tags, embeddings, sourceLLM)

	// Save processed file
	err = outputManager.SaveProcessedFile(processedFile)
	if err != nil {
		slog.Error("Error saving processed file", "file", file.Path, "error", err)
		return err
	}

	slog.Info("Successfully processed file", "path", file.Path)
	return nil
}
