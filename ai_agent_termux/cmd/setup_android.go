package cmd

import (
	"fmt"
	"time"

	"ai_agent_termux/android"
	"github.com/spf13/cobra"
)

var setupAndroidCmd = &cobra.Command{
	Use:   "setup-android",
	Short: "Configure G.I.D.A for Termux/Android",
	Long:  `Sets up widgets, notifications, and background scheduling for autonomous operation.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🤖 G.I.D.A Android Setup\n")

		// 1. Create widget shortcuts
		fmt.Println("📱 Creating Termux:Widget shortcuts...")
		if err := android.CreateWidgetScripts(); err != nil {
			fmt.Printf("❌ Widget creation failed: %v\n", err)
		} else {
			fmt.Println("✅ Widgets created in ~/.shortcuts/")
		}

		// 2. Schedule autonomous gardening
		fmt.Println("\n🌱 Scheduling autonomous gardening...")
		scheduler := android.NewJobScheduler()
		if err := scheduler.ScheduleGardening(1 * time.Hour); err != nil {
			fmt.Printf("❌ Scheduling failed: %v\n", err)
		} else {
			fmt.Println("✅ Gardener scheduled (hourly)")
		}

		// 3. Schedule cloud sync
		fmt.Println("\n☁️  Scheduling Turso cloud sync...")
		if err := scheduler.ScheduleSync(2 * time.Hour); err != nil {
			fmt.Printf("❌ Sync scheduling failed: %v\n", err)
		} else {
			fmt.Println("✅ Cloud sync scheduled (every 2 hours)")
		}

		// 4. Test notification
		fmt.Println("\n🔔 Testing native notifications...")
		if err := android.NotifyUser("G.I.D.A Setup", "Android integration complete!"); err != nil {
			fmt.Printf("⚠️  Notifications unavailable (install termux-api)\n")
		} else {
			fmt.Println("✅ Notifications working")
		}

		fmt.Println("\n🚀 G.I.D.A is now fully integrated with your Android device!")
		fmt.Println("   - Add widgets from Termux:Widget app")
		fmt.Println("   - Background jobs are now active")
		fmt.Println("   - Run 'ai_agent interactive' for the command center")
	},
}

func init() {
	rootCmd.AddCommand(setupAndroidCmd)
}
