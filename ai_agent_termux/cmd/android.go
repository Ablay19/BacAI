package cmd

import (
	"fmt"
	"strings"

	"ai_agent_termux/android"
	"ai_agent_termux/config"

	"github.com/spf13/cobra"
)

// ValidateFilePath validates that a file path is acceptable
func ValidateFilePath(filePath string) error {
	// Check if file path is empty
	if filePath == "" {
		return fmt.Errorf("file path cannot be empty")
	}

	// Check if file path starts with a valid Android path prefix
	validPrefixes := []string{"/sdcard/", "/storage/", "/data/media/", "/data/data/com.termux/"}

	valid := false
	for _, prefix := range validPrefixes {
		if strings.HasPrefix(filePath, prefix) {
			valid = true
			break
		}
	}

	if !valid {
		return fmt.Errorf("invalid file path: %s. Path must start with /sdcard/, /storage/, /data/media/, or /data/data/com.termux/", filePath)
	}

	return nil
}

// androidCmd represents the android command
var androidCmd = &cobra.Command{
	Use:   "android [command]",
	Short: "Android device management and integration",
	Long: `Manage Android device integration with ADB, Termux, and Shizuku.
	Features:
	- File discovery across Android storage
	- Package management
	- Device information retrieval
	- Battery monitoring
	- Screenshot capture
	- Notification sending (Termux)
	- Permission management (Shizuku)
	- Google Lens integration
	- Tasker/AutoInput automation`,
}

func init() {
	// Add subcommands
	androidCmd.AddCommand(androidDiscoverCmd)
	androidCmd.AddCommand(androidInfoCmd)
	androidCmd.AddCommand(androidPackagesCmd)
	androidCmd.AddCommand(androidBatteryCmd)
	androidCmd.AddCommand(androidScreenshotCmd)
	androidCmd.AddCommand(androidNotifyCmd)
	androidCmd.AddCommand(androidPullCmd)
	androidCmd.AddCommand(androidPushCmd)
	androidCmd.AddCommand(androidLensCmd)
	androidCmd.AddCommand(androidTaskerCmd)
}

// androidDiscoverCmd represents the android discover command
var androidDiscoverCmd = &cobra.Command{
	Use:   "discover",
	Short: "Discover files on Android device",
	Long:  `Discover files across Android storage locations using ADB and Termux.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.LoadConfig()
		manager := android.NewAndroidManager(cfg)

		if !manager.HasADB() && !manager.HasTermux() {
			fmt.Println("❌ No Android tools available (ADB or Termux)")
			return
		}

		fmt.Println("🔍 Discovering Android files...")

		files, err := manager.DiscoverAndroidFiles()
		if err != nil {
			fmt.Printf("❌ Error discovering files: %v\n", err)
			return
		}

		fmt.Printf("✅ Discovered %d files\n", len(files))

		// Show sample files
		if len(files) > 0 {
			fmt.Println("\nSample files found:")
			max := len(files)
			if max > 10 {
				max = 10
			}

			for i := 0; i < max; i++ {
				fmt.Printf("  %s (%s)\n", files[i].Filename, files[i].Path)
			}

			if len(files) > 10 {
				fmt.Printf("  ... and %d more files\n", len(files)-10)
			}
		}
	},
}

// androidInfoCmd represents the android info command
var androidInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get Android device information",
	Long:  `Retrieve detailed information about connected Android device.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.LoadConfig()
		manager := android.NewAndroidManager(cfg)

		if !manager.HasADB() {
			fmt.Println("❌ ADB not available")
			return
		}

		fmt.Println("📱 Retrieving device information...")

		info, err := manager.GetDeviceInfo()
		if err != nil {
			fmt.Printf("❌ Error getting device info: %v\n", err)
			return
		}

		fmt.Println("Device Information:")
		fmt.Println("==================")
		for key, value := range info {
			fmt.Printf("%s: %s\n", strings.Title(key), value)
		}
	},
}

