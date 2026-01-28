package cmd

import (
	"fmt"
	"strings"

	"ai_agent_termux/config"
	"ai_agent_termux/realtime"
	"ai_agent_termux/tools"

	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
)

// toolsCmd represents the tools command
var toolsCmd = &cobra.Command{
	Use:   "tools [command]",
	Short: "Manage system tools and utilities",
	Long: `Manage and utilize system tools including aichat, ffmpeg, ffprobe, and 30+ useful packages.
	Features:
	- Tool detection and version checking
	- Media file analysis with ffprobe
	- Media conversion with ffmpeg
	- aichat model management and querying
	- Package discovery and installation suggestions`,
}

// realtimeCmd represents the realtime command
var realtimeCmd = &cobra.Command{
	Use:   "realtime [command]",
	Short: "Real-time logging and notification system",
	Long: `Monitor and manage real-time logging with notifications via Termux and ntfy.sh.
	Features:
	- Live log viewing and filtering
	- System notification sending
	- Log statistics and analysis
	- File change watching with alerts`,
}

func init() {
	// Add tools subcommands
	toolsCmd.AddCommand(toolsListCmd)
	toolsCmd.AddCommand(toolsCheckCmd)
	toolsCmd.AddCommand(toolsAichatCmd)
	toolsCmd.AddCommand(toolsMediaCmd)
	toolsCmd.AddCommand(toolsPackagesCmd)

	// Add realtime subcommands
	realtimeCmd.AddCommand(realtimeLogCmd)
	realtimeCmd.AddCommand(realtimeNotifyCmd)
	realtimeCmd.AddCommand(realtimeWatchCmd)
	realtimeCmd.AddCommand(realtimeStatsCmd)

	// Add to root
	rootCmd.AddCommand(toolsCmd)
	rootCmd.AddCommand(realtimeCmd)
}

// toolsListCmd represents the tools list command
var toolsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available system tools",
	Long:  `Discover and list all available system tools with version information.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.LoadConfig()
		toolManager := tools.NewToolManager(cfg)

		fmt.Println("🔍 Detecting system tools...")
		allTools := toolManager.GetAllTools()

		if len(allTools) == 0 {
			fmt.Println("❌ No tools detected")
			return
		}

		fmt.Printf("✅ Found %d tools:\n", len(allTools))
		fmt.Println()

		// Group tools by category (simplified)
		categories := make(map[string][]string)
		for name, tool := range allTools {
			if tool.Availability {
				category := "General"
				// Very basic categorization
				if strings.Contains(name, "ffmpeg") || strings.Contains(name, "ffprobe") {
					category = "Media"
				} else if strings.Contains(name, "aichat") {
					category = "AI"
				} else if strings.Contains(name, "termux") {
					category = "Android"
				}

				categories[category] = append(categories[category], fmt.Sprintf("  %s (%s)", name, tool.Version))
			}
		}

		for category, tools := range categories {
			fmt.Printf("%s Tools:\n", category)
			for _, tool := range tools {
				fmt.Println(tool)
			}
			fmt.Println()
		}
	},
}

// toolsCheckCmd represents the tools check command
var toolsCheckCmd = &cobra.Command{
	Use:   "check [tool-name]",
	Short: "Check if a specific tool is available",
	Long:  `Verify if a specific tool is installed and accessible.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		toolName := args[0]
		cfg := config.LoadConfig()
		toolManager := tools.NewToolManager(cfg)

		fmt.Printf("🔍 Checking tool: %s\n", toolName)

		if toolManager.IsToolAvailable(toolName) {
			toolInfo, err := toolManager.GetToolInfo(toolName)
			if err != nil {
				fmt.Printf("✅ %s is available (but info unavailable: %v)\n", toolName, err)
				return
			}

			fmt.Printf("✅ %s is available:\n", toolName)
			fmt.Printf("  Version: %s\n", toolInfo.Version)
			fmt.Printf("  Path: %s\n", toolInfo.Path)
		} else {
			fmt.Printf("❌ %s is not available\n", toolName)
		}
	},
}

// toolsAichatCmd represents the tools aichat command
var toolsAichatCmd = &cobra.Command{
	Use:   "aichat [subcommand]",
	Short: "Manage aichat integration",
	Long:  `Work with aichat for AI-powered conversations and queries.`,
}

func init() {
	toolsAichatCmd.AddCommand(toolsAichatListCmd)
	toolsAichatCmd.AddCommand(toolsAichatQueryCmd)
}

