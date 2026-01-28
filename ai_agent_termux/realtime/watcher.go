package realtime

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"ai_agent_termux/automation"
	"ai_agent_termux/notifications"
	"github.com/fsnotify/fsnotify"
	"golang.org/x/exp/slog"
)

// FileWatcher monitors directories for new image files and triggers Google Lens processing
type FileWatcher struct {
	watcher          *fsnotify.Watcher
	googleLens       *automation.GoogleLensProcessor
	notifySystem     *notifications.NotifySystem
	watchedPaths     map[string]bool
	pathsMutex       sync.RWMutex
	debounceMap      map[string]time.Time
	debounceMutex    sync.RWMutex
	debounceInterval time.Duration
	ctx              context.Context
	cancel           context.CancelFunc
	processingQueue  chan FileEvent
	wg               sync.WaitGroup
	maxWorkers       int
	operation        string
	autoPreprocess   bool
	notifyOnComplete bool
	enabled          bool
}

// FileEvent represents a file event
type FileEvent struct {
	Path      string    `json:"path"`
	Operation string    `json:"operation"`
	Timestamp time.Time `json:"timestamp"`
	Size      int64     `json:"size"`
}

// WatcherConfig contains configuration for the file watcher
type WatcherConfig struct {
	Paths            []string      `json:"paths"`
	Operation        string        `json:"operation"`          // Google Lens operation to perform
	DebounceInterval time.Duration `json:"debounce_interval"`  // Debounce interval for file events
	MaxWorkers       int           `json:"max_workers"`        // Maximum concurrent processing workers
	AutoPreprocess   bool          `json:"auto_preprocess"`    // Auto-preprocess images
	NotifyOnComplete bool          `json:"notify_on_complete"` // Send notification on processing complete
	FileExtensions   []string      `json:"file_extensions"`    // Supported file extensions
	IgnorePatterns   []string      `json:"ignore_patterns"`    // File patterns to ignore
}

// DefaultWatcherConfig returns default configuration for the file watcher
func DefaultWatcherConfig() *WatcherConfig {
	return &WatcherConfig{
		Operation:        "extract_text",
		DebounceInterval: 2 * time.Second,
		MaxWorkers:       3,
		AutoPreprocess:   true,
		NotifyOnComplete: true,
		FileExtensions:   []string{".jpg", ".jpeg", ".png", ".bmp", ".webp", ".tiff"},
		IgnorePatterns:   []string{".tmp", ".temp", "~", "Thumbs.db", ".DS_Store"},
	}
}

// NewFileWatcher creates a new file watcher instance
func NewFileWatcher(googleLens *automation.GoogleLensProcessor, notifySystem *notifications.NotifySystem, config *WatcherConfig) (*FileWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create file watcher: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	fw := &FileWatcher{
		watcher:          watcher,
		googleLens:       googleLens,
		notifySystem:     notifySystem,
		watchedPaths:     make(map[string]bool),
		debounceMap:      make(map[string]time.Time),
		debounceInterval: config.DebounceInterval,
		ctx:              ctx,
		cancel:           cancel,
		processingQueue:  make(chan FileEvent, 100),
		maxWorkers:       config.MaxWorkers,
		operation:        config.Operation,
		autoPreprocess:   config.AutoPreprocess,
		notifyOnComplete: config.NotifyOnComplete,
		enabled:          true,
	}

	// Start the worker pool
	fw.startWorkers()

	// Start the event processor
	fw.startEventProcessor(config)

	slog.Info("File watcher initialized", "config", config)
	return fw, nil
}

// AddPath adds a directory path to watch for new files
func (fw *FileWatcher) AddPath(path string) error {
	fw.pathsMutex.Lock()
	defer fw.pathsMutex.Unlock()

	if fw.watchedPaths[path] {
		slog.Debug("Path already being watched", "path", path)
		return nil
	}

	// Add the path to the watcher
	err := fw.watcher.Add(path)
	if err != nil {
		return fmt.Errorf("failed to add watch path %s: %v", path, err)
	}

	fw.watchedPaths[path] = true
	slog.Info("Added watch path", "path", path)

	// Process existing files in the directory
	go fw.processExistingFiles(path)

	return nil
}

// RemovePath removes a directory path from watching
func (fw *FileWatcher) RemovePath(path string) error {
	fw.pathsMutex.Lock()
	defer fw.pathsMutex.Unlock()

	if !fw.watchedPaths[path] {
		slog.Debug("Path not being watched", "path", path)
		return nil
	}

	err := fw.watcher.Remove(path)
	if err != nil {
		return fmt.Errorf("failed to remove watch path %s: %v", path, err)
	}

	delete(fw.watchedPaths, path)
	slog.Info("Removed watch path", "path", path)
	return nil
}

