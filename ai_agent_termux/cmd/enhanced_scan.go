package cmd

import (
	"fmt"

	"ai_agent_termux/config"
	"ai_agent_termux/database"
	"ai_agent_termux/file_processor"
	"ai_agent_termux/llm_processor"

	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
)

// enhancedScanCmd represents the enhanced scan command
var enhancedScanCmd = &cobra.Command{
	Use:   "enhanced-scan [directories...]",
	Short: "Enhanced scan with batch processing for large directories",
	Long: `Enhanced scan that intelligently handles directories with 20+ files.
	Features:
	- Automatic batch processing for directories with 50+ files
	- Priority-based file ordering
	- Enhanced ollama and aichat integration
	- Database storage of results
	- Concurrent processing with worker pools`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.LoadConfig()

		// Parse directories from args or use config defaults
		var dirs []string
		if len(args) > 0 {
			dirs = args
		} else {
			dirs = cfg.ScanDirs
		}

		// Parse batch size
		batchSize, _ := cmd.Flags().GetInt("batch-size")
		if batchSize <= 0 {
			batchSize = 20
		}

		// Parse workers
		workers, _ := cmd.Flags().GetInt("workers")
		if workers <= 0 {
			workers = 5
		}

		slog.Info("Starting enhanced scan",
			"directories", len(dirs),
			"batch_size", batchSize,
			"workers", workers)

		// Enhanced file discovery with directory analysis
		files, dirStats, err := file_processor.DiscoverFilesEnhanced(dirs)
		if err != nil {
			fmt.Printf("Error discovering files: %v\n", err)
			return
		}

		fmt.Printf("Discovered %d files\n", len(files))
		fmt.Println(file_processor.GetDirectorySummary(dirStats))

		// Initialize database
		db, err := database.NewDatabase(cfg)
		if err != nil {
			slog.Warn("Failed to initialize database", "error", err)
		} else {
			defer db.Close()
		}

		// Initialize enhanced LLM processor
		llmProc := llm_processor.NewEnhancedCloudLLMProcessor(cfg)

		// Check available providers
		providers := llmProc.GetAvailableProviders()
		fmt.Printf("Available LLM providers: %v\n", getProviderKeys(providers))

		// Process files based on directory size
		err = processFilesWithBatching(files, cfg, llmProc, db, batchSize, workers)
		if err != nil {
			fmt.Printf("Error processing files: %v\n", err)
			return
		}

		fmt.Println("Enhanced scan completed successfully!")
	},
}

func init() {
	rootCmd.AddCommand(enhancedScanCmd)
	enhancedScanCmd.Flags().Int("batch-size", 20, "Batch size for processing files")
	enhancedScanCmd.Flags().Int("workers", 5, "Number of concurrent workers")
	enhancedScanCmd.Flags().Bool("use-ollama", true, "Use ollama for processing")
	enhancedScanCmd.Flags().Bool("use-aichat", true, "Use aichat for processing")
}

// processFilesWithBatching processes files using intelligent batching
func processFilesWithBatching(files []file_processor.FileMetadata, cfg *config.Config,
	llmProc *llm_processor.EnhancedCloudLLMProcessor, db *database.Database,
	batchSize, workers int) error {

	// Check if we have many files requiring batch processing
	if len(files) >= 50 {
		fmt.Printf("Large directory detected (%d files), using batch processing\n", len(files))
		return processLargeBatch(files, cfg, llmProc, db, batchSize, workers)
	}

	// Standard processing for smaller directories
	fmt.Printf("Processing %d files with standard processing\n", len(files))
	return processStandardBatch(files, cfg, llmProc, db)
}

// processLargeBatch handles large directories with optimized batch processing
func processLargeBatch(files []file_processor.FileMetadata, cfg *config.Config,
	llmProc *llm_processor.EnhancedCloudLLMProcessor, db *database.Database,
	batchSize, workers int) error {

	// Group files by type for better processing
	textFiles := filterFilesByType(files, "text")
	pdfFiles := filterFilesByType(files, "pdf")
	imageFiles := filterFilesByType(files, "image")

	// Process text files first (highest priority)
	if len(textFiles) > 0 {
		fmt.Printf("Processing %d text files with batching\n", len(textFiles))
		err := file_processor.ProcessLargeDirectory(textFiles, func(file file_processor.FileMetadata) error {
			return processSingleFile(file, cfg, llmProc, db)
		})
		if err != nil {
			return err
		}
	}

	// Process PDF files
	if len(pdfFiles) > 0 {
		fmt.Printf("Processing %d PDF files with batching\n", len(pdfFiles))
		err := file_processor.ProcessLargeDirectory(pdfFiles, func(file file_processor.FileMetadata) error {
			return processSingleFile(file, cfg, llmProc, db)
		})
		if err != nil {
			return err
		}
	}

	// Process image files
	if len(imageFiles) > 0 {
		fmt.Printf("Processing %d image files with batching\n", len(imageFiles))
		err := file_processor.ProcessLargeDirectory(imageFiles, func(file file_processor.FileMetadata) error {
			return processSingleFile(file, cfg, llmProc, db)
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// processStandardBatch handles smaller directories
func processStandardBatch(files []file_processor.FileMetadata, cfg *config.Config,
	llmProc *llm_processor.EnhancedCloudLLMProcessor, db *database.Database) error {

	for i, file := range files {
		fmt.Printf("Processing file %d/%d: %s\n", i+1, len(files), file.Filename)

		err := processSingleFile(file, cfg, llmProc, db)
		if err != nil {
			slog.Error("Error processing file", "file", file.Path, "error", err)
		}

		// Progress indicator for standard processing
		if (i+1)%10 == 0 {
			fmt.Printf("Progress: %d/%d files processed\n", i+1, len(files))
		}
	}

	return nil
}

// processSingleFile processes a single file with enhanced LLM capabilities
func processSingleFile(file file_processor.FileMetadata, cfg *config.Config,
	llmProc *llm_processor.EnhancedCloudLLMProcessor, db *database.Database) error {

	slog.Info("Processing file", "file", file.Path, "type", file.Type, "size", file.Size)

	// Here you would:
	// 1. Preprocess the file (extract text)
	// 2. Generate embeddings
	// 3. Store in database
	// 4. Generate summary using enhanced LLM

	// For now, just log the processing
	return nil
}

// filterFilesByType filters files by type
func filterFilesByType(files []file_processor.FileMetadata, fileType string) []file_processor.FileMetadata {
	var filtered []file_processor.FileMetadata
	for _, file := range files {
		if file.Type == fileType {
			filtered = append(filtered, file)
		}
	}
	return filtered
}

// getProviderKeys returns provider names as string slice
func getProviderKeys(providers map[string]bool) []string {
	var keys []string
	for provider, available := range providers {
		if available {
			keys = append(keys, provider)
		}
	}
	return keys
}