// androidPackagesCmd represents the android packages command
var androidPackagesCmd = &cobra.Command{
	Use:   "packages",
	Short: "List installed Android packages",
	Long:  `List all installed packages on the Android device.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.LoadConfig()
		manager := android.NewAndroidManager(cfg)

		if !manager.HasADB() {
			fmt.Println("❌ ADB not available")
			return
		}

		fmt.Println("📦 Retrieving installed packages...")

		packages, err := manager.ListPackages()
		if err != nil {
			fmt.Printf("❌ Error listing packages: %v\n", err)
			return
		}

		fmt.Printf("✅ Found %d installed packages\n", len(packages))

		// Show sample packages
		if len(packages) > 0 {
			fmt.Println("\nSample packages:")
			max := len(packages)
			if max > 20 {
				max = 20
			}

			for i := 0; i < max; i++ {
				fmt.Printf("  %s\n", packages[i])
			}

			if len(packages) > 20 {
				fmt.Printf("  ... and %d more packages\n", len(packages)-20)
			}
		}
	},
}

// androidBatteryCmd represents the android battery command
var androidBatteryCmd = &cobra.Command{
	Use:   "battery",
	Short: "Monitor Android battery status",
	Long:  `Monitor battery level, status, and health of the Android device.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.LoadConfig()
		manager := android.NewAndroidManager(cfg)

		if !manager.HasADB() {
			fmt.Println("❌ ADB not available")
			return
		}

		fmt.Println("🔋 Monitoring battery status...")

		status, err := manager.MonitorBattery()
		if err != nil {
			fmt.Printf("❌ Error monitoring battery: %v\n", err)
			return
		}

		fmt.Println("Battery Status:")
		fmt.Println("===============")
		for key, value := range status {
			fmt.Printf("%s: %v\n", strings.Title(key), value)
		}
	},
}

// androidScreenshotCmd represents the android screenshot command
var androidScreenshotCmd = &cobra.Command{
	Use:   "screenshot [local-path]",
	Short: "Take screenshot of Android device",
	Long:  `Capture and download a screenshot from the Android device.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		localPath := args[0]
		cfg := config.LoadConfig()
		manager := android.NewAndroidManager(cfg)

		if !manager.HasADB() {
			fmt.Println("❌ ADB not available")
			return
		}

		fmt.Printf("📸 Taking screenshot and saving to %s...\n", localPath)

		err := manager.TakeScreenshot(localPath)
		if err != nil {
			fmt.Printf("❌ Error taking screenshot: %v\n", err)
			return
		}

		fmt.Println("✅ Screenshot saved successfully!")
	},
}

// androidNotifyCmd represents the android notify command
var androidNotifyCmd = &cobra.Command{
	Use:   "notify [title] [content]",
	Short: "Send notification to Android device",
	Long:  `Send a notification to the Android device using Termux API.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		title := args[0]
		content := args[1]
		cfg := config.LoadConfig()
		manager := android.NewAndroidManager(cfg)

		if !manager.HasTermux() {
			fmt.Println("❌ Termux API not available")
			return
		}

		fmt.Printf("🔔 Sending notification: %s\n", title)

		err := manager.SendNotification(title, content)
		if err != nil {
			fmt.Printf("❌ Error sending notification: %v\n", err)
			return
		}

		fmt.Println("✅ Notification sent successfully!")
	},
}

// androidPullCmd represents the android pull command
var androidPullCmd = &cobra.Command{
	Use:   "pull [remote-path] [local-path]",
	Short: "Pull file from Android device",
	Long:  `Copy a file from the Android device to local storage using ADB.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		remotePath := args[0]
		localPath := args[1]
		cfg := config.LoadConfig()
		manager := android.NewAndroidManager(cfg)

		if !manager.HasADB() {
			fmt.Println("❌ ADB not available")
			return
		}

		fmt.Printf("📥 Pulling %s to %s...\n", remotePath, localPath)

		err := manager.PullFile(remotePath, localPath)
		if err != nil {
			fmt.Printf("❌ Error pulling file: %v\n", err)
			return
		}

		fmt.Println("✅ File pulled successfully!")
	},
}

// androidPushCmd represents the android push command
var androidPushCmd = &cobra.Command{
	Use:   "push [local-path] [remote-path]",
	Short: "Push file to Android device",
	Long:  `Copy a file from local storage to the Android device using ADB.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		localPath := args[0]
		remotePath := args[1]
		cfg := config.LoadConfig()
		manager := android.NewAndroidManager(cfg)

		if !manager.HasADB() {
			fmt.Println("❌ ADB not available")
			return
		}

		fmt.Printf("📤 Pushing %s to %s...\n", localPath, remotePath)

		err := manager.PushFile(localPath, remotePath)
		if err != nil {
			fmt.Printf("❌ Error pushing file: %v\n", err)
			return
		}

		fmt.Println("✅ File pushed successfully!")
	},
}

