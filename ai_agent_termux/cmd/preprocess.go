package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"ai_agent_termux/config"
	"ai_agent_termux/file_processor"
	"ai_agent_termux/preprocessor"

	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
)

// preprocessCmd represents the preprocess command
var preprocessCmd = &cobra.Command{
	Use:   "preprocess [file]",
	Short: "Preprocess a specific file",
	Long: `Preprocess a specific file and output extracted content.
	Supports text files, PDFs, and images with OCR capabilities.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]

		// Check if file exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			slog.Error("File does not exist", "file", filePath)
			fmt.Printf("Error: File does not exist: %s\n", filePath)
			return
		}

		// Create file metadata
		ext := filepath.Ext(filePath)
		fileType := "other"
		switch strings.ToLower(ext) {
		case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp":
			fileType = "image"
		case ".pdf":
			fileType = "pdf"
		case ".txt", ".md", ".log", ".csv", ".json", ".xml", ".html", ".htm":
			fileType = "text"
		}

		fileMeta := file_processor.FileMetadata{
			Path:     filePath,
			Filename: filepath.Base(filePath),
			Ext:      ext,
			Type:     fileType,
		}

		// Load config
		cfg := config.LoadConfig()

		// Preprocess the file
		preprocessedContent, err := preprocessor.PreprocessFile(fileMeta, cfg)
		if err != nil {
			slog.Error("Error preprocessing file", "file", filePath, "error", err)
			fmt.Printf("Error preprocessing file: %v\n", err)
			return
		}

		fmt.Printf("Preprocessed content:\n%s\n", preprocessedContent.Text)
	},
}