// GetWatchedPaths returns all currently watched paths
func (fw *FileWatcher) GetWatchedPaths() []string {
	fw.pathsMutex.RLock()
	defer fw.pathsMutex.RUnlock()

	paths := make([]string, 0, len(fw.watchedPaths))
	for path := range fw.watchedPaths {
		paths = append(paths, path)
	}
	return paths
}

// SetEnabled enables or disables the file watcher
func (fw *FileWatcher) SetEnabled(enabled bool) {
	fw.enabled = enabled
	slog.Info("File watcher enabled status changed", "enabled", enabled)
}

// IsEnabled returns whether the file watcher is enabled
func (fw *FileWatcher) IsEnabled() bool {
	return fw.enabled
}

// startEventProcessor starts the main event processing loop
func (fw *FileWatcher) startEventProcessor(config *WatcherConfig) {
	fw.wg.Add(1)
	go func() {
		defer fw.wg.Done()
		slog.Info("File watcher event processor started")

		for {
			select {
			case <-fw.ctx.Done():
				slog.Info("File watcher event processor stopping")
				return

			case event, ok := <-fw.watcher.Events:
				if !ok {
					return
				}

				if !fw.enabled {
					continue
				}

				fw.handleFileSystemEvent(event, config)

			case err, ok := <-fw.watcher.Errors:
				if !ok {
					return
				}
				slog.Error("File watcher error", "error", err)
			}
		}
	}()
}

// handleFileSystemEvent handles individual file system events
func (fw *FileWatcher) handleFileSystemEvent(event fsnotify.Event, config *WatcherConfig) {
	// Only handle file creation events
	if event.Op&fsnotify.Create != fsnotify.Create {
		return
	}

	// Check if it's an image file
	if !fw.isImageFile(event.Name, config.FileExtensions) {
		return
	}

	// Check if it should be ignored
	if fw.shouldIgnoreFile(event.Name, config.IgnorePatterns) {
		return
	}

	// Debounce the event to handle rapid file creation
	if fw.isDebounced(event.Name) {
		return
	}

	slog.Debug("New image file detected", "path", event.Name)

	// Get file info
	fileInfo, err := os.Stat(event.Name)
	if err != nil {
		slog.Warn("Failed to get file info", "path", event.Name, "error", err)
		return
	}

	// Queue the file for processing
	fileEvent := FileEvent{
		Path:      event.Name,
		Operation: fw.operation,
		Timestamp: time.Now(),
		Size:      fileInfo.Size(),
	}

	select {
	case fw.processingQueue <- fileEvent:
		slog.Debug("File queued for processing", "path", event.Name)
	default:
		slog.Warn("Processing queue full, dropping file", "path", event.Name)
	}
}

// isImageFile checks if a file has a supported image extension
func (fw *FileWatcher) isImageFile(path string, extensions []string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	for _, supportedExt := range extensions {
		if ext == strings.ToLower(supportedExt) {
			return true
		}
	}
	return false
}

// shouldIgnoreFile checks if a file should be ignored based on patterns
func (fw *FileWatcher) shouldIgnoreFile(path string, patterns []string) bool {
	filename := filepath.Base(path)

	// Check against ignore patterns
	for _, pattern := range patterns {
		if strings.Contains(filename, pattern) {
			return true
		}
	}

	// Check for temporary files
	if strings.HasPrefix(filename, ".") || strings.HasSuffix(filename, "~") {
		return true
	}

	return false
}

// isDebounced checks if a file event should be debounced
func (fw *FileWatcher) isDebounced(path string) bool {
	fw.debounceMutex.Lock()
	defer fw.debounceMutex.Unlock()

	if lastTime, exists := fw.debounceMap[path]; exists {
		if time.Since(lastTime) < fw.debounceInterval {
			return true
		}
	}

	fw.debounceMap[path] = time.Now()
	return false
}

// processExistingFiles processes existing files in a directory
func (fw *FileWatcher) processExistingFiles(directory string) {
	entries, err := os.ReadDir(directory)
	if err != nil {
		slog.Warn("Failed to read directory for existing files", "directory", directory, "error", err)
		return
	}

	// Use default config for existing file processing
	config := DefaultWatcherConfig()

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filePath := filepath.Join(directory, entry.Name())

		if fw.isImageFile(filePath, config.FileExtensions) && !fw.shouldIgnoreFile(filePath, config.IgnorePatterns) {
			slog.Debug("Processing existing file", "path", filePath)

			fileInfo, err := entry.Info()
			if err != nil {
				continue
			}

			fileEvent := FileEvent{
				Path:      filePath,
				Operation: fw.operation,
				Timestamp: time.Now(),
				Size:      fileInfo.Size(),
			}

			select {
			case fw.processingQueue <- fileEvent:
				slog.Debug("Existing file queued for processing", "path", filePath)
			default:
				slog.Warn("Processing queue full, dropping existing file", "path", filePath)
			}
		}
	}
}

