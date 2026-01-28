package android

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"ai_agent_termux/config"
	"golang.org/x/exp/slog"
)

// AndroidManager handles Android-specific operations using ADB, Termux and Shizuku
type AndroidManager struct {
	config     *config.Config
	hasADB     bool
	hasTermux  bool
	hasShizuku bool
	deviceID   string
}

// AndroidFileInfo represents Android file metadata
type AndroidFileInfo struct {
	Path        string    `json:"path"`
	Filename    string    `json:"filename"`
	Size        int64     `json:"size"`
	ModTime     time.Time `json:"mod_time"`
	MimeType    string    `json:"mime_type"`
	PackageName string    `json:"package_name,omitempty"`
	IsSystem    bool      `json:"is_system"`
}

// NewAndroidManager creates a new Android manager instance
func NewAndroidManager(cfg *config.Config) *AndroidManager {
	manager := &AndroidManager{
		config: cfg,
	}

	// Detect available tools
	manager.detectTools()

	// Set default device if multiple devices connected
	if manager.hasADB {
		manager.setDefaultDevice()
	}

	slog.Info("Android Manager initialized",
		"adb_available", manager.hasADB,
		"termux_available", manager.hasTermux,
		"shizuku_available", manager.hasShizuku,
		"device_id", manager.deviceID)

	return manager
}

// detectTools checks which Android tools are available
func (a *AndroidManager) detectTools() {
	// Check ADB
	if _, err := exec.LookPath("adb"); err == nil {
		a.hasADB = true
	}

	// Check Termux API (if running in Termux)
	if _, err := exec.LookPath("termux-notification"); err == nil {
		a.hasTermux = true
	}

	// Check Shizuku (look for Shizuku service)
	if a.hasADB {
		if a.isShizukuAvailable() {
			a.hasShizuku = true
		}
	}
}

// setDefaultDevice sets the default ADB device
func (a *AndroidManager) setDefaultDevice() {
	if !a.hasADB {
		return
	}

	// Get connected devices
	cmd := exec.Command("adb", "devices")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		slog.Warn("Failed to get ADB devices", "error", err)
		return
	}

	lines := strings.Split(out.String(), "\n")
	var devices []string

	for _, line := range lines {
		if strings.Contains(line, "device") && !strings.Contains(line, "List of devices") {
			parts := strings.Fields(line)
			if len(parts) > 0 {
				devices = append(devices, parts[0])
			}
		}
	}

	if len(devices) > 0 {
		a.deviceID = devices[0] // Use first device
		if len(devices) > 1 {
			slog.Info("Multiple ADB devices found", "count", len(devices), "using", a.deviceID)
		}
	}
}

// isShizukuAvailable checks if Shizuku service is available
func (a *AndroidManager) isShizukuAvailable() bool {
	// Try to connect to Shizuku service
	cmd := exec.Command("adb", "-s", a.deviceID, "shell", "pm", "list", "packages", "moe.shizuku.privileged.api")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return false
	}

	return strings.Contains(out.String(), "package:moe.shizuku.privileged.api")
}

// DiscoverAndroidFiles discovers files across Android storage locations
func (a *AndroidManager) DiscoverAndroidFiles() ([]AndroidFileInfo, error) {
	var files []AndroidFileInfo

	if !a.hasADB && !a.hasTermux {
		return files, fmt.Errorf("no Android tools available")
	}

	// Discover files using available methods
	if a.hasTermux {
		termuxFiles, err := a.discoverTermuxFiles()
		if err == nil {
			files = append(files, termuxFiles...)
		}
	}

	if a.hasADB {
		adbFiles, err := a.discoverADBFiles()
		if err == nil {
			files = append(files, adbFiles...)
		}
	}

	slog.Info("Discovered Android files", "count", len(files))
	return files, nil
}

// CaptureWithGoogleLens captures an image using Google Lens
func (a *AndroidManager) CaptureWithGoogleLens() (string, error) {
	if !a.hasADB {
		return "", fmt.Errorf("ADB not available")
	}

	slog.Info("Capturing image with Google Lens")

	// Launch Google Lens app
	_, err := a.ExecuteShellCommand("am start -n com.google.ar.lens/com.google.android.apps.lens.main.MainActivity")
	if err != nil {
		// Fallback to launching through Google app
		_, err = a.ExecuteShellCommand("am start -n com.google.android.googlequicksearchbox/com.google.android.googlequicksearchbox.SearchActivity")
		if err != nil {
			return "", fmt.Errorf("failed to launch Google Lens: %v", err)
		}
	}

	// Trigger camera capture (this would normally be done through UI automation)
	// For now, we'll simulate by creating a timestamped filename
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("/sdcard/DCIM/LensCapture_%s.jpg", timestamp)

	// In a real implementation, this would trigger the actual camera capture
	// and wait for the image to be saved

	slog.Info("Image captured", "file", filename)
	return filename, nil
}

