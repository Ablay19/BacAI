package utils

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"golang.org/x/exp/slog"
)

var (
	logger *slog.Logger
)

func InitLogger(logPath string) {
	// Create or open log file
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	// Configure slog with custom options
	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Customize the time format
			if a.Key == slog.TimeKey {
				return slog.String(slog.TimeKey, a.Value.Time().Format(time.RFC3339))
			}
			return a
		},
	}

	// Use different handlers based on environment
	var handler slog.Handler
	if os.Getenv("DEBUG") == "true" {
		handler = slog.NewTextHandler(os.Stdout, opts)
	} else {
		handler = slog.NewJSONHandler(file, opts)
	}

	logger = slog.New(handler)
	slog.SetDefault(logger)
}

func LogInfo(message string, args ...any) {
	if logger != nil {
		logger.Info(message, args...)
	} else {
		slog.Info(message, args...)
	}
}

func LogError(message string, args ...any) {
	if logger != nil {
		logger.Error(message, args...)
	} else {
		slog.Error(message, args...)
	}
}

func LogWarning(message string, args ...any) {
	if logger != nil {
		logger.Warn(message, args...)
	} else {
		slog.Warn(message, args...)
	}
}

func LogDebug(message string, args ...any) {
	if logger != nil {
		logger.Debug(message, args...)
	} else {
		slog.Debug(message, args...)
	}
}

// ExecuteCommand runs a shell command and returns its stdout or an error.
func ExecuteCommand(name string, arg ...string) (string, error) {
	return ExecuteCommandWithTimeout(30*time.Second, name, arg...)
}

// ExecuteCommandWithTimeout runs a shell command with a timeout and returns its stdout or an error.
func ExecuteCommandWithTimeout(timeout time.Duration, name string, arg ...string) (string, error) {
	LogInfo("Executing command", "command", name, "args", strings.Join(arg, " "))

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, arg...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	// Check if the error was due to timeout
	if ctx.Err() == context.DeadlineExceeded {
		LogError("Command timed out", "command", name, "args", strings.Join(arg, " "), "timeout", timeout)
		return "", fmt.Errorf("command timed out after %v: %s %s", timeout, name, strings.Join(arg, " "))
	}

	if err != nil {
		LogError("Command failed",
			"command", name,
			"args", strings.Join(arg, " "),
			"error", err,
			"stderr", stderr.String())
		return "", fmt.Errorf("command execution failed: %v\nStderr: %s", err, stderr.String())
	}

	result := strings.TrimSpace(stdout.String())
	LogInfo("Command successful", "command", name, "output_length", len(result))
	return result, nil
}

// ExecuteCommandWithEnv runs a shell command with environment variables
func ExecuteCommandWithEnv(env map[string]string, name string, arg ...string) (string, error) {
	LogInfo("Executing command with env", "command", name, "args", strings.Join(arg, " "))

	cmd := exec.Command(name, arg...)

	// Set environment variables
	if env != nil {
		for k, v := range env {
			cmd.Env = append(cmd.Environ(), fmt.Sprintf("%s=%s", k, v))
		}
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		LogError("Command with env failed",
			"command", name,
			"args", strings.Join(arg, " "),
			"error", err,
			"stderr", stderr.String())
		return "", fmt.Errorf("command execution failed: %v\nStderr: %s", err, stderr.String())
	}

	result := strings.TrimSpace(stdout.String())
	LogInfo("Command with env successful", "command", name, "output_length", len(result))
	return result, nil
}

// HashString generates SHA256 hash for a given string
func HashString(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// FilterEmpty removes empty strings from a slice
func FilterEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if strings.TrimSpace(str) != "" {
			r = append(r, str)
		}
	}
	return r
}

// Min returns the smaller of two integers
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// InterfaceSliceToStringSlice converts []interface{} to []string
func InterfaceSliceToStringSlice(in []interface{}) []string {
	var out []string
	for _, v := range in {
		if s, ok := v.(string); ok {
			out = append(out, s)
		}
	}
	return out
}

// GetCurrentOS returns the current operating system
func GetCurrentOS() string {
	return runtime.GOOS
}

// FileExists checks if a file exists
func FileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	return !os.IsNotExist(err)
}

// DirectoryExists checks if a directory exists
func DirectoryExists(dirpath string) bool {
	info, err := os.Stat(dirpath)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}