// toolsAichatListCmd represents the tools aichat list command
var toolsAichatListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available aichat models",
	Long:  `Show all AI models available through aichat.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.LoadConfig()
		toolManager := tools.NewToolManager(cfg)

		if !toolManager.IsToolAvailable("aichat") {
			fmt.Println("❌ aichat is not installed or not in PATH")
			fmt.Println("💡 Install with: cargo install aichat")
			return
		}

		fmt.Println("🤖 Getting available aichat models...")
		models, err := toolManager.GetAichatModels()
		if err != nil {
			fmt.Printf("❌ Error getting models: %v\n", err)
			return
		}

		fmt.Printf("✅ Found %d models:\n", len(models))
		for _, model := range models {
			fmt.Printf("  • %s\n", model)
		}
	},
}

// toolsAichatQueryCmd represents the tools aichat query command
var toolsAichatQueryCmd = &cobra.Command{
	Use:   "query [prompt]",
	Short: "Execute a query using aichat",
	Long:  `Send a prompt to aichat and get AI-generated response.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		prompt := args[0]
		model, _ := cmd.Flags().GetString("model")

		cfg := config.LoadConfig()
		toolManager := tools.NewToolManager(cfg)

		if !toolManager.IsToolAvailable("aichat") {
			fmt.Println("❌ aichat is not installed or not in PATH")
			fmt.Println("💡 Install with: cargo install aichat")
			return
		}

		fmt.Printf("🤖 Executing query: %s\n", prompt)
		if model != "" {
			fmt.Printf("   Using model: %s\n", model)
		}

		response, err := toolManager.ExecuteAichatQuery(prompt, model)
		if err != nil {
			fmt.Printf("❌ Error executing query: %v\n", err)
			return
		}

		fmt.Println("✅ Response:")
		fmt.Println(response)
	},
}

func init() {
	toolsAichatQueryCmd.Flags().String("model", "", "Specific AI model to use")
}

// toolsMediaCmd represents the tools media command
var toolsMediaCmd = &cobra.Command{
	Use:   "media [subcommand]",
	Short: "Media processing tools (ffmpeg/ffprobe)",
	Long:  `Analyze and process media files using ffmpeg and ffprobe.`,
}

func init() {
	toolsMediaCmd.AddCommand(toolsMediaInfoCmd)
	toolsMediaCmd.AddCommand(toolsMediaConvertCmd)
}

