//go:build standalone_agent

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"ai_agent_termux/ai"
	"ai_agent_termux/automation"
	"ai_agent_termux/config"
	"ai_agent_termux/goroutines"
	"ai_agent_termux/interactive"
	"ai_agent_termux/notifications"
	"ai_agent_termux/realtime"
)

// Version and build information
var (
	Version   = "dev"
	GitCommit = "unknown"
	BuildTime = "unknown"
)

// Command line flags
var (
	configFile = flag.String("config", "", "Configuration file path")
	mode       = flag.String("mode", "serve", "Operation mode: serve, worker, batch, interactive")
	inputFiles = flag.String("files", "", "Input files for batch processing (comma-separated)")
	operation  = flag.String("operation", "extract_text", "Google Lens operation")
	watchDir   = flag.String("watch", "", "Directory to watch for new images")
	serverPort = flag.String("port", "8080", "Server port")
	verbose    = flag.Bool("verbose", false, "Enable verbose logging")
	version    = flag.Bool("version", false, "Show version information")
)

// App represents the main application
type App struct {
	config       *config.Config
	googleLens   *automation.GoogleLensProcessor
	aiProcessor  *ai.ImageProcessor
	notifySystem *notifications.NotifySystem
	taskManager  *goroutines.TaskManager
	fileWatcher  *realtime.FileWatcher
	cli          *interactive.InteractiveCLI
}

func main() {
	flag.Parse()

	if *version {
		printVersion()
		return
	}

	// Initialize application
	app, err := initializeApp()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}
	defer cleanup(app)

	// Setup signal handling
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Run based on mode
	switch *mode {
	case "serve":
		runServer(ctx, app)
	case "worker":
		runWorker(ctx, app)
	case "batch":
		runBatch(ctx, app)
	case "interactive":
		runInteractive(ctx, app)
	default:
		log.Fatalf("Unknown mode: %s", *mode)
	}

	// Wait for shutdown signal
	<-sigChan
	log.Println("Received shutdown signal")
}

// initializeApp initializes the application components
func initializeApp() (*App, error) {
	// Load configuration
	cfg, err := loadConfiguration()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %v", err)
	}

	// Initialize Google Lens processor
	googleLens := automation.NewGoogleLensProcessor(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Google Lens processor: %v", err)
	}

	// Initialize AI processor (Ollama)
	aiConfig := ai.DefaultAIConfig()
	aiConfig.OllamaBaseURL = cfg.OllamaBaseURL
	aiConfig.OllamaModel = cfg.OllamaModel

	aiProcessor, err := ai.NewImageProcessor(aiConfig)
	if err != nil {
		log.Printf("Warning: Failed to initialize AI processor: %v", err)
		aiProcessor = nil
	}

	// Initialize notification system
	notifyConfig := notifications.DefaultNotifyConfig()
	notifyConfig.NtfyTopic = cfg.NtfyTopic
	notifyConfig.NtfyServer = cfg.NtfyServer

	notifySystem := notifications.NewNotifySystem(notifyConfig)

	// Initialize task manager
	taskConfig := goroutines.DefaultManagerConfig()
	taskConfig.MaxConcurrent = cfg.MaxConcurrentTasks

	taskManager := goroutines.NewTaskManager(taskConfig)

	// Initialize file watcher
	var fileWatcher *realtime.FileWatcher
	if aiProcessor != nil {
		watcherConfig := realtime.DefaultWatcherConfig()
		watcherConfig.Operation = *operation
		fileWatcher, err = realtime.NewFileWatcher(googleLens, notifySystem, watcherConfig)
		if err != nil {
			log.Printf("Warning: Failed to initialize file watcher: %v", err)
			fileWatcher = nil
		}
	}

	// Initialize interactive CLI
	cliConfig := interactive.DefaultCLIConfig()
	cliConfig.GoogleLens = googleLens
	cliConfig.TaskManager = taskManager

	cli := interactive.NewInteractiveCLI(cliConfig)

	app := &App{
		config:       cfg,
		googleLens:   googleLens,
		aiProcessor:  aiProcessor,
		notifySystem: notifySystem,
		taskManager:  taskManager,
		fileWatcher:  fileWatcher,
		cli:          cli,
	}

	log.Printf("Application initialized successfully")
	log.Printf("Google Lens API: %s", func() string {
		if googleLens != nil {
			return "Available"
		}
		return "Unavailable"
	}())
	log.Printf("AI Processor: %s", func() string {
		if aiProcessor != nil {
			return "Available"
		}
		return "Unavailable"
	}())
	log.Printf("Notifications: %s", func() string {
		if notifySystem != nil {
			return "Enabled"
		}
		return "Disabled"
	}())

	return app, nil
}

