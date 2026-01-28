package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"ai_agent_termux/config"

	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config [command]",
	Short: "Manage AI Agent configuration",
	Long: `View, edit, or validate AI Agent configuration.
	Configuration includes scan directories, LLM settings, and API keys.`,
}

func init() {
	// Add subcommands
	configCmd.AddCommand(configViewCmd)
	configCmd.AddCommand(configEditCmd)
	configCmd.AddCommand(configValidateCmd)
}

// configViewCmd represents the config view command
var configViewCmd = &cobra.Command{
	Use:   "view",
	Short: "View current configuration",
	Long:  `Display the current AI Agent configuration settings.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.LoadConfig()

		fmt.Println("Current AI Agent Configuration:")
		fmt.Println("===============================")
		fmt.Printf("Scan Directories: %v\n", cfg.ScanDirs)
		fmt.Printf("Output Directory: %s\n", cfg.OutputDir)
		fmt.Printf("LLaMA.cpp Model Path: %s\n", cfg.LLaMACppModelPath)
		fmt.Printf("Tesseract Language: %s\n", cfg.TesseractLang)
		fmt.Printf("Python Scripts Directory: %s\n", cfg.PythonScriptsDir)
		fmt.Printf("FAISS Index Path: %s\n", cfg.FaissIndexPath)
		fmt.Printf("Minimum Summary Length: %d\n", cfg.MinSummaryLength)
		fmt.Printf("Log File Path: %s\n", cfg.LogFilePath)

		// Show database configuration
		fmt.Println("\nDatabase Configuration:")
		if cfg.UseLocalSQLite {
			fmt.Printf("  Local SQLite: Enabled (%s)\n", cfg.SQLiteDBPath)
		} else {
			fmt.Println("  Local SQLite: Disabled")
		}

		if cfg.TursoURL != "" {
			fmt.Printf("  Turso URL: Configured\n")
		} else {
			fmt.Println("  Turso URL: Not configured")
		}

		if cfg.TursoAuthToken != "" {
			fmt.Printf("  Turso Auth: Configured\n")
		} else {
			fmt.Println("  Turso Auth: Not configured")
		}

		// Show Android configuration
		fmt.Println("\nAndroid Configuration:")
		if cfg.AndroidEnabled {
			fmt.Printf("  Android Integration: Enabled\n")
		} else {
			fmt.Println("  Android Integration: Disabled")
		}

		fmt.Printf("  Preferred Storage: %s\n", cfg.PreferredStorage)
		fmt.Printf("  Auto Pull Files: %t\n", cfg.AutoPullFiles)
		fmt.Printf("  Android Scan Dirs: %v\n", cfg.AndroidScanDirs)

		// Show API key status (not actual keys for security)
		fmt.Println("\nAPI Keys:")
		if cfg.GeminiAPIKey != "" {
			fmt.Println("  Gemini: Configured")
		} else {
			fmt.Println("  Gemini: Not configured")
		}

		if cfg.OllamaCloudAPIKey != "" {
			fmt.Println("  Ollama Model: Configured")
		} else {
			fmt.Println("  Ollama Model: Using default (llama2)")
		}

		if cfg.ClawDBotAPIKey != "" {
			fmt.Println("  Claude: Configured")
		} else {
			fmt.Println("  Claude: Not configured")
		}

		if cfg.AIChatAPIKey != "" {
			fmt.Println("  AIChat: Configured")
		} else {
			fmt.Println("  AIChat: Using system installation")
		}
	},
}

// configEditCmd represents the config edit command
var configEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit configuration file",
	Long:  `Open the configuration file in your default editor.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Determine config file path
		configPath := filepath.Join(os.Getenv("HOME"), ".ai_agent.yaml")

		// Get config file flag from root command
		if cmd.Flags().Changed("config") {
			if configFilePath, err := cmd.Flags().GetString("config"); err == nil && configFilePath != "" {
				configPath = configFilePath
			}
		}

		fmt.Printf("Editing configuration file: %s\n", configPath)
		fmt.Println("Configuration editing is not yet implemented.")
		fmt.Println("Please manually edit the config.go file for now.")

		// TODO: Implement actual config file editing
		// This would typically involve:
		// 1. Creating a YAML config file
		// 2. Opening it in an editor
		// 3. Validating changes
		// 4. Saving updates
	},
}