// toolsMediaInfoCmd represents the tools media info command
var toolsMediaInfoCmd = &cobra.Command{
	Use:   "info [file-path]",
	Short: "Get media file information with ffprobe",
	Long:  `Extract detailed information from media files using ffprobe.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]
		cfg := config.LoadConfig()
		toolManager := tools.NewToolManager(cfg)

		if !toolManager.IsToolAvailable("ffprobe") {
			fmt.Println("❌ ffprobe is not installed or not in PATH")
			fmt.Println("💡 Install with: pkg install ffmpeg")
			return
		}

		fmt.Printf("🔍 Analyzing media file: %s\n", filePath)

		mediaInfo, err := toolManager.GetMediaFileInfo(filePath)
		if err != nil {
			fmt.Printf("❌ Error analyzing file: %v\n", err)
			return
		}

		fmt.Println("✅ Media Information:")
		fmt.Printf("  Format: %s\n", mediaInfo.Format)
		fmt.Printf("  Duration: %.2f seconds\n", mediaInfo.Duration)
		fmt.Printf("  Bit Rate: %d bps\n", mediaInfo.BitRate)
		fmt.Printf("  Size: %d bytes\n", mediaInfo.Size)
		fmt.Printf("  Streams: %d\n", len(mediaInfo.Streams))

		for i, stream := range mediaInfo.Streams {
			fmt.Printf("  Stream %d:\n", i+1)
			fmt.Printf("    Type: %s\n", stream.CodecType)
			fmt.Printf("    Codec: %s\n", stream.CodecName)

			if stream.CodecType == "video" {
				fmt.Printf("    Resolution: %dx%d\n", stream.Width, stream.Height)
			} else if stream.CodecType == "audio" {
				fmt.Printf("    Sample Rate: %d Hz\n", stream.SampleRate)
				fmt.Printf("    Channels: %d\n", stream.Channels)
			}
		}
	},
}

// toolsMediaConvertCmd represents the tools media convert command
var toolsMediaConvertCmd = &cobra.Command{
	Use:   "convert [input] [output] [format]",
	Short: "Convert media files with ffmpeg",
	Long:  `Convert media files to different formats using ffmpeg.`,
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		inputFile := args[0]
		outputFile := args[1]
		format := args[2]

		cfg := config.LoadConfig()
		toolManager := tools.NewToolManager(cfg)

		if !toolManager.IsToolAvailable("ffmpeg") {
			fmt.Println("❌ ffmpeg is not installed or not in PATH")
			fmt.Println("💡 Install with: pkg install ffmpeg")
			return
		}

		fmt.Printf("🔄 Converting %s to %s format (%s)\n", inputFile, format, outputFile)

		// Additional options can be passed via flags
		options, _ := cmd.Flags().GetStringSlice("options")

		err := toolManager.ConvertMedia(inputFile, outputFile, format, options...)
		if err != nil {
			fmt.Printf("❌ Error converting file: %v\n", err)
			return
		}

		fmt.Println("✅ Conversion completed successfully!")
	},
}

func init() {
	toolsMediaConvertCmd.Flags().StringSlice("options", []string{}, "Additional ffmpeg options (e.g., --options='-crf 23,-preset fast')")
}

// toolsPackagesCmd represents the tools packages command
var toolsPackagesCmd = &cobra.Command{
	Use:   "packages",
	Short: "Discover useful packages and tools",
	Long:  `Discover 30+ useful packages that can enhance AI Agent capabilities.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.LoadConfig()
		toolManager := tools.NewToolManager(cfg)

		fmt.Println("📦 Discovering useful packages...")
		packages := toolManager.GetUsefulPackages()

		fmt.Printf("✅ Found %d packages:\n", len(packages))
		fmt.Println()

		// Group by priority and category
		highPriority := []tools.UsefulPackage{}
		mediumPriority := []tools.UsefulPackage{}
		lowPriority := []tools.UsefulPackage{}

		for _, pkg := range packages {
			switch {
			case pkg.Priority >= 9:
				highPriority = append(highPriority, pkg)
			case pkg.Priority >= 7:
				mediumPriority = append(mediumPriority, pkg)
			default:
				lowPriority = append(lowPriority, pkg)
			}
		}

		fmt.Println("⭐ High Priority Packages (Essential):")
		for _, pkg := range highPriority {
			status := "❌ Not Installed"
			if pkg.Installed {
				status = "✅ Installed"
			}
			fmt.Printf("  %s - %s\n", pkg.Name, pkg.Description)
			fmt.Printf("    %s | Install: %s\n", status, pkg.InstallCmd)
			fmt.Println()
		}

		fmt.Println("⭐ Medium Priority Packages (Recommended):")
		for _, pkg := range mediumPriority[:10] { // Show top 10
			status := "❌ Not Installed"
			if pkg.Installed {
				status = "✅ Installed"
			}
			fmt.Printf("  %s - %s\n", pkg.Name, pkg.Description)
			fmt.Printf("    %s | Install: %s\n", status, pkg.InstallCmd)
			fmt.Println()
		}

		fmt.Printf("📦 Total packages: %d (%d installed)\n", len(packages),
			len(highPriority)+len(mediumPriority)+len(lowPriority)-len(filterNotInstalled(packages)))
	},
}

func filterNotInstalled(packages []tools.UsefulPackage) []tools.UsefulPackage {
	var notInstalled []tools.UsefulPackage
	for _, pkg := range packages {
		if !pkg.Installed {
			notInstalled = append(notInstalled, pkg)
		}
	}
	return notInstalled
}

// realtimeLogCmd represents the realtime log command
var realtimeLogCmd = &cobra.Command{
	Use:   "log [subcommand]",
	Short: "View and manage real-time logs",
	Long:  `View, search, and analyze real-time application logs.`,
}

func init() {
	realtimeLogCmd.AddCommand(realtimeLogViewCmd)
	realtimeLogCmd.AddCommand(realtimeLogSearchCmd)
}

// realtimeLogViewCmd represents the realtime log view command
var realtimeLogViewCmd = &cobra.Command{
	Use:   "view",
	Short: "View recent log entries",
	Long:  `Display recent log entries with filtering options.`,
	Run: func(cmd *cobra.Command, args []string) {
		count, _ := cmd.Flags().GetInt("count")
		level, _ := cmd.Flags().GetString("level")

		cfg := config.LoadConfig()
		logManager := realtime.NewLogManager(cfg)
		defer logManager.Close()

		fmt.Printf("📋 Recent Logs (last %d entries):\n", count)

		var logs []realtime.LogEntry
		if level != "" {
			logs = logManager.SearchLogs("", level)
		} else {
			logs = logManager.GetRecentLogs(count)
		}

		if len(logs) == 0 {
			fmt.Println("  No logs found")
			return
		}

		for _, log := range logs {
			fmt.Printf("[%s] %s: %s (%s)\n",
				log.Timestamp.Format("15:04:05"),
				log.Level,
				log.Message,
				log.Source)
		}
	},
}

