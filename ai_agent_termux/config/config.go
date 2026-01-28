package config

import (
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	ScanDirs          []string
	OutputDir         string
	LLaMACppModelPath string
	TesseractLang     string
	PythonScriptsDir  string
	FaissIndexPath    string
	GeminiAPIKey      string
	OllamaCloudAPIKey string
	ClawDBotAPIKey    string
	AIChatAPIKey      string
	MinSummaryLength  int
	LogFilePath       string
	// Database configuration
	TursoURL       string
	TursoAuthToken string
	SQLiteDBPath   string
	UseLocalSQLite bool
	// Android configuration
	AndroidEnabled   bool
	ADBDeviceID      string
	PreferredStorage string   // "internal", "sdcard", "both"
	AutoPullFiles    bool     // Automatically pull discovered files
	AndroidScanDirs  []string // Android directories to scan
	// Notification configuration
	NtfyURL   string
	NtfyTopic string
}

func LoadConfig() *Config {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get user home directory: %v", err)
	}

	cfg := &Config{
		ScanDirs:          []string{filepath.Join(homeDir, "storage/shared/Download"), filepath.Join(homeDir, "storage/shared/Documents"), filepath.Join(homeDir, "test_docs"), filepath.Join(homeDir, "demo_docs")},
		OutputDir:         filepath.Join(homeDir, "processed_data"),
		LLaMACppModelPath: filepath.Join(homeDir, "llama.cpp/models/tinyllama.gguf"),
		TesseractLang:     "eng",
		PythonScriptsDir:  filepath.Join(homeDir, "ai_agent_termux/python_scripts"),
		FaissIndexPath:    filepath.Join(homeDir, "faiss_index.bin"),
		GeminiAPIKey:      os.Getenv("GEMINI_API_KEY"),
		OllamaCloudAPIKey: os.Getenv("OLLAMA_CLOUD_API_KEY"),
		ClawDBotAPIKey:    os.Getenv("CLAWDBOT_API_KEY"),
		AIChatAPIKey:      os.Getenv("AICHAT_API_KEY"),
		MinSummaryLength:  50,
		LogFilePath:       filepath.Join(homeDir, "ai_agent.log"),
		// Database configuration
		TursoURL:       os.Getenv("TURSO_URL"),
		TursoAuthToken: os.Getenv("TURSO_AUTH_TOKEN"),
		SQLiteDBPath:   filepath.Join(homeDir, "ai_agent.db"),
		UseLocalSQLite: true, // Default to local SQLite for offline capability
		// Android configuration
		AndroidEnabled:   true,
		ADBDeviceID:      "",
		PreferredStorage: "both",
		AutoPullFiles:    false,
		AndroidScanDirs:  []string{"/sdcard/", "/storage/emulated/0/"},
		// Notification configuration
		NtfyURL:   os.Getenv("NTFY_URL"),
		NtfyTopic: os.Getenv("NTFY_TOPIC"),
	}

	// Create output directories if they don't exist
	if _, err := os.Stat(cfg.OutputDir); os.IsNotExist(err) {
		os.MkdirAll(cfg.OutputDir, 0755)
	}

	return cfg
}
