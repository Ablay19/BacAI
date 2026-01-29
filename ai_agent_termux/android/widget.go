package android

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// WidgetCommand creates Termux:Widget compatible scripts
type WidgetCommand struct {
	Name        string
	Description string
	Script      string
}

// CreateWidgetScripts generates all widget shortcuts in ~/.shortcuts/
func CreateWidgetScripts() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	shortcutsDir := filepath.Join(homeDir, ".shortcuts")
	if err := os.MkdirAll(shortcutsDir, 0755); err != nil {
		return err
	}

	widgets := []WidgetCommand{
		{
			Name:        "G.I.D.A Quick Scan",
			Description: "One-tap full device scan",
			Script: `#!/data/data/com.termux/files/usr/bin/bash
cd ~/storage/shared
ai_agent scan --enhanced --workers 16
termux-notification --title "G.I.D.A" --content "Scan complete!"
`,
		},
		{
			Name:        "G.I.D.A Garden",
			Description: "Run autonomous gardener",
			Script: `#!/data/data/com.termux/files/usr/bin/bash
ai_agent garden --interval 1h &
termux-notification --title "G.I.D.A" --content "Gardener activated"
`,
		},
		{
			Name:        "G.I.D.A TUI",
			Description: "Launch premium command center",
			Script: `#!/data/data/com.termux/files/usr/bin/bash
ai_agent interactive
`,
		},
	}

	for _, w := range widgets {
		scriptPath := filepath.Join(shortcutsDir, w.Name)
		if err := os.WriteFile(scriptPath, []byte(w.Script), 0755); err != nil {
			return fmt.Errorf("failed to create widget '%s': %v", w.Name, err)
		}
	}

	return nil
}

// NotifyUser sends a native Android notification
func NotifyUser(title, content string) error {
	cmd := exec.Command("termux-notification", "--title", title, "--content", content)
	return cmd.Run()
}

// NotifyWithAction sends a notification with a clickable action
func NotifyWithAction(title, content, action string) error {
	cmd := exec.Command("termux-notification",
		"--title", title,
		"--content", content,
		"--action", action,
	)
	return cmd.Run()
}