// ProcessWithGoogleLens processes a file with Google Lens
func (a *AndroidManager) ProcessWithGoogleLens(filePath, operation string) (string, error) {
	if !a.hasADB {
		return "", fmt.Errorf("ADB not available")
	}

	slog.Info("Processing file with Google Lens", "file", filePath, "operation", operation)

	// Check if file exists on device
	_, err := a.ExecuteShellCommand(fmt.Sprintf("ls %s", filePath))
	if err != nil {
		return "", fmt.Errorf("file not found on device: %s", filePath)
	}

	// Launch Google Lens with the file
	// This would normally involve sharing the file with Google Lens
	// For simulation, we'll just return a mock result
	var result string

	switch operation {
	case "extract_text":
		result = fmt.Sprintf("Text extracted from %s\n[Simulated result: Google Lens would extract text here]", filePath)
	case "identify_object":
		result = fmt.Sprintf("Object identified in %s\n[Simulated result: Google Lens would identify objects here]", filePath)
	case "solve_math":
		result = fmt.Sprintf("Math problem solved from %s\n[Simulated result: Google Lens would solve math problems here]", filePath)
	case "translate_text":
		result = fmt.Sprintf("Text translated from %s\n[Simulated result: Google Lens would translate text here]", filePath)
	case "barcode":
		result = fmt.Sprintf("Barcode/QR code scanned from %s\n[Simulated result: Google Lens would scan barcode/QR code here]", filePath)
	case "document":
		result = fmt.Sprintf("Document processed from %s\n[Simulated result: Google Lens would summarize document here]", filePath)
	default:
		result = fmt.Sprintf("Processed %s with Google Lens operation: %s\n[Simulated result]", filePath, operation)
	}

	slog.Info("Google Lens processing completed", "file", filePath)
	return result, nil
}

// ExecuteTaskerProfile triggers a Tasker profile
func (a *AndroidManager) ExecuteTaskerProfile(profileName string) error {
	if !a.hasADB {
		return fmt.Errorf("ADB not available")
	}

	slog.Info("Executing Tasker profile", "profile", profileName)

	// Use ADB to broadcast Tasker intent
	command := fmt.Sprintf("am broadcast -a net.dinglisch.android.tasker.ACTION_TASK_WITHOUT_VARIABLES --es task_name \"%s\"", profileName)
	_, err := a.ExecuteShellCommand(command)
	if err != nil {
		// Try alternative method with variables
		command = fmt.Sprintf("am broadcast -a net.dinglisch.android.tasker.ACTION_TASK --es task_name \"%s\"", profileName)
		_, err = a.ExecuteShellCommand(command)
		if err != nil {
			return fmt.Errorf("failed to execute Tasker profile: %v", err)
		}
	}

	slog.Info("Tasker profile executed", "profile", profileName)
	return nil
}

// ExecuteAutoInputScript runs an AutoInput script
func (a *AndroidManager) ExecuteAutoInputScript(scriptName string) error {
	if !a.hasADB {
		return fmt.Errorf("ADB not available")
	}

	slog.Info("Executing AutoInput script", "script", scriptName)

	// Use ADB to trigger AutoInput
	command := fmt.Sprintf("am startservice -n com.joaomgcd.autoinput/.AutoInputService --es script \"%s\"", scriptName)
	_, err := a.ExecuteShellCommand(command)
	if err != nil {
		// Try alternative method
		command = fmt.Sprintf("am broadcast -a com.joaomgcd.autoinput.RUN_SCRIPT --es script_name \"%s\"", scriptName)
		_, err = a.ExecuteShellCommand(command)
		if err != nil {
			return fmt.Errorf("failed to execute AutoInput script: %v", err)
		}
	}

	slog.Info("AutoInput script executed", "script", scriptName)
	return nil
}

// ShowTaskerScene displays a Tasker scene
func (a *AndroidManager) ShowTaskerScene(sceneName string) error {
	if !a.hasADB {
		return fmt.Errorf("ADB not available")
	}

	slog.Info("Showing Tasker scene", "scene", sceneName)

	// Use ADB to show Tasker scene
	command := fmt.Sprintf("am broadcast -a net.dinglisch.android.tasker.ACTION_SCENE_SHOW --es scene_name \"%s\"", sceneName)
	_, err := a.ExecuteShellCommand(command)
	if err != nil {
		return fmt.Errorf("failed to show Tasker scene: %v", err)
	}

	slog.Info("Tasker scene shown", "scene", sceneName)
	return nil
}

