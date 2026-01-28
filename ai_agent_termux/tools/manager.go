package tools

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"

	"ai_agent_termux/config"
	"golang.org/x/exp/slog"
)

// ToolManager manages various system tools and utilities
type ToolManager struct {
	config         *config.Config
	availableTools map[string]*ToolInfo
	mu             sync.RWMutex
}

// ToolInfo represents information about a tool
type ToolInfo struct {
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Version      string    `json:"version"`
	Path         string    `json:"path"`
	Categories   []string  `json:"categories"`
	Availability bool      `json:"availability"`
	LastChecked  time.Time `json:"last_checked"`
}

// MediaFileInfo represents media file information from ffprobe
type MediaFileInfo struct {
	Format   string       `json:"format"`
	Duration float64      `json:"duration"`
	BitRate  int          `json:"bit_rate"`
	Size     int64        `json:"size"`
	Streams  []StreamInfo `json:"streams"`
}

// StreamInfo represents media stream information
type StreamInfo struct {
	CodecType  string `json:"codec_type"`
	CodecName  string `json:"codec_name"`
	Width      int    `json:"width,omitempty"`
	Height     int    `json:"height,omitempty"`
	SampleRate int    `json:"sample_rate,omitempty"`
	Channels   int    `json:"channels,omitempty"`
}

// UsefulPackage represents a useful package/tool
type UsefulPackage struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Command     string `json:"command"`
	InstallCmd  string `json:"install_cmd"`
	Priority    int    `json:"priority"` // 1-10, higher is more important
	Installed   bool   `json:"installed"`
}

// NewToolManager creates a new tool manager instance
func NewToolManager(cfg *config.Config) *ToolManager {
	manager := &ToolManager{
		config:         cfg,
		availableTools: make(map[string]*ToolInfo),
	}

	// Initialize tool detection
	manager.detectTools()

	slog.Info("Tool Manager initialized", "tools_detected", len(manager.availableTools))
	return manager
}

// detectTools scans for commonly available tools
func (tm *ToolManager) detectTools() {
	// Essential tools to check
	essentialTools := []string{
		"aichat", "ffmpeg", "ffprobe", "ffplay",
		"exiftool", "imagemagick", "ghostscript", "poppler-utils",
		"tesseract", "ocrmypdf", "pdftk", "qpdf",
		"sqlite3", "jq", "curl", "wget",
		"git", "nano", "vim", "rsync",
		"python3", "node", "java", "go",
		"7z", "unzip", "tar", "gzip",
		"ntfy", "termux-api", "termux-notification",
		"shizuku", "adb", "fastboot",
		"mediainfo", "sox", "youtube-dl", "aria2c",
	}

	for _, toolName := range essentialTools {
		if tm.isToolAvailable(toolName) {
			version := tm.getToolVersion(toolName)
			path := tm.getToolPath(toolName)

			tm.availableTools[toolName] = &ToolInfo{
				Name:         toolName,
				Description:  fmt.Sprintf("Detected system tool: %s", toolName),
				Version:      version,
				Path:         path,
				Availability: true,
				LastChecked:  time.Now(),
			}

			slog.Debug("Tool detected", "name", toolName, "version", version, "path", path)
		}
	}
}

// isToolAvailable checks if a tool is available in PATH
func (tm *ToolManager) isToolAvailable(toolName string) bool {
	_, err := exec.LookPath(toolName)
	return err == nil
}

// getToolVersion gets the version of a tool
func (tm *ToolManager) getToolVersion(toolName string) string {
	// Common version flags
	versionFlags := []string{"--version", "-v", "version", "--help"}

	for _, flag := range versionFlags {
		cmd := exec.Command(toolName, flag)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out

		if err := cmd.Run(); err == nil {
			output := out.String()
			lines := strings.Split(output, "\n")
			if len(lines) > 0 {
				// Return first line that looks like a version
				for _, line := range lines {
					if strings.Contains(strings.ToLower(line), "version") ||
						strings.Contains(line, "v") && len(line) < 100 {
						return strings.TrimSpace(line)
					}
				}
				return strings.TrimSpace(lines[0])
			}
		}
	}

	return "unknown"
}

