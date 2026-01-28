package notifications

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"golang.org/x/exp/slog"
)

// NotifySystem handles notifications via ntfy.sh and Termux API
type NotifySystem struct {
	ntfyTopic   string
	ntfyServer  string
	httpClient  *http.Client
	termuxAPI   bool
	tags        []string
	priority    string
	clickAction string
}

// NotificationPriority represents notification priority levels
type NotificationPriority string

const (
	PriorityMax     NotificationPriority = "max"
	PriorityHigh    NotificationPriority = "high"
	PriorityDefault NotificationPriority = "default"
	PriorityLow     NotificationPriority = "low"
	PriorityMin     NotificationPriority = "min"
)

// NotificationMessage represents a notification message
type NotificationMessage struct {
	Topic    string               `json:"topic,omitempty"`
	Title    string               `json:"title,omitempty"`
	Message  string               `json:"message,omitempty"`
	Priority string               `json:"priority,omitempty"`
	Tags     []string             `json:"tags,omitempty"`
	Click    string               `json:"click,omitempty"`
	Attach   string               `json:"attach,omitempty"`
	Filename string               `json:"filename,omitempty"`
	Delay    string               `json:"delay,omitempty"`
	Email    string               `json:"email,omitempty"`
	Actions  []NotificationAction `json:"actions,omitempty"`
}

// NotificationAction represents an action button in a notification
type NotificationAction struct {
	Action string `json:"action"`
	Label  string `json:"label"`
	URL    string `json:"url,omitempty"`
	Clear  bool   `json:"clear,omitempty"`
}

// NotifyConfig contains configuration for the notification system
type NotifyConfig struct {
	NtfyTopic   string               `json:"ntfy_topic"`
	NtfyServer  string               `json:"ntfy_server"`
	TermuxAPI   bool                 `json:"termux_api"`
	DefaultTags []string             `json:"default_tags"`
	Priority    NotificationPriority `json:"priority"`
	ClickAction string               `json:"click_action"`
	Timeout     time.Duration        `json:"timeout"`
}

// DefaultNotifyConfig returns default configuration for notifications
func DefaultNotifyConfig() *NotifyConfig {
	return &NotifyConfig{
		NtfyServer:  "https://ntfy.sh",
		TermuxAPI:   true, // Enable Termux API by default in Termux environment
		DefaultTags: []string{"robot", "information"},
		Priority:    PriorityDefault,
		Timeout:     30 * time.Second,
	}
}

// NewNotifySystem creates a new notification system instance
func NewNotifySystem(config *NotifyConfig) *NotifySystem {
	if config == nil {
		config = DefaultNotifyConfig()
	}

	// Check if running in Termux environment
	isTermux := os.Getenv("TERMUX_VERSION") != "" || isTermuxInstalled()
	if !isTermux {
		config.TermuxAPI = false
		slog.Info("Termux API not available, disabled")
	}

	return &NotifySystem{
		ntfyTopic:   config.NtfyTopic,
		ntfyServer:  config.NtfyServer,
		httpClient:  &http.Client{Timeout: config.Timeout},
		termuxAPI:   config.TermuxAPI,
		tags:        config.DefaultTags,
		priority:    string(config.Priority),
		clickAction: config.ClickAction,
	}
}

// Send sends a notification with the specified title and message
func (ns *NotifySystem) Send(title, message string) error {
	return ns.SendWithPriority(title, message, PriorityDefault)
}

// SendWithPriority sends a notification with the specified priority
func (ns *NotifySystem) SendWithPriority(title, message string, priority NotificationPriority) error {
	return ns.sendNotification(title, message, priority, []string{}, "")
}

// SendSuccess sends a success notification
func (ns *NotifySystem) SendSuccess(title, message string) error {
	tags := append(ns.tags, "white_check_mark", "green_circle")
	return ns.sendNotification(title, message, PriorityDefault, tags, "")
}

// SendError sends an error notification
func (ns *NotifySystem) SendError(title, message string) error {
	tags := append(ns.tags, "warning", "red_circle")
	return ns.sendNotification(title, message, PriorityHigh, tags, "")
}