func init() {
	realtimeLogViewCmd.Flags().Int("count", 20, "Number of recent logs to display")
	realtimeLogViewCmd.Flags().String("level", "", "Filter by log level (ERROR, WARN, INFO, DEBUG)")
}

// realtimeLogSearchCmd represents the realtime log search command
var realtimeLogSearchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search logs for specific terms",
	Long:  `Search through log entries for specific keywords or phrases.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		query := args[0]
		level, _ := cmd.Flags().GetString("level")

		cfg := config.LoadConfig()
		logManager := realtime.NewLogManager(cfg)
		defer logManager.Close()

		fmt.Printf("🔍 Searching logs for: '%s'\n", query)
		if level != "" {
			fmt.Printf("  Filtered by level: %s\n", level)
		}

		var levels []string
		if level != "" {
			levels = []string{level}
		}

		logs := logManager.SearchLogs(query, levels...)

		if len(logs) == 0 {
			fmt.Println("  No matching logs found")
			return
		}

		fmt.Printf("✅ Found %d matching entries:\n", len(logs))
		for _, log := range logs {
			fmt.Printf("[%s] %s: %s (%s)\n",
				log.Timestamp.Format("15:04:05"),
				log.Level,
				log.Message,
				log.Source)
		}
	},
}

func init() {
	realtimeLogSearchCmd.Flags().String("level", "", "Filter by log level")
}

// realtimeNotifyCmd represents the realtime notify command
var realtimeNotifyCmd = &cobra.Command{
	Use:   "notify [title] [message]",
	Short: "Send system notification",
	Long:  `Send notification via Termux and/or ntfy.sh services.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		title := args[0]
		message := args[1]
		priority, _ := cmd.Flags().GetString("priority")

		cfg := config.LoadConfig()
		notificationManager := realtime.NewNotificationManager(cfg)

		fmt.Printf("🔔 Sending notification: %s\n", title)
		if priority != "" {
			fmt.Printf("   Priority: %s\n", priority)
		}

		err := notificationManager.SendNotification(title, message, priority)
		if err != nil {
			fmt.Printf("❌ Error sending notification: %v\n", err)
			return
		}

		fmt.Println("✅ Notification sent successfully!")
	},
}

func init() {
	realtimeNotifyCmd.Flags().String("priority", "default", "Notification priority (low, default, high, urgent)")
}

// realtimeWatchCmd represents the realtime watch command
var realtimeWatchCmd = &cobra.Command{
	Use:   "watch [paths...]",
	Short: "Watch files/directories for changes",
	Long:  `Monitor files and directories for changes with notifications.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		paths := args
		interval, _ := cmd.Flags().GetInt("interval")

		cfg := config.LoadConfig()
		notificationManager := realtime.NewNotificationManager(cfg)

		fmt.Printf("👀 Watching %d paths for changes (checking every %d seconds):\n", len(paths), interval)
		for _, path := range paths {
			fmt.Printf("  • %s\n", path)
		}

		// This would normally set up file watching
		// For now, we'll simulate
		err := notificationManager.WatchFileChanges(paths, func(path, event string) {
			slog.Info("File change detected", "path", path, "event", event)
			notificationManager.SendNotification("File Change",
				fmt.Sprintf("File changed: %s (%s)", path, event), "default")
		})

		if err != nil {
			fmt.Printf("❌ Error setting up file watching: %v\n", err)
			return
		}

		fmt.Println("✅ File watching started (press Ctrl+C to stop)")

		// Keep running until interrupted
		select {}
	},
}

func init() {
	realtimeWatchCmd.Flags().Int("interval", 5, "Check interval in seconds")
}

// realtimeStatsCmd represents the realtime stats command
var realtimeStatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show log statistics",
	Long:  `Display statistics about application logs.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.LoadConfig()
		logManager := realtime.NewLogManager(cfg)
		defer logManager.Close()

		stats := logManager.GetLogStatistics()

		fmt.Println("📊 Log Statistics:")
		fmt.Printf("  Total Entries: %d\n", stats["total_entries"])
		fmt.Printf("  Errors: %d\n", stats["error_count"])
		fmt.Printf("  Warnings: %d\n", stats["warn_count"])
		fmt.Printf("  Info Messages: %d\n", stats["info_count"])
		fmt.Printf("  Debug Messages: %d\n", stats["debug_count"])
	},
}