// loadConfiguration loads configuration from file and environment
func loadConfiguration() (*config.Config, error) {
	// Load configuration file if specified
	var cfg *config.Config
	if *configFile != "" {
		cfg, err := config.LoadFromFile(*configFile)
		if err != nil {
			return nil, fmt.Errorf("failed to load config file: %v", err)
		}
	} else {
		cfg = config.DefaultConfig()
	}

	// Override with environment variables
	if port := os.Getenv("PORT"); port != "" {
		cfg.ServerPort = port
	}
	if apiKey := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"); apiKey != "" {
		cfg.GoogleCredentialsPath = apiKey
	}
	if serpAPIKey := os.Getenv("SERPAPI_API_KEY"); serpAPIKey != "" {
		cfg.SerpAPIKey = serpAPIKey
	}
	if ntfyTopic := os.Getenv("NTFY_TOPIC"); ntfyTopic != "" {
		cfg.NtfyTopic = ntfyTopic
	}
	if ollamaURL := os.Getenv("OLLAMA_BASE_URL"); ollamaURL != "" {
		cfg.OllamaBaseURL = ollamaURL
	}
	if ollamaModel := os.Getenv("OLLAMA_MODEL"); ollamaModel != "" {
		cfg.OllamaModel = ollamaModel
	}

	return cfg, nil
}

// runServer starts the application in server mode
func runServer(ctx context.Context, app *App) {
	log.Printf("Starting server mode on port %s", *serverPort)

	// TODO: Implement HTTP server
	// For now, just start file watching if directory provided
	if *watchDir != "" {
		if app.fileWatcher != nil {
			err := app.fileWatcher.AddPath(*watchDir)
			if err != nil {
				log.Printf("Failed to add watch directory: %v", err)
			} else {
				log.Printf("Watching directory: %s", *watchDir)
			}
		}
	}

	// Keep running until context is cancelled
	<-ctx.Done()
	log.Println("Server mode stopped")
}

// runWorker starts the application in worker mode
func runWorker(ctx context.Context, app *App) {
	log.Println("Starting worker mode")

	// Start background processing
	// TODO: Implement worker logic

	<-ctx.Done()
	log.Println("Worker mode stopped")
}

// runBatch starts the application in batch processing mode
func runBatch(ctx context.Context, app *App) {
	if *inputFiles == "" {
		log.Fatal("No input files specified for batch processing")
	}

	files := strings.Split(*inputFiles, ",")
	log.Printf("Starting batch processing for %d files", len(files))

	// Process files using Google Lens
	progressCallback := func(taskID string, progress float64, message string) {
		log.Printf("Progress: %.1f%% - %s", progress*100, message)
	}

	results, err := app.googleLens.BatchProcessImagesWithProgress(files, *operation, progressCallback)
	if err != nil {
		log.Printf("Batch processing failed: %v", err)
		return
	}

	log.Printf("Batch processing completed for %d files", len(results))
	for i, result := range results {
		log.Printf("File %d: %s", i+1, filepath.Base(result.ImagePath))
		log.Printf("  Result: %s", result.ResultText)
	}
}

// runInteractive starts the application in interactive mode
func runInteractive(ctx context.Context, app *App) {
	log.Println("Starting interactive mode")

	err := app.cli.Start()
	if err != nil {
		log.Printf("Interactive mode error: %v", err)
	}
}

// cleanup performs application cleanup
func cleanup(app *App) {
	log.Println("Performing application cleanup")

	if app.fileWatcher != nil {
		app.fileWatcher.Close()
	}

	if app.googleLens != nil {
		app.googleLens.Close()
	}

	if app.taskManager != nil {
		app.taskManager.Shutdown(10 * time.Second)
	}

	log.Println("Cleanup completed")
}

// printVersion prints version information
func printVersion() {
	fmt.Printf("Google Lens AI Agent\n")
	fmt.Printf("Version: %s\n", Version)
	fmt.Printf("Git Commit: %s\n", GitCommit)
	fmt.Printf("Build Time: %s\n", BuildTime)
}
