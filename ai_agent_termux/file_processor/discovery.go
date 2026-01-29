package file_processor

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"ai_agent_termux/pkg/lowlevel"

	"golang.org/x/exp/slog"
)

type FileMetadata struct {
	Path        string    `json:"path"`
	Filename    string    `json:"filename"`
	Ext         string    `json:"extension"`
	Type        string    `json:"type"` // "image", "pdf", "text", "other"
	Size        int64     `json:"size"`
	ModTime     time.Time `json:"mod_time"`
	IsDir       bool      `json:"is_dir"`
	Processed   bool      `json:"processed"` // Track if processed by the agent
	LastProcess time.Time `json:"last_process,omitempty"`
	ContentHash string    `json:"content_hash,omitempty"`
	Priority    int       `json:"priority"` // Priority for processing (higher = more important)
}

type DirectoryStats struct {
	Path          string `json:"path"`
	TotalFiles    int    `json:"total_files"`
	ImageFiles    int    `json:"image_files"`
	PDFFiles      int    `json:"pdf_files"`
	TextFiles     int    `json:"text_files"`
	OtherFiles    int    `json:"other_files"`
	TotalSize     int64  `json:"total_size"`
	IsLargeDir    bool   `json:"is_large_dir"`   // true if >= 20 files
	RequiresBatch bool   `json:"requires_batch"` // true if needs special batch handling
}

// DiscoverFilesFast is a high-speed alternative to filepath.Walk using low-level syscalls
func DiscoverFilesFast(dir string) ([]string, error) {
	var allFiles []string

	names, err := lowlevel.FastScanDir(dir)
	if err != nil {
		return nil, err
	}

	for _, name := range names {
		if name == "." || name == ".." {
			continue
		}
		path := filepath.Join(dir, name)
		info, err := os.Lstat(path)
		if err != nil {
			continue
		}

		if info.IsDir() {
			subFiles, _ := DiscoverFilesFast(path)
			allFiles = append(allFiles, subFiles...)
		} else {
			allFiles = append(allFiles, path)
		}
	}
	return allFiles, nil
}

// DiscoverFilesEnhanced provides enhanced file discovery with directory analysis
func DiscoverFilesEnhanced(dirs []string) ([]FileMetadata, []DirectoryStats, error) {
	var allFiles []FileMetadata
	var dirStats []DirectoryStats

	for _, dir := range dirs {
		slog.Info("Scanning directory with fast syscalls", "directory", dir)

		// Use fast scanner for initial file list if it's potentially large
		paths, err := DiscoverFilesFast(dir)
		if err != nil {
			slog.Warn("Fast discovery failed, falling back to standard walk", "directory", dir, "error", err)
			files, stats, err := discoverDirectoryEnhanced(dir)
			if err == nil {
				allFiles = append(allFiles, files...)
				dirStats = append(dirStats, stats)
			}
			continue
		}

		// Process the discovered paths into metadata
		stats := DirectoryStats{Path: dir}
		for _, path := range paths {
			info, err := os.Stat(path)
			if err != nil {
				continue
			}

			ext := strings.ToLower(filepath.Ext(path))
			fileType := "other"
			switch ext {
			case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp":
				fileType = "image"
				stats.ImageFiles++
			case ".pdf":
				fileType = "pdf"
				stats.PDFFiles++
			case ".mp3", ".wav", ".flac", ".m4a", ".ogg":
				fileType = "audio"
				stats.OtherFiles++
			case ".mp4", ".mkv", ".avi", ".mov", ".wmv":
				fileType = "video"
				stats.OtherFiles++
			case ".txt", ".md", ".log", ".csv", ".json", ".xml", ".html", ".htm",
				".go", ".py", ".js", ".ts", ".c", ".cpp", ".h", ".cs", ".java", ".sh", ".rs":
				fileType = "text"
				stats.TextFiles++
			default:
				stats.OtherFiles++
			}

			stats.TotalFiles++
			stats.TotalSize += info.Size()
			priority := calculateFilePriority(fileType, info.Size(), info.ModTime())

			// Zero-copy hashing (already implemented in Phase 0)
			hash, _ := CalculateFileHash(path)

			allFiles = append(allFiles, FileMetadata{
				Path:        path,
				Filename:    info.Name(),
				Ext:         ext,
				Type:        fileType,
				Size:        info.Size(),
				ModTime:     info.ModTime(),
				ContentHash: hash,
				Priority:    priority,
			})
		}

		stats.IsLargeDir = stats.TotalFiles >= 20
		stats.RequiresBatch = stats.TotalFiles >= 50
		dirStats = append(dirStats, stats)
	}

	return allFiles, dirStats, nil
}