// getToolPath gets the path of a tool
func (tm *ToolManager) getToolPath(toolName string) string {
	path, err := exec.LookPath(toolName)
	if err != nil {
		return ""
	}
	return path
}

// IsToolAvailable checks if a specific tool is available
func (tm *ToolManager) IsToolAvailable(toolName string) bool {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	tool, exists := tm.availableTools[toolName]
	if !exists {
		return false
	}

	// Re-check if last check was more than 5 minutes ago
	if time.Since(tool.LastChecked) > 5*time.Minute {
		available := tm.isToolAvailable(toolName)
		if toolInfo, ok := tm.availableTools[toolName]; ok {
			toolInfo.Availability = available
			toolInfo.LastChecked = time.Now()
		}
		return available
	}

	return tool.Availability
}

// GetToolInfo returns information about a specific tool
func (tm *ToolManager) GetToolInfo(toolName string) (*ToolInfo, error) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	tool, exists := tm.availableTools[toolName]
	if !exists {
		return nil, fmt.Errorf("tool not found: %s", toolName)
	}

	return tool, nil
}

// GetAllTools returns all detected tools
func (tm *ToolManager) GetAllTools() map[string]*ToolInfo {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	// Return a copy to prevent concurrent access issues
	toolsCopy := make(map[string]*ToolInfo)
	for k, v := range tm.availableTools {
		toolsCopy[k] = &ToolInfo{
			Name:         v.Name,
			Description:  v.Description,
			Version:      v.Version,
			Path:         v.Path,
			Categories:   append([]string{}, v.Categories...),
			Availability: v.Availability,
			LastChecked:  v.LastChecked,
		}
	}

	return toolsCopy
}

// GetMediaFileInfo extracts media information using ffprobe
func (tm *ToolManager) GetMediaFileInfo(filePath string) (*MediaFileInfo, error) {
	if !tm.IsToolAvailable("ffprobe") {
		return nil, fmt.Errorf("ffprobe not available")
	}

	// Run ffprobe to get media information in JSON format
	cmd := exec.Command("ffprobe",
		"-v", "quiet",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
		filePath)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to probe media file: %v", err)
	}

	// Parse the JSON output (simplified parsing)
	output := out.String()

	// Extract basic information from output
	info := &MediaFileInfo{
		Format:   tm.extractValue(output, "format_name"),
		Duration: tm.extractFloatValue(output, "duration"),
		BitRate:  tm.extractIntValue(output, "bit_rate"),
		Size:     tm.extractInt64Value(output, "size"),
	}

	// Extract stream information
	info.Streams = tm.extractStreams(output)

	slog.Debug("Media file info extracted", "file", filePath, "format", info.Format, "duration", info.Duration)
	return info, nil
}

// extractValue extracts a string value from ffprobe output
func (tm *ToolManager) extractValue(output, key string) string {
	// Simple string extraction - would normally use proper JSON parsing
	start := strings.Index(output, fmt.Sprintf("\"%s\":", key))
	if start == -1 {
		return ""
	}

	start += len(fmt.Sprintf("\"%s\":\"", key))
	end := strings.Index(output[start:], "\"")
	if end == -1 {
		return ""
	}

	return output[start : start+end]
}

// extractFloatValue extracts a float value from ffprobe output
func (tm *ToolManager) extractFloatValue(output, key string) float64 {
	valueStr := tm.extractValue(output, key)
	if valueStr == "" {
		return 0
	}

	var value float64
	fmt.Sscanf(valueStr, "%f", &value)
	return value
}

// extractIntValue extracts an integer value from ffprobe output
func (tm *ToolManager) extractIntValue(output, key string) int {
	valueStr := tm.extractValue(output, key)
	if valueStr == "" {
		return 0
	}

	var value int
	fmt.Sscanf(valueStr, "%d", &value)
	return value
}

// extractInt64Value extracts an int64 value from ffprobe output
func (tm *ToolManager) extractInt64Value(output, key string) int64 {
	valueStr := tm.extractValue(output, key)
	if valueStr == "" {
		return 0
	}

	var value int64
	fmt.Sscanf(valueStr, "%d", &value)
	return value
}