// configValidateCmd represents the config validate command
var configValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate configuration",
	Long:  `Check if the current configuration is valid and all dependencies are available.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.LoadConfig()

		fmt.Println("Validating AI Agent Configuration...")
		fmt.Println("====================================")

		// Check scan directories
		allValid := true
		for _, dir := range cfg.ScanDirs {
			if _, err := os.Stat(dir); os.IsNotExist(err) {
				fmt.Printf("❌ Scan directory does not exist: %s\n", dir)
				allValid = false
			} else {
				fmt.Printf("✅ Scan directory exists: %s\n", dir)
			}
		}

		// Check output directory
		if _, err := os.Stat(cfg.OutputDir); os.IsNotExist(err) {
			fmt.Printf("⚠️  Output directory does not exist (will be created): %s\n", cfg.OutputDir)
		} else {
			fmt.Printf("✅ Output directory exists: %s\n", cfg.OutputDir)
		}

		// Check LLM dependencies
		fmt.Println("\nChecking LLM dependencies...")

		// Check local LLM
		if _, err := os.Stat(cfg.LLaMACppModelPath); os.IsNotExist(err) {
			fmt.Printf("⚠️  LLaMA.cpp model not found: %s\n", cfg.LLaMACppModelPath)
			fmt.Println("   Local LLM processing will be skipped.")
		} else {
			fmt.Printf("✅ LLaMA.cpp model found: %s\n", cfg.LLaMACppModelPath)
		}

		// Check aichat
		aichatFound := false
		if _, err := os.Stat("/usr/local/bin/aichat"); err == nil {
			aichatFound = true
		} else if _, err := os.Stat("/usr/bin/aichat"); err == nil {
			aichatFound = true
		}

		if aichatFound {
			fmt.Printf("✅ aichat found in system PATH\n")
		} else {
			fmt.Printf("⚠️  aichat not found in system PATH\n")
			fmt.Println("   Will use cloud APIs as fallback.")
		}

		// Check ollama
		ollamaFound := false
		if _, err := os.Stat("/usr/local/bin/ollama"); err == nil {
			ollamaFound = true
		} else if _, err := os.Stat("/usr/bin/ollama"); err == nil {
			ollamaFound = true
		}

		if ollamaFound {
			fmt.Printf("✅ Ollama found in system PATH\n")
		} else {
			fmt.Printf("⚠️  Ollama not found in system PATH\n")
			fmt.Println("   Will use other providers as fallback.")
		}

		// Validate database configuration
		fmt.Println("\nChecking database configuration...")

		if cfg.UseLocalSQLite {
			// Check if SQLite DB directory exists and is writable
			dbDir := filepath.Dir(cfg.SQLiteDBPath)
			if _, err := os.Stat(dbDir); os.IsNotExist(err) {
				fmt.Printf("⚠️  SQLite DB directory does not exist (will be created): %s\n", dbDir)
			} else {
				fmt.Printf("✅ SQLite DB directory accessible: %s\n", dbDir)
			}
			fmt.Printf("✅ Local SQLite database enabled: %s\n", cfg.SQLiteDBPath)
		} else {
			fmt.Println("ℹ️  Local SQLite database disabled")
		}

		if cfg.TursoURL != "" && cfg.TursoAuthToken != "" {
			fmt.Printf("✅ Turso database configuration complete\n")
		} else if cfg.TursoURL != "" || cfg.TursoAuthToken != "" {
			fmt.Printf("⚠️  Turso database partially configured (need both URL and auth token)\n")
		} else {
			fmt.Println("ℹ️  Turso database not configured")
		}

		// Validate Android configuration
		fmt.Println("\nChecking Android configuration...")

		if cfg.AndroidEnabled {
			fmt.Printf("✅ Android integration enabled\n")

			// Check for required tools
			_, err := exec.LookPath("adb")
			if err == nil {
				fmt.Printf("✅ ADB available\n")
			} else {
				fmt.Printf("⚠️  ADB not found - limited Android functionality\n")
			}

			_, err = exec.LookPath("termux-notification")
			if err == nil {
				fmt.Printf("✅ Termux API available\n")
			} else {
				fmt.Printf("⚠️  Termux API not found - notifications disabled\n")
			}
		} else {
			fmt.Println("ℹ️  Android integration disabled")
		}

		// Validate API keys
		fmt.Println("\nChecking API key configurations...")
		if cfg.GeminiAPIKey != "" {
			fmt.Printf("✅ Gemini API key configured\n")
		} else {
			fmt.Printf("ℹ️  Gemini API key not configured\n")
		}

		if cfg.ClawDBotAPIKey != "" {
			fmt.Printf("✅ Claude API key configured\n")
		} else {
			fmt.Printf("ℹ️  Claude API key not configured\n")
		}

		if allValid {
			fmt.Println("\n✅ Configuration validation completed successfully!")
		} else {
			fmt.Println("\n⚠️  Configuration validation completed with warnings.")
		}
	},
}