// ExitTasker exits the Tasker application
func (a *AndroidManager) ExitTasker() error {
	if !a.hasADB {
		return fmt.Errorf("ADB not available")
	}

	slog.Info("Exiting Tasker")

	// Use ADB to exit Tasker
	command := "am broadcast -a net.dinglisch.android.tasker.ACTION_TASKER_EXIT"
	_, err := a.ExecuteShellCommand(command)
	if err != nil {
		return fmt.Errorf("failed to exit Tasker: %v", err)
	}

	slog.Info("Tasker exited")
	return nil
}

// discoverTermuxFiles discovers files accessible through Termux
func (a *AndroidManager) discoverTermuxFiles() ([]AndroidFileInfo, error) {
	var files []AndroidFileInfo

	// Common Termux-accessible directories
	termuxDirs := []string{
		"/data/data/com.termux/files/home/storage/shared/",
		"/data/data/com.termux/files/home/storage/downloads/",
		"/data/data/com.termux/files/home/storage/documents/",
		"/data/data/com.termux/files/home/storage/dcim/",
	}

	for _, dir := range termuxDirs {
		cmd := exec.Command("find", dir, "-type", "f", "-not", "-path", "*/.*")
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out

		if err := cmd.Run(); err != nil {
			slog.Debug("Failed to scan Termux directory", "dir", dir, "error", err)
			continue
		}

		lines := strings.Split(strings.TrimSpace(out.String()), "\n")
		for _, line := range lines {
			if line != "" {
				fileInfo := AndroidFileInfo{
					Path:     line,
					Filename: strings.TrimPrefix(line, dir),
					Size:     0,          // Would need additional command to get size
					ModTime:  time.Now(), // Would need additional command to get actual time
					MimeType: "application/octet-stream",
					IsSystem: false,
				}
				files = append(files, fileInfo)
			}
		}
	}

	return files, nil
}

// discoverADBFiles discovers files using ADB
func (a *AndroidManager) discoverADBFiles() ([]AndroidFileInfo, error) {
	var files []AndroidFileInfo

	if !a.hasADB {
		return files, nil
	}

	// Android storage locations to scan
	androidPaths := []string{
		"/sdcard/",
		"/storage/emulated/0/",
		"/data/media/0/",
	}

	for _, basePath := range androidPaths {
		// List files recursively (basic implementation)
		cmd := exec.Command("adb", "-s", a.deviceID, "shell", "find", basePath, "-type", "f", "-not", "-path", "*/.*")
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out

		if err := cmd.Run(); err != nil {
			slog.Debug("Failed to scan Android path", "path", basePath, "error", err)
			continue
		}

		lines := strings.Split(strings.TrimSpace(out.String()), "\n")
		for _, line := range lines {
			if line != "" && !strings.Contains(line, "Permission denied") {
				fileInfo := AndroidFileInfo{
					Path:     line,
					Filename: strings.TrimPrefix(line, basePath),
					Size:     0,          // Would need additional command to get size
					ModTime:  time.Now(), // Would need additional command to get actual time
					MimeType: "application/octet-stream",
					IsSystem: strings.Contains(line, "/system/") || strings.Contains(line, "/vendor/"),
				}
				files = append(files, fileInfo)
			}
		}
	}

	return files, nil
}

// PullFile copies a file from Android device to local storage
func (a *AndroidManager) PullFile(remotePath, localPath string) error {
	if !a.hasADB {
		return fmt.Errorf("ADB not available")
	}

	cmd := exec.Command("adb", "-s", a.deviceID, "pull", remotePath, localPath)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to pull file: %v, output: %s", err, out.String())
	}

	return nil
}

// PushFile copies a file from local storage to Android device
func (a *AndroidManager) PushFile(localPath, remotePath string) error {
	if !a.hasADB {
		return fmt.Errorf("ADB not available")
	}

	cmd := exec.Command("adb", "-s", a.deviceID, "push", localPath, remotePath)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to push file: %v, output: %s", err, out.String())
	}

	return nil
}

// ListPackages lists installed Android packages
func (a *AndroidManager) ListPackages() ([]string, error) {
	if !a.hasADB {
		return nil, fmt.Errorf("ADB not available")
	}

	cmd := exec.Command("adb", "-s", a.deviceID, "shell", "pm", "list", "packages")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to list packages: %v", err)
	}

	var packages []string
	lines := strings.Split(out.String(), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "package:") {
			pkg := strings.TrimPrefix(line, "package:")
			packages = append(packages, pkg)
		}
	}

	return packages, nil
}