// androidLensCmd represents the android lens command
var androidLensCmd = &cobra.Command{
	Use:   "lens [operation] [file-path]",
	Short: "Process files with Google Lens",
	Long: `Process files with Google Lens for various operations:
	- extract_text: Extract text from images
	- identify_object: Identify objects in images
	- solve_math: Solve math problems from images
	- translate_text: Translate text from images
	- capture: Capture image using Google Lens
	- barcode: Scan barcode/QR code from images
	- document: Process document for summarization
	
Examples:
	ai_agent android lens capture
	ai_agent android lens extract_text /sdcard/Download/receipt.jpg
	ai_agent android lens solve_math /sdcard/DCIM/Camera/math_homework.jpg
	ai_agent android lens barcode /sdcard/Pictures/barcode.png`,
	Args: cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		operation := args[0]

		// Validate operation
		validOperations := map[string]bool{
			"extract_text":    true,
			"identify_object": true,
			"solve_math":      true,
			"translate_text":  true,
			"capture":         true,
			"barcode":         true,
			"document":        true,
		}

		if !validOperations[operation] {
			fmt.Printf("❌ Invalid operation: %s\n", operation)
			fmt.Println("Valid operations: extract_text, identify_object, solve_math, translate_text, capture, barcode, document")
			fmt.Println("\nExamples:")
			fmt.Println("  ai_agent android lens capture")
			fmt.Println("  ai_agent android lens extract_text /sdcard/Download/receipt.jpg")
			fmt.Println("  ai_agent android lens solve_math /sdcard/DCIM/Camera/math_homework.jpg")
			return
		}

		cfg := config.LoadConfig()
		manager := android.NewAndroidManager(cfg)

		if !manager.HasADB() {
			fmt.Println("❌ ADB not available")
			return
		}

		if operation == "capture" {
			fmt.Println("📸 Capturing image with Google Lens...")

			filePath, err := manager.CaptureWithGoogleLens()
			if err != nil {
				fmt.Printf("❌ Error capturing with Google Lens: %v\n", err)
				return
			}

			fmt.Printf("✅ Image captured: %s\n", filePath)
			return
		}

		if len(args) < 2 {
			fmt.Printf("❌ File path required for '%s' operation\n", operation)
			fmt.Println("\nExamples:")
			switch operation {
			case "extract_text":
				fmt.Println("  ai_agent android lens extract_text /sdcard/Download/receipt.jpg")
			case "identify_object":
				fmt.Println("  ai_agent android lens identify_object /sdcard/DCIM/Camera/object.jpg")
			case "solve_math":
				fmt.Println("  ai_agent android lens solve_math /sdcard/DCIM/Camera/math.jpg")
			case "translate_text":
				fmt.Println("  ai_agent android lens translate_text /sdcard/Download/foreign_text.jpg")
			case "barcode":
				fmt.Println("  ai_agent android lens barcode /sdcard/Pictures/barcode.png")
			case "document":
				fmt.Println("  ai_agent android lens document /sdcard/Documents/contract.pdf")
			}
			return
		}

		filePath := args[1]

		// Validate file path
		if err := ValidateFilePath(filePath); err != nil {
			fmt.Printf("❌ %v\n", err)
			return
		}

		fmt.Printf("🔍 Processing %s with Google Lens (%s)...\n", filePath, operation)

		result, err := manager.ProcessWithGoogleLens(filePath, operation)
		if err != nil {
			fmt.Printf("❌ Error processing with Google Lens: %v\n", err)
			return
		}

		fmt.Println("✅ Google Lens processing completed!")
		fmt.Println("Result:")
		fmt.Println(result)
	},
}