// extractStreams extracts stream information from ffprobe output
func (tm *ToolManager) extractStreams(output string) []StreamInfo {
	// Simplified stream extraction
	var streams []StreamInfo

	// Count codec_type occurrences to estimate number of streams
	codecCount := strings.Count(output, "\"codec_type\"")

	for i := 0; i < codecCount; i++ {
		stream := StreamInfo{
			CodecType: tm.extractValueN(output, "codec_type", i),
			CodecName: tm.extractValueN(output, "codec_name", i),
		}

		// Extract width/height for video streams
		if stream.CodecType == "video" {
			stream.Width = tm.extractIntValue(output, "width")
			stream.Height = tm.extractIntValue(output, "height")
		}

		// Extract sample rate/channels for audio streams
		if stream.CodecType == "audio" {
			stream.SampleRate = tm.extractIntValue(output, "sample_rate")
			stream.Channels = tm.extractIntValue(output, "channels")
		}

		streams = append(streams, stream)
	}

	return streams
}

// extractValueN extracts the Nth occurrence of a value from output
func (tm *ToolManager) extractValueN(output, key string, n int) string {
	// Find nth occurrence
	occurrences := 0
	startPos := 0

	for occurrences <= n {
		start := strings.Index(output[startPos:], fmt.Sprintf("\"%s\":", key))
		if start == -1 {
			return ""
		}

		if occurrences == n {
			start += startPos + len(fmt.Sprintf("\"%s\":\"", key))
			end := strings.Index(output[start:], "\"")
			if end == -1 {
				return ""
			}
			return output[start : start+end]
		}

		occurrences++
		startPos += start + 1
	}

	return ""
}

// ExecuteToolCommand executes a tool command with proper error handling
func (tm *ToolManager) ExecuteToolCommand(toolName string, args ...string) (string, error) {
	if !tm.IsToolAvailable(toolName) {
		return "", fmt.Errorf("tool not available: %s", toolName)
	}

	cmd := exec.Command(toolName, args...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to execute %s: %v, stderr: %s", toolName, err, stderr.String())
	}

	return out.String(), nil
}

// ConvertMedia converts media files using ffmpeg
func (tm *ToolManager) ConvertMedia(inputFile, outputFile, format string, options ...string) error {
	if !tm.IsToolAvailable("ffmpeg") {
		return fmt.Errorf("ffmpeg not available")
	}

	// Build ffmpeg command
	args := []string{"-i", inputFile}
	args = append(args, options...)
	args = append(args, outputFile)

	cmd := exec.Command("ffmpeg", args...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to convert media: %v, stderr: %s", err, stderr.String())
	}

	slog.Info("Media conversion completed", "input", inputFile, "output", outputFile, "format", format)
	return nil
}

