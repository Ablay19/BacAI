package android

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// JobScheduler manages background tasks via cron/termux-job-scheduler
type JobScheduler struct {
	cronPath string
}

// NewJobScheduler creates a new job scheduler instance
func NewJobScheduler() *JobScheduler {
	homeDir, _ := os.UserHomeDir()
	return &JobScheduler{
		cronPath: filepath.Join(homeDir, ".termux", "crontabs", "default"),
	}
}

// ScheduleGardening enables autonomous background gardening
func (js *JobScheduler) ScheduleGardening(interval time.Duration) error {
	cronExpr := "0 */1 * * *" // Every hour by default

	if interval == 30*time.Minute {
		cronExpr = "*/30 * * * *"
	} else if interval == 6*time.Hour {
		cronExpr = "0 */6 * * *"
	}

	jobLine := fmt.Sprintf("%s ai_agent garden --interval %s &>> /data/data/com.termux/files/home/.ai_garden.log\n",
		cronExpr, interval.String())

	// Read existing crontab
	existing, _ := os.ReadFile(js.cronPath)

	// Append job if not already present
	newContent := string(existing)
	if !contains(newContent, "ai_agent garden") {
		newContent += jobLine
		if err := os.WriteFile(js.cronPath, []byte(newContent), 0644); err != nil {
			return err
		}
	}

	// Reload cron
	return exec.Command("crond", "-c", filepath.Dir(js.cronPath)).Run()
}

// ScheduleSync enables periodic Turso cloud sync
func (js *JobScheduler) ScheduleSync(interval time.Duration) error {
	cronExpr := "0 */2 * * *" // Every 2 hours

	jobLine := fmt.Sprintf("%s ai_agent sync --turso &>> /data/data/com.termux/files/home/.ai_sync.log\n",
		cronExpr)

	existing, _ := os.ReadFile(js.cronPath)
	newContent := string(existing)

	if !contains(newContent, "ai_agent sync") {
		newContent += jobLine
		if err := os.WriteFile(js.cronPath, []byte(newContent), 0644); err != nil {
			return err
		}
	}

	return exec.Command("crond", "-c", filepath.Dir(js.cronPath)).Run()
}

func contains(haystack, needle string) bool {
	return strings.Contains(haystack, needle)
}
