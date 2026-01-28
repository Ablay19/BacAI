package realtime

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"ai_agent_termux/config"
	"golang.org/x/exp/slog"
)

// LogManager handles real-time logging and notifications
type LogManager struct {
	config     *config.Config
	ctx        context.Context
	cancel     context.CancelFunc
	mu         sync.RWMutex
	logEntries []LogEntry
}

// LogEntry represents a log entry with metadata
type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	Source    string    `json:"source"`
	Tags      []string  `json:"tags"`
}

// NotificationManager handles various notification methods
type NotificationManager struct {
	config    *config.Config
	ntfyURL   string
	ntfyTopic string
	useTermux bool
	useNtfy   bool
}

// NewLogManager creates a new log manager instance
func NewLogManager(cfg *config.Config) *LogManager {
	ctx, cancel := context.WithCancel(context.Background())

	manager := &LogManager{
		config:     cfg,
		ctx:        ctx,
		cancel:     cancel,
		logEntries: make([]LogEntry, 0, 1000), // Pre-allocate for performance
	}

	// Start background log processing
	go manager.processLogs()

	slog.Info("Log Manager initialized")
	return manager
}

// NewNotificationManager creates a new notification manager instance
func NewNotificationManager(cfg *config.Config) *NotificationManager {
	manager := &NotificationManager{
		config:    cfg,
		ntfyURL:   "https://ntfy.sh", // Default ntfy server
		ntfyTopic: "",                // Will be set from config or generated
		useTermux: false,
		useNtfy:   false,
	}

	// Check available notification methods
	manager.detectNotificationMethods()

	slog.Info("Notification Manager initialized",
		"termux_available", manager.useTermux,
		"ntfy_available", manager.useNtfy)

	return manager
}

// detectNotificationMethods checks what notification methods are available
func (nm *NotificationManager) detectNotificationMethods() {
	// Check Termux API
	if _, err := exec.LookPath("termux-notification"); err == nil {
		nm.useTermux = true
		slog.Debug("Termux notification API detected")
	}

	// Check ntfy
	if _, err := exec.LookPath("ntfy"); err == nil {
		nm.useNtfy = true
		slog.Debug("ntfy client detected")
	}

	// Check if we have ntfy.sh URL in config or environment
	if nm.config.NtfyURL != "" {
		nm.ntfyURL = nm.config.NtfyURL
		nm.useNtfy = true
	} else if ntfyURL := getEnvWithDefault("NTFY_URL", ""); ntfyURL != "" {
		nm.ntfyURL = ntfyURL
		nm.useNtfy = true
	}

	// Check for topic
	if nm.config.NtfyTopic != "" {
		nm.ntfyTopic = nm.config.NtfyTopic
	} else if topic := getEnvWithDefault("NTFY_TOPIC", ""); topic != "" {
		nm.ntfyTopic = topic
	} else {
		// Generate default topic based on hostname
		if hostname, err := exec.Command("hostname").Output(); err == nil {
			nm.ntfyTopic = fmt.Sprintf("ai_agent_%s", strings.TrimSpace(string(hostname)))
		} else {
			nm.ntfyTopic = "ai_agent_notifications"
		}
	}
}

// getEnvWithDefault gets environment variable or returns default value
func getEnvWithDefault(key, defaultValue string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return defaultValue
}

// Log adds a new log entry
func (lm *LogManager) Log(level, message, source string, tags ...string) {
	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   message,
		Source:    source,
		Tags:      tags,
	}

	lm.mu.Lock()
	lm.logEntries = append(lm.logEntries, entry)

	// Keep only last 1000 entries to prevent memory issues
	if len(lm.logEntries) > 1000 {
		lm.logEntries = lm.logEntries[len(lm.logEntries)-1000:]
	}
	lm.mu.Unlock()

	// Also log using standard slog
	switch level {
	case "ERROR":
		slog.Error(message, "source", source, "tags", tags)
	case "WARN":
		slog.Warn(message, "source", source, "tags", tags)
	case "INFO":
		slog.Info(message, "source", source, "tags", tags)
	case "DEBUG":
		slog.Debug(message, "source", source, "tags", tags)
	default:
		slog.Info(message, "source", source, "tags", tags)
	}
}

// GetRecentLogs returns recent log entries
func (lm *LogManager) GetRecentLogs(count int) []LogEntry {
	lm.mu.RLock()
	defer lm.mu.RUnlock()

	if count > len(lm.logEntries) {
		count = len(lm.logEntries)
	}

	if count <= 0 {
		return []LogEntry{}
	}

	// Return last 'count' entries
	start := len(lm.logEntries) - count
	return append([]LogEntry(nil), lm.logEntries[start:]...)
}

// SearchLogs searches for log entries matching criteria
func (lm *LogManager) SearchLogs(query string, levels ...string) []LogEntry {
	lm.mu.RLock()
	defer lm.mu.RUnlock()

	var results []LogEntry

	for _, entry := range lm.logEntries {
		// Match query in message or tags
		matchesQuery := query == "" ||
			strings.Contains(strings.ToLower(entry.Message), strings.ToLower(query)) ||
			containsTag(entry.Tags, query)

		// Match levels if specified
		matchesLevel := len(levels) == 0
		for _, level := range levels {
			if strings.EqualFold(entry.Level, level) {
				matchesLevel = true
				break
			}
		}

		if matchesQuery && (matchesLevel || len(levels) == 0) {
			results = append(results, entry)
		}
	}

	return results
}

// containsTag checks if tags contain the query string
func containsTag(tags []string, query string) bool {
	queryLower := strings.ToLower(query)
	for _, tag := range tags {
		if strings.Contains(strings.ToLower(tag), queryLower) {
			return true
		}
	}
	return false
}