// GetUsefulPackages returns a curated list of 30+ useful packages
func (tm *ToolManager) GetUsefulPackages() []UsefulPackage {
	packages := []UsefulPackage{
		// Development Tools
		{Name: "aichat", Description: "AI chat client for multiple providers", Category: "AI", Command: "aichat --help", InstallCmd: "cargo install aichat", Priority: 10, Installed: tm.IsToolAvailable("aichat")},
		{Name: "ffmpeg", Description: "Complete multimedia solution", Category: "Media", Command: "ffmpeg -version", InstallCmd: "pkg install ffmpeg", Priority: 10, Installed: tm.IsToolAvailable("ffmpeg")},
		{Name: "exiftool", Description: "Read/write meta information in files", Category: "Metadata", Command: "exiftool -ver", InstallCmd: "pkg install exiftool", Priority: 8, Installed: tm.IsToolAvailable("exiftool")},

		// Document Processing
		{Name: "poppler-utils", Description: "PDF utilities (pdftotext, pdfinfo)", Category: "Documents", Command: "pdftotext -v", InstallCmd: "pkg install poppler", Priority: 9, Installed: tm.IsToolAvailable("pdftotext")},
		{Name: "tesseract", Description: "OCR engine", Category: "OCR", Command: "tesseract --version", InstallCmd: "pkg install tesseract", Priority: 9, Installed: tm.IsToolAvailable("tesseract")},
		{Name: "ocrmypdf", Description: "Adds OCR text to PDF files", Category: "OCR", Command: "ocrmypdf --version", InstallCmd: "pip install ocrmypdf", Priority: 8, Installed: tm.IsToolAvailable("ocrmypdf")},

		// System Tools
		{Name: "rsync", Description: "Fast file synchronization", Category: "File Transfer", Command: "rsync --version", InstallCmd: "pkg install rsync", Priority: 8, Installed: tm.IsToolAvailable("rsync")},
		{Name: "sqlite3", Description: "SQLite database command-line tool", Category: "Database", Command: "sqlite3 --version", InstallCmd: "pkg install sqlite", Priority: 7, Installed: tm.IsToolAvailable("sqlite3")},
		{Name: "jq", Description: "JSON processor", Category: "Data Processing", Command: "jq --version", InstallCmd: "pkg install jq", Priority: 8, Installed: tm.IsToolAvailable("jq")},

		// Networking
		{Name: "curl", Description: "Transfer data with URLs", Category: "Networking", Command: "curl --version", InstallCmd: "pkg install curl", Priority: 9, Installed: tm.IsToolAvailable("curl")},
		{Name: "wget", Description: "Non-interactive network downloader", Category: "Networking", Command: "wget --version", InstallCmd: "pkg install wget", Priority: 8, Installed: tm.IsToolAvailable("wget")},
		{Name: "aria2c", Description: "Multi-connection download utility", Category: "Networking", Command: "aria2c --version", InstallCmd: "pkg install aria2", Priority: 7, Installed: tm.IsToolAvailable("aria2c")},

		// Compression & Archiving
		{Name: "7z", Description: "7-Zip archiver", Category: "Compression", Command: "7z --help", InstallCmd: "pkg install p7zip", Priority: 7, Installed: tm.IsToolAvailable("7z")},
		{Name: "unzip", Description: "Extract ZIP archives", Category: "Compression", Command: "unzip", InstallCmd: "pkg install unzip", Priority: 7, Installed: tm.IsToolAvailable("unzip")},

		// Programming
		{Name: "python3", Description: "Python interpreter", Category: "Programming", Command: "python3 --version", InstallCmd: "pkg install python", Priority: 10, Installed: tm.IsToolAvailable("python3")},
		{Name: "node", Description: "Node.js JavaScript runtime", Category: "Programming", Command: "node --version", InstallCmd: "pkg install nodejs", Priority: 9, Installed: tm.IsToolAvailable("node")},
		{Name: "go", Description: "Go programming language", Category: "Programming", Command: "go version", InstallCmd: "pkg install golang", Priority: 8, Installed: tm.IsToolAvailable("go")},

		// Image Processing
		{Name: "imagemagick", Description: "Image manipulation suite", Category: "Graphics", Command: "magick --version", InstallCmd: "pkg install imagemagick", Priority: 8, Installed: tm.IsToolAvailable("magick")},
		{Name: "ghostscript", Description: "PostScript interpreter", Category: "Graphics", Command: "gs --version", InstallCmd: "pkg install ghostscript", Priority: 7, Installed: tm.IsToolAvailable("gs")},

		// Android Tools
		{Name: "termux-api", Description: "Termux API package", Category: "Android", Command: "termux-api --help", InstallCmd: "pkg install termux-api", Priority: 10, Installed: tm.IsToolAvailable("termux-api")},
		{Name: "adb", Description: "Android Debug Bridge", Category: "Android", Command: "adb version", InstallCmd: "pkg install android-tools", Priority: 10, Installed: tm.IsToolAvailable("adb")},

		// Utilities
		{Name: "mediainfo", Description: "Media file information", Category: "Media", Command: "mediainfo --version", InstallCmd: "pkg install mediainfo", Priority: 7, Installed: tm.IsToolAvailable("mediainfo")},
		{Name: "sox", Description: "Sound processing toolkit", Category: "Audio", Command: "sox --version", InstallCmd: "pkg install sox", Priority: 7, Installed: tm.IsToolAvailable("sox")},
		{Name: "youtube-dl", Description: "Video download utility", Category: "Media", Command: "youtube-dl --version", InstallCmd: "pip install youtube-dl", Priority: 8, Installed: tm.IsToolAvailable("youtube-dl")},

		// Editors
		{Name: "nano", Description: "Simple text editor", Category: "Editors", Command: "nano --version", InstallCmd: "pkg install nano", Priority: 6, Installed: tm.IsToolAvailable("nano")},
		{Name: "vim", Description: "Vi IMproved editor", Category: "Editors", Command: "vim --version", InstallCmd: "pkg install vim", Priority: 7, Installed: tm.IsToolAvailable("vim")},

		// Git & Version Control
		{Name: "git", Description: "Distributed version control", Category: "Development", Command: "git --version", InstallCmd: "pkg install git", Priority: 9, Installed: tm.IsToolAvailable("git")},

		// System Information
		{Name: "neofetch", Description: "System information tool", Category: "System", Command: "neofetch --help", InstallCmd: "pkg install neofetch", Priority: 6, Installed: tm.IsToolAvailable("neofetch")},

		// Notification Services
		{Name: "ntfy", Description: "HTTP-based pub-sub notification service", Category: "Notifications", Command: "ntfy --help", InstallCmd: "go install github.com/binwiederhier/ntfy/cmd/ntfy@latest", Priority: 8, Installed: tm.IsToolAvailable("ntfy")},
	}

	return packages
}