// SendInfo sends an info notification
func (ns *NotifySystem) SendInfo(title, message string) error {
	tags := append(ns.tags, "information_source", "blue_circle")
	return ns.sendNotification(title, message, PriorityDefault, tags, "")
}

// SendWarning sends a warning notification
func (ns *NotifySystem) SendWarning(title, message string) error {
	tags := append(ns.tags, "warning", "yellow_circle")
	return ns.sendNotification(title, message, PriorityHigh, tags, "")
}

// SendWithAttachment sends a notification with an attachment
func (ns *NotifySystem) SendWithAttachment(title, message, attachmentPath, filename string) error {
	tags := append(ns.tags, "attachment")
	return ns.sendNotification(title, message, PriorityDefault, tags, attachmentPath)
}

// SendWithActions sends a notification with action buttons
func (ns *NotifySystem) SendWithActions(title, message string, actions []NotificationAction) error {
	return ns.sendNotificationWithActions(title, message, PriorityDefault, ns.tags, "", actions)
}

// sendNotification sends the actual notification using available methods
func (ns *NotifySystem) sendNotification(title, message string, priority NotificationPriority, tags []string, attachment string) error {
	var errors []string

	// Try ntfy.sh if configured
	if ns.ntfyTopic != "" {
		if err := ns.sendNtfyNotification(title, message, priority, tags, attachment); err != nil {
			slog.Warn("Failed to send ntfy notification", "error", err)
			errors = append(errors, fmt.Sprintf("ntfy: %v", err))
		}
	}

	// Try Termux API if available
	if ns.termuxAPI {
		if err := ns.sendTermuxNotification(title, message, priority, tags, attachment); err != nil {
			slog.Warn("Failed to send Termux notification", "error", err)
			errors = append(errors, fmt.Sprintf("termux: %v", err))
		}
	}

	// If all methods failed, return error
	if len(errors) > 0 && len(errors) == 2 { // Both ntfy and Termux failed
		return fmt.Errorf("all notification methods failed: %s", strings.Join(errors, "; "))
	}

	slog.Debug("Notification sent successfully", "title", title, "methods", len(errors)+1)
	return nil
}

// sendNtfyNotification sends notification via ntfy.sh
func (ns *NotifySystem) sendNtfyNotification(title, message string, priority NotificationPriority, tags []string, attachment string) error {
	notification := NotificationMessage{
		Topic:    ns.ntfyTopic,
		Title:    title,
		Message:  message,
		Priority: string(priority),
		Tags:     tags,
	}

	if ns.clickAction != "" {
		notification.Click = ns.clickAction
	}

	if attachment != "" {
		notification.Attach = attachment
		notification.Filename = attachment
	}

	jsonData, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %v", err)
	}

	url := fmt.Sprintf("%s/%s", ns.ntfyServer, ns.ntfyTopic)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := ns.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("ntfy server returned status %d: %s", resp.StatusCode, string(body))
	}

	slog.Debug("ntfy notification sent successfully", "topic", ns.ntfyTopic)
	return nil
}

// sendTermuxNotification sends notification via Termux API
func (ns *NotifySystem) sendTermuxNotification(title, message string, priority NotificationPriority, tags []string, attachment string) error {
	// Termux notification command: termux-notification --title "Title" --content "Content"
	cmd := exec.Command("termux-notification", "--title", title, "--content", message)

	// Add priority if specified
	if priority != "" && priority != PriorityDefault {
		cmd.Args = append(cmd.Args, "--priority", string(priority))
	}

	// Add icon based on first tag if available
	if len(tags) > 0 {
		cmd.Args = append(cmd.Args, "--icon", tags[0])
	}

	// Add sound for high priority notifications
	if priority == PriorityHigh || priority == PriorityMax {
		cmd.Args = append(cmd.Args, "--sound")
	}

	// Add action if configured
	if ns.clickAction != "" {
		cmd.Args = append(cmd.Args, "--action", ns.clickAction)
	}

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to execute termux-notification: %v", err)
	}

	slog.Debug("Termux notification sent successfully")
	return nil
}