// processLogs runs background log processing
func (lm *LogManager) processLogs() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-lm.ctx.Done():
			return
		case <-ticker.C:
			lm.cleanupOldLogs()
		}
	}
}

// cleanupOldLogs removes logs older than 24 hours
func (lm *LogManager) cleanupOldLogs() {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	cutoff := time.Now().Add(-24 * time.Hour)

	// Find first entry that's within the cutoff time
	firstValid := len(lm.logEntries)
	for i, entry := range lm.logEntries {
		if entry.Timestamp.After(cutoff) {
			firstValid = i
			break
		}
	}

	// Remove old entries
	if firstValid > 0 && firstValid < len(lm.logEntries) {
		lm.logEntries = lm.logEntries[firstValid:]
		slog.Debug("Cleaned up old log entries", "removed", firstValid)
	}
}

// SendNotification sends a notification using available methods
func (nm *NotificationManager) SendNotification(title, message string, priority string) error {
	var errors []string

	// Send via Termux if available
	if nm.useTermux {
		if err := nm.sendTermuxNotification(title, message); err != nil {
			errors = append(errors, fmt.Sprintf("termux: %v", err))
		} else {
			slog.Debug("Notification sent via Termux", "title", title)
		}
	}

	// Send via ntfy if available
	if nm.useNtfy {
		if err := nm.sendNtfyNotification(title, message, priority); err != nil {
			errors = append(errors, fmt.Sprintf("ntfy: %v", err))
		} else {
			slog.Debug("Notification sent via ntfy", "title", title)
		}
	}

	if len(errors) > 0 && !(nm.useTermux || nm.useNtfy) {
		return fmt.Errorf("no notification methods available")
	}

	if len(errors) > 0 {
		return fmt.Errorf("some notifications failed: %s", strings.Join(errors, ", "))
	}

	return nil
}

// sendTermuxNotification sends notification via Termux API
func (nm *NotificationManager) sendTermuxNotification(title, message string) error {
	if !nm.useTermux {
		return fmt.Errorf("Termux notification not available")
	}

	// Build command
	args := []string{}
	if title != "" {
		args = append(args, "--title", title)
	}
	args = append(args, "--content", message)

	cmd := exec.Command("termux-notification", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to send Termux notification: %v, stderr: %s", err, stderr.String())
	}

	return nil
}

// sendNtfyNotification sends notification via ntfy.sh
func (nm *NotificationManager) sendNtfyNotification(title, message, priority string) error {
	if !nm.useNtfy {
		return fmt.Errorf("ntfy notification not available")
	}

	// If we have ntfy binary, use it
	if _, err := exec.LookPath("ntfy"); err == nil {
		return nm.sendNtfyBinaryNotification(title, message, priority)
	}

	// Otherwise use curl/HTTP if URL is available
	if nm.ntfyURL != "" {
		return nm.sendNtfyHTTPNotification(title, message, priority)
	}

	return fmt.Errorf("no ntfy method available")
}

// sendNtfyBinaryNotification sends notification using ntfy binary
func (nm *NotificationManager) sendNtfyBinaryNotification(title, message, priority string) error {
	args := []string{"publish"}

	if nm.ntfyTopic != "" {
		args = append(args, "--topic", nm.ntfyTopic)
	}

	if title != "" {
		args = append(args, "--title", title)
	}

	if priority != "" {
		args = append(args, "--priority", priority)
	}

	args = append(args, message)

	cmd := exec.Command("ntfy", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to send ntfy notification: %v, stderr: %s", err, stderr.String())
	}

	return nil
}

// sendNtfyHTTPNotification sends notification using HTTP via curl
func (nm *NotificationManager) sendNtfyHTTPNotification(title, message, priority string) error {
	url := fmt.Sprintf("%s/%s", nm.ntfyURL, nm.ntfyTopic)

	// Build curl command
	args := []string{"-X", "POST", url, "-d", message}

	if title != "" {
		args = append(args, "-H", fmt.Sprintf("Title: %s", title))
	}

	if priority != "" {
		args = append(args, "-H", fmt.Sprintf("Priority: %s", priority))
	}

	cmd := exec.Command("curl", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to send ntfy HTTP notification: %v, stderr: %s", err, stderr.String())
	}

	return nil
}

// WatchFileChanges watches for file changes and sends notifications
func (nm *NotificationManager) WatchFileChanges(paths []string, callback func(string, string)) error {
	// This would typically use fsnotify or similar
	// For now, we'll simulate with a simple polling approach

	go func() {
		lastCheck := time.Now()
		for {
			select {
			case <-time.After(5 * time.Second):
				// Check file modification times
				for _, path := range paths {
					// In real implementation, you'd check actual file mtimes
					// For simulation, we'll just trigger occasionally
					if time.Since(lastCheck) > 30*time.Second {
						callback(path, "modified")
						lastCheck = time.Now()
					}
				}
			}
		}
	}()

	return nil
}

// Close shuts down the log manager
func (lm *LogManager) Close() error {
	lm.cancel()
	return nil
}

// GetLogStatistics returns log statistics
func (lm *LogManager) GetLogStatistics() map[string]int {
	lm.mu.RLock()
	defer lm.mu.RUnlock()

	stats := map[string]int{
		"total_entries": len(lm.logEntries),
		"error_count":   0,
		"warn_count":    0,
		"info_count":    0,
		"debug_count":   0,
	}

	for _, entry := range lm.logEntries {
		switch strings.ToUpper(entry.Level) {
		case "ERROR":
			stats["error_count"]++
		case "WARN":
			stats["warn_count"]++
		case "INFO":
			stats["info_count"]++
		case "DEBUG":
			stats["debug_count"]++
		}
	}

	return stats
}