// androidTaskerCmd represents the android tasker command
var androidTaskerCmd = &cobra.Command{
	Use:   "tasker [subcommand] [arguments...]",
	Short: "Execute Tasker profiles and AutoInput scripts",
	Long: `Execute Tasker profiles and AutoInput scripts on Android device:
	- profile [profile-name]: Execute a Tasker profile
	- autoinput [script-name]: Execute an AutoInput script
	- scene [scene-name]: Show a Tasker scene
	- exit: Exit Tasker
	
Examples:
	ai_agent android tasker profile "Silent Mode"
	ai_agent android tasker autoinput "Login Script"
	ai_agent android tasker scene "Quick Menu"`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		subcommand := args[0]

		// Validate subcommand
		validSubcommands := map[string]bool{
			"profile":   true,
			"autoinput": true,
			"scene":     true,
			"exit":      true,
		}

		if !validSubcommands[subcommand] {
			fmt.Printf("❌ Invalid subcommand: %s\n", subcommand)
			fmt.Println("Valid subcommands: profile, autoinput, scene, exit")
			fmt.Println("\nExamples:")
			fmt.Println("  ai_agent android tasker profile \"Silent Mode\"")
			fmt.Println("  ai_agent android tasker autoinput \"Login Script\"")
			fmt.Println("  ai_agent android tasker scene \"Quick Menu\"")
			fmt.Println("  ai_agent android tasker exit")
			return
		}

		cfg := config.LoadConfig()
		manager := android.NewAndroidManager(cfg)

		if !manager.HasADB() {
			fmt.Println("❌ ADB not available")
			return
		}

		switch subcommand {
		case "profile":
			if len(args) < 2 {
				fmt.Println("❌ Profile name required")
				fmt.Println("\nExample: ai_agent android tasker profile \"Silent Mode\"")
				return
			}

			name := args[1]
			fmt.Printf("▶️ Executing Tasker profile: %s\n", name)

			err := manager.ExecuteTaskerProfile(name)
			if err != nil {
				fmt.Printf("❌ Error executing Tasker profile: %v\n", err)
				return
			}

			fmt.Println("✅ Tasker profile executed successfully!")

		case "autoinput":
			if len(args) < 2 {
				fmt.Println("❌ Script name required")
				fmt.Println("\nExample: ai_agent android tasker autoinput \"Login Script\"")
				return
			}

			name := args[1]
			fmt.Printf("▶️ Executing AutoInput script: %s\n", name)

			err := manager.ExecuteAutoInputScript(name)
			if err != nil {
				fmt.Printf("❌ Error executing AutoInput script: %v\n", err)
				return
			}

			fmt.Println("✅ AutoInput script executed successfully!")

		case "scene":
			if len(args) < 2 {
				fmt.Println("❌ Scene name required")
				fmt.Println("\nExample: ai_agent android tasker scene \"Quick Menu\"")
				return
			}

			name := args[1]
			fmt.Printf("🎨 Showing Tasker scene: %s\n", name)

			// Execute scene via Tasker intent
			command := fmt.Sprintf("am broadcast -a net.dinglisch.android.tasker.ACTION_SCENE_SHOW --es scene_name \"%s\"", name)
			_, err := manager.ExecuteShellCommand(command)
			if err != nil {
				fmt.Printf("❌ Error showing Tasker scene: %v\n", err)
				return
			}

			fmt.Println("✅ Tasker scene shown successfully!")

		case "exit":
			fmt.Println("⏹️ Exiting Tasker...")

			// Exit Tasker via intent
			command := "am broadcast -a net.dinglisch.android.tasker.ACTION_TASKER_EXIT"
			_, err := manager.ExecuteShellCommand(command)
			if err != nil {
				fmt.Printf("❌ Error exiting Tasker: %v\n", err)
				return
			}

			fmt.Println("✅ Tasker exited successfully!")
		}
	},
}