// startWorkers starts the worker pool for processing files
func (fw *FileWatcher) startWorkers() {
	for i := 0; i < fw.maxWorkers; i++ {
		fw.wg.Add(1)
		go fw.worker(i)
	}
	slog.Info("File watcher worker pool started", "workers", fw.maxWorkers)
}

// worker processes files from the queue
func (fw *FileWatcher) worker(id int) {
	defer fw.wg.Done()
	slog.Debug("File watcher worker started", "worker_id", id)

	for {
		select {
		case <-fw.ctx.Done():
			slog.Debug("File watcher worker stopping", "worker_id", id)
			return

		case fileEvent := <-fw.processingQueue:
			fw.processFile(fileEvent)
		}
	}
}

// processFile processes a single file event
func (fw *FileWatcher) processFile(event FileEvent) {
	slog.Info("Processing file with Google Lens",
		"path", event.Path,
		"operation", event.Operation,
		"size", event.Size)

	// Create progress callback
	progressCallback := func(progress float64, message string) {
		slog.Debug("Processing progress", "path", event.Path, "progress", progress, "message", message)
	}

	var result automation.LensResult
	var err error

	// Apply preprocessing if enabled
	if fw.autoPreprocess && event.Operation == "extract_text" {
		slog.Debug("Applying OCR preprocessing", "path", event.Path)
		result, err = fw.googleLens.ProcessImageWithPreprocessing(event.Path, event.Operation, "ocr", progressCallback)
	} else {
		result, err = fw.googleLens.ProcessImageWithProgress(event.Path, event.Operation, progressCallback)
	}

	if err != nil {
		slog.Error("Failed to process file", "path", event.Path, "error", err)
		if fw.notifySystem != nil && fw.notifyOnComplete {
			fw.notifySystem.SendError("Google Lens processing failed", fmt.Sprintf("File: %s\nError: %v", event.Path, err))
		}
		return
	}

	slog.Info("File processing completed",
		"path", event.Path,
		"operation", event.Operation,
		"processing_time", result.ProcessingTime,
		"cached", result.Cached)

	// Send notification if enabled and not cached
	if fw.notifySystem != nil && fw.notifyOnComplete && !result.Cached {
		title := fmt.Sprintf("Google Lens: %s", event.Operation)
		message := fmt.Sprintf("File: %s\nResult: %s", filepath.Base(event.Path), result.ResultText)
		if len(result.ResultText) > 100 {
			message = fmt.Sprintf("File: %s\nResult: %s...", filepath.Base(event.Path), result.ResultText[:97])
		}
		fw.notifySystem.SendSuccess(title, message)
	}
}

// Close stops the file watcher and cleans up resources
func (fw *FileWatcher) Close() error {
	slog.Info("Closing file watcher")

	// Cancel context to stop all goroutines
	fw.cancel()

	// Wait for workers to finish
	done := make(chan struct{})
	go func() {
		fw.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		slog.Debug("All file watcher workers stopped")
	case <-time.After(10 * time.Second):
		slog.Warn("File watcher workers did not stop gracefully within timeout")
	}

	// Close the watcher
	if err := fw.watcher.Close(); err != nil {
		slog.Error("Failed to close file watcher", "error", err)
		return err
	}

	slog.Info("File watcher closed successfully")
	return nil
}

// GetStats returns statistics about the file watcher
func (fw *FileWatcher) GetStats() map[string]interface{} {
	fw.pathsMutex.RLock()
	watchedPaths := len(fw.watchedPaths)
	fw.pathsMutex.RUnlock()

	fw.debounceMutex.RLock()
	debounceCount := len(fw.debounceMap)
	fw.debounceMutex.RUnlock()

	return map[string]interface{}{
		"watched_paths":      watchedPaths,
		"debounce_entries":   debounceCount,
		"queue_length":       len(fw.processingQueue),
		"max_workers":        fw.maxWorkers,
		"operation":          fw.operation,
		"auto_preprocess":    fw.autoPreprocess,
		"notify_on_complete": fw.notifyOnComplete,
		"enabled":            fw.enabled,
	}
}