// CalculateFileHash computes a fast hash of the file content using zero-copy mmap
func CalculateFileHash(path string) (string, error) {
	data, err := lowlevel.MmapFile(path)
	if err != nil {
		// Fallback to traditional ReadFile if mmap fails
		data, err = os.ReadFile(path)
		if err != nil {
			return "", err
		}
		hash := lowlevel.FastHashASM(data)
		return fmt.Sprintf("%x", hash), nil
	}
	defer lowlevel.Munmap(data)

	hash := lowlevel.FastHashASM(data)
	return fmt.Sprintf("%x", hash), nil
}

// discoverDirectoryEnhanced discovers files in a directory with detailed statistics
func discoverDirectoryEnhanced(dir string) ([]FileMetadata, DirectoryStats, error) {
	var files []FileMetadata
	stats := DirectoryStats{Path: dir}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			slog.Warn("Error accessing path", "path", path, "error", err)
			return filepath.SkipDir
		}

		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		fileType := "other"

		switch ext {
		case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp":
			fileType = "image"
			stats.ImageFiles++
		case ".pdf":
			fileType = "pdf"
			stats.PDFFiles++
		case ".txt", ".md", ".log", ".csv", ".json", ".xml", ".html", ".htm":
			fileType = "text"
			stats.TextFiles++
		default:
			stats.OtherFiles++
		}

		stats.TotalFiles++
		stats.TotalSize += info.Size()

		// Calculate priority based on file type and size
		priority := calculateFilePriority(fileType, info.Size(), info.ModTime())

		// Calculate fast hash for changed/new files (optional: could be deferred)
		hash, _ := CalculateFileHash(path)

		metadata := FileMetadata{
			Path:        path,
			Filename:    info.Name(),
			Ext:         ext,
			Type:        fileType,
			Size:        info.Size(),
			ModTime:     info.ModTime(),
			IsDir:       false,
			ContentHash: hash,
			Priority:    priority,
		}

		files = append(files, metadata)
		return nil
	})

	if err != nil {
		return files, stats, err
	}

	// Determine if this is a large directory requiring batch processing
	stats.IsLargeDir = stats.TotalFiles >= 20
	stats.RequiresBatch = stats.TotalFiles >= 50 // Batch processing for 50+ files

	slog.Info("Directory scan complete",
		"directory", dir,
		"total_files", stats.TotalFiles,
		"is_large", stats.IsLargeDir,
		"requires_batch", stats.RequiresBatch)

	return files, stats, nil
}

// calculateFilePriority determines processing priority for files
func calculateFilePriority(fileType string, size int64, modTime time.Time) int {
	basePriority := 1

	// Higher priority for text files (easier to process)
	switch fileType {
	case "text":
		basePriority = 10
	case "pdf":
		basePriority = 8
	case "image":
		basePriority = 6
	default:
		basePriority = 3
	}

	// Adjust for size - medium-sized files get highest priority
	switch {
	case size < 1024: // < 1KB
		basePriority += 1
	case size < 1024*1024: // < 1MB
		basePriority += 3
	case size < 10*1024*1024: // < 10MB
		basePriority += 2
	default: // >= 10MB
		basePriority -= 1
	}

	// Boost recently modified files
	weeksSinceModified := time.Since(modTime).Hours() / 24 / 7
	switch {
	case weeksSinceModified < 1:
		basePriority += 2
	case weeksSinceModified < 4:
		basePriority += 1
	}

	return basePriority
}