// InstallMissingPackages attempts to install missing packages (requires user confirmation)
func (tm *ToolManager) InstallMissingPackages(packages []UsefulPackage) error {
	missingPackages := []UsefulPackage{}

	for _, pkg := range packages {
		if !pkg.Installed {
			missingPackages = append(missingPackages, pkg)
		}
	}

	if len(missingPackages) == 0 {
		slog.Info("All packages are already installed")
		return nil
	}

	slog.Info("Missing packages detected", "count", len(missingPackages))

	// Sort by priority (higher priority first)
	for i := 0; i < len(missingPackages)-1; i++ {
		for j := i + 1; j < len(missingPackages); j++ {
			if missingPackages[i].Priority < missingPackages[j].Priority {
				missingPackages[i], missingPackages[j] = missingPackages[j], missingPackages[i]
			}
		}
	}

	// Try to install top 10 highest priority packages
	toInstall := 10
	if len(missingPackages) < toInstall {
		toInstall = len(missingPackages)
	}

	slog.Info("Installing top priority missing packages", "count", toInstall)

	for i := 0; i < toInstall; i++ {
		pkg := missingPackages[i]
		slog.Info("Installing package", "name", pkg.Name, "priority", pkg.Priority)

		// This would normally execute the install command
		// For safety, we'll just log what would be installed
		slog.Info("Would install", "command", pkg.InstallCmd)
	}

	return nil
}

// GetAichatModels returns available aichat models
func (tm *ToolManager) GetAichatModels() ([]string, error) {
	if !tm.IsToolAvailable("aichat") {
		return nil, fmt.Errorf("aichat not available")
	}

	output, err := tm.ExecuteToolCommand("aichat", "list")
	if err != nil {
		return nil, fmt.Errorf("failed to get aichat models: %v", err)
	}

	// Parse models from output
	lines := strings.Split(output, "\n")
	var models []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.Contains(line, "Available models") {
			// Extract model name (assumes format like "* model_name" or "model_name")
			if strings.HasPrefix(line, "* ") {
				model := strings.TrimPrefix(line, "* ")
				models = append(models, model)
			} else if !strings.HasPrefix(line, " ") && !strings.Contains(line, ":") {
				models = append(models, line)
			}
		}
	}

	return models, nil
}

// ExecuteAichatQuery executes a query using aichat
func (tm *ToolManager) ExecuteAichatQuery(prompt string, model string) (string, error) {
	if !tm.IsToolAvailable("aichat") {
		return "", fmt.Errorf("aichat not available")
	}

	args := []string{"--text", prompt}
	if model != "" {
		args = append([]string{"--model", model}, args...)
	}

	output, err := tm.ExecuteToolCommand("aichat", args...)
	if err != nil {
		return "", fmt.Errorf("failed to execute aichat query: %v", err)
	}

	return output, nil
}