// sendNotificationWithActions sends notification with action buttons
func (ns *NotifySystem) sendNotificationWithActions(title, message string, priority NotificationPriority, tags []string, attachment string, actions []NotificationAction) error {
	var errors []string

	// Try ntfy.sh if configured (ntfy supports actions)
	if ns.ntfyTopic != "" {
		if err := ns.sendNtfyNotificationWithActions(title, message, priority, tags, attachment, actions); err != nil {
			slog.Warn("Failed to send ntfy notification with actions", "error", err)
			errors = append(errors, fmt.Sprintf("ntfy: %v", err))
		}
	}

	// For Termux, actions are not easily supported, so send a simple notification
	if ns.termuxAPI {
		if err := ns.sendTermuxNotification(title, message, priority, tags, attachment); err != nil {
			slog.Warn("Failed to send Termux notification", "error", err)
			errors = append(errors, fmt.Sprintf("termux: %v", err))
		}
	}

	if len(errors) > 0 && len(errors) == 2 {
		return fmt.Errorf("all notification methods failed: %s", strings.Join(errors, "; "))
	}

	return nil
}

// sendNtfyNotificationWithActions sends ntfy notification with actions
func (ns *NotifySystem) sendNtfyNotificationWithActions(title, message string, priority NotificationPriority, tags []string, attachment string, actions []NotificationAction) error {
	notification := NotificationMessage{
		Topic:    ns.ntfyTopic,
		Title:    title,
		Message:  message,
		Priority: string(priority),
		Tags:     tags,
		Actions:  actions,
	}

	if ns.clickAction != "" {
		notification.Click = ns.clickAction
	}

	if attachment != "" {
		notification.Attach = attachment
		notification.Filename = attachment
	}

	jsonData, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %v", err)
	}

	url := fmt.Sprintf("%s/%s", ns.ntfyServer, ns.ntfyTopic)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := ns.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("ntfy server returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// isTermuxInstalled checks if Termux API is installed
func isTermuxInstalled() bool {
	_, err := exec.LookPath("termux-notification")
	return err == nil
}

// TestNotifications sends test notifications to verify the system is working
func (ns *NotifySystem) TestNotifications(ctx context.Context) error {
	testTitle := "Google Lens Test"
	testMessage := fmt.Sprintf("Test notification sent at %s", time.Now().Format("15:04:05"))

	// Test basic notification
	if err := ns.SendInfo(testTitle, testMessage); err != nil {
		return fmt.Errorf("failed to send test notification: %v", err)
	}

	// Test notification with attachment (if we can create a test file)
	testFile := "/tmp/test_notification.txt"
	if err := os.WriteFile(testFile, []byte(testMessage), 0644); err == nil {
		if err := ns.SendWithAttachment("Test with Attachment", "This is a test with attachment", testFile, "test.txt"); err != nil {
			slog.Warn("Failed to send test notification with attachment", "error", err)
		}
		os.Remove(testFile)
	}

	// Test notification with actions
	actions := []NotificationAction{
		{Action: "view", Label: "View Result", URL: "https://example.com"},
		{Action: "dismiss", Label: "Dismiss", Clear: true},
	}
	if err := ns.SendWithActions("Test with Actions", "This is a test with action buttons", actions); err != nil {
		slog.Warn("Failed to send test notification with actions", "error", err)
	}

	slog.Info("Test notifications sent successfully")
	return nil
}

// GetCapabilities returns the capabilities of the notification system
func (ns *NotifySystem) GetCapabilities() map[string]interface{} {
	capabilities := make(map[string]interface{})

	capabilities["ntfy_enabled"] = ns.ntfyTopic != ""
	capabilities["termux_api_enabled"] = ns.termuxAPI
	capabilities["supports_attachments"] = ns.ntfyTopic != "" || ns.termuxAPI
	capabilities["supports_actions"] = ns.ntfyTopic != ""
	capabilities["default_priority"] = ns.priority
	capabilities["default_tags"] = ns.tags

	return capabilities
}