// GroupFilesForBatchProcessing groups files for efficient batch processing
func GroupFilesForBatchProcessing(files []FileMetadata, batchSize int) [][]FileMetadata {
	if batchSize <= 0 {
		batchSize = 20 // Default batch size
	}

	// Sort files by priority (highest first) and then by modification time
	sort.Slice(files, func(i, j int) bool {
		if files[i].Priority != files[j].Priority {
			return files[i].Priority > files[j].Priority
		}
		return files[i].ModTime.After(files[j].ModTime)
	})

	var batches [][]FileMetadata

	for i := 0; i < len(files); i += batchSize {
		end := i + batchSize
		if end > len(files) {
			end = len(files)
		}

		batch := files[i:end]
		batches = append(batches, batch)
	}

	slog.Info("Files grouped for batch processing",
		"total_files", len(files),
		"batch_size", batchSize,
		"num_batches", len(batches))

	return batches
}

// ProcessLargeDirectory processes directories with many files efficiently
func ProcessLargeDirectory(files []FileMetadata, workers int, batchSize int, processor func(FileMetadata) error) error {
	if workers <= 0 {
		workers = runtime.NumCPU()
		if workers > 8 {
			workers = 8 // Cap for mobile devices
		}
	}
	if batchSize <= 0 {
		batchSize = 20
	}

	if len(files) < batchSize {
		// Process sequentially for small batches
		for _, file := range files {
			if err := processor(file); err != nil {
				slog.Error("Error processing file", "file", file.Path, "error", err)
			}
		}
		return nil
	}

	slog.Info("Starting batch processing", "total_files", len(files), "workers", workers)

	// Group files into batches
	batches := GroupFilesForBatchProcessing(files, batchSize)

	// Create worker pool
	var wg sync.WaitGroup
	errChan := make(chan error, len(batches))
	semaphore := make(chan struct{}, workers)

	for i, batch := range batches {
		wg.Add(1)
		go func(batchIndex int, batchFiles []FileMetadata) {
			defer wg.Done()

			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			slog.Info("Processing batch", "batch", batchIndex, "files", len(batchFiles))

			for _, file := range batchFiles {
				if err := processor(file); err != nil {
					errChan <- fmt.Errorf("batch %d, file %s: %v", batchIndex, file.Path, err)
					return
				}
			}

			slog.Info("Completed batch", "batch", batchIndex, "files", len(batchFiles))
		}(i, batch)
	}

	// Wait for all workers to complete
	wg.Wait()
	close(errChan)

	// Check for errors
	var errors []string
	for err := range errChan {
		errors = append(errors, err.Error())
	}

	if len(errors) > 0 {
		slog.Warn("Batch processing completed with errors", "error_count", len(errors))
		// Continue processing despite some errors
	}

	slog.Info("Batch processing completed", "total_files", len(files), "errors", len(errors))
	return nil
}

// DiscoverFiles maintains backward compatibility
func DiscoverFiles(dirs []string) ([]FileMetadata, error) {
	files, _, err := DiscoverFilesEnhanced(dirs)
	return files, err
}

// GetDirectorySummary provides a summary of discovered directories
func GetDirectorySummary(stats []DirectoryStats) string {
	var totalFiles int
	var totalSize int64
	var largeDirs int
	var batchDirs int

	for _, stat := range stats {
		totalFiles += stat.TotalFiles
		totalSize += stat.TotalSize
		if stat.IsLargeDir {
			largeDirs++
		}
		if stat.RequiresBatch {
			batchDirs++
		}
	}

	return fmt.Sprintf(`Directory Scan Summary:
========================
Total Directories: %d
Total Files: %d
Total Size: %.2f MB
Large Directories (≥20 files): %d
Directories Requiring Batch Processing (≥50 files): %d`,
		len(stats), totalFiles, float64(totalSize)/1024/1024, largeDirs, batchDirs)
}