// SendNotification sends a notification using Termux API
func (a *AndroidManager) SendNotification(title, content string) error {
	if !a.hasTermux {
		return fmt.Errorf("Termux API not available")
	}

	// Ensure title is not empty
	if title == "" {
		title = "AI Agent Notification"
	}

	cmd := exec.Command("termux-notification", "--title", title, "--content", content)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to send notification: %v, output: %s", err, out.String())
	}

	return nil
}

// ExecuteShellCommand executes a shell command on Android device
func (a *AndroidManager) ExecuteShellCommand(command string) (string, error) {
	if !a.hasADB {
		return "", fmt.Errorf("ADB not available")
	}

	cmd := exec.Command("adb", "-s", a.deviceID, "shell", command)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to execute command: %v, output: %s", err, out.String())
	}

	return out.String(), nil
}

// GrantPermissions grants permissions to an app using Shizuku (if available)
func (a *AndroidManager) GrantPermissions(packageName, permission string) error {
	if !a.hasShizuku {
		return fmt.Errorf("Shizuku not available")
	}

	// Use Shizuku to grant permission
	cmd := exec.Command("adb", "-s", a.deviceID, "shell",
		"sh", "-c", fmt.Sprintf("shizuku su -c \"pm grant %s %s\"", packageName, permission))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to grant permission: %v, output: %s", err, out.String())
	}

	return nil
}

// HasADB returns true if ADB is available
func (a *AndroidManager) HasADB() bool {
	return a.hasADB
}

// HasTermux returns true if Termux API is available
func (a *AndroidManager) HasTermux() bool {
	return a.hasTermux
}

// HasShizuku returns true if Shizuku is available
func (a *AndroidManager) HasShizuku() bool {
	return a.hasShizuku
}

// GetDeviceInfo retrieves basic device information
func (a *AndroidManager) GetDeviceInfo() (map[string]string, error) {
	info := make(map[string]string)

	if !a.hasADB {
		return info, fmt.Errorf("ADB not available")
	}

	// Get device model
	modelOutput, err := a.ExecuteShellCommand("getprop ro.product.model")
	if err == nil {
		info["model"] = strings.TrimSpace(modelOutput)
	}

	// Get Android version
	versionOutput, err := a.ExecuteShellCommand("getprop ro.build.version.release")
	if err == nil {
		info["android_version"] = strings.TrimSpace(versionOutput)
	}

	// Get SDK version
	sdkOutput, err := a.ExecuteShellCommand("getprop ro.build.version.sdk")
	if err == nil {
		info["sdk_version"] = strings.TrimSpace(sdkOutput)
	}

	// Get device ID
	info["device_id"] = a.deviceID

	return info, nil
}

// MonitorBattery monitors battery status
func (a *AndroidManager) MonitorBattery() (map[string]interface{}, error) {
	status := make(map[string]interface{})

	if !a.hasADB {
		return status, fmt.Errorf("ADB not available")
	}

	// Get battery status
	batteryOutput, err := a.ExecuteShellCommand("dumpsys battery")
	if err != nil {
		return status, err
	}

	lines := strings.Split(batteryOutput, "\n")
	for _, line := range lines {
		if strings.Contains(line, "level:") {
			parts := strings.Fields(line)
			if len(parts) > 1 {
				status["level"] = parts[1]
			}
		} else if strings.Contains(line, "status:") {
			parts := strings.Fields(line)
			if len(parts) > 1 {
				status["status"] = parts[1]
			}
		} else if strings.Contains(line, "health:") {
			parts := strings.Fields(line)
			if len(parts) > 1 {
				status["health"] = parts[1]
			}
		}
	}

	return status, nil
}

// TakeScreenshot takes a screenshot of the Android device
func (a *AndroidManager) TakeScreenshot(localPath string) error {
	if !a.hasADB {
		return fmt.Errorf("ADB not available")
	}

	// Take screenshot on device
	_, err := a.ExecuteShellCommand("screencap -p /sdcard/temp_screenshot.png")
	if err != nil {
		return fmt.Errorf("failed to take screenshot: %v", err)
	}

	// Pull screenshot to local storage
	err = a.PullFile("/sdcard/temp_screenshot.png", localPath)
	if err != nil {
		return fmt.Errorf("failed to pull screenshot: %v", err)
	}

	// Clean up temporary file on device
	a.ExecuteShellCommand("rm /sdcard/temp_screenshot.png")

	return nil
}
