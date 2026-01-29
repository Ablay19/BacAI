package adapter

import (
	"fmt"
	"os"
	"path/filepath"
)

// NushellAdapter generates completion and integration scripts for Nushell
type NushellAdapter struct {
	configDir string
}

// NewNushellAdapter creates a new Nushell adapter
func NewNushellAdapter() *NushellAdapter {
	homeDir, _ := os.UserHomeDir()
	return &NushellAdapter{
		configDir: filepath.Join(homeDir, ".config", "nushell"),
	}
}

// GenerateCompletionScript creates a nushell completion script for ai_agent
func (na *NushellAdapter) GenerateCompletionScript() error {
	script := `
# G.I.D.A Nushell Completions
export extern "ai_agent" [
    command?: string@ai_agent_commands # The command to run
    --debug(-d)                         # Enable debug logging
]

def ai_agent_commands [] {
    [ "scan", "search", "garden", "interactive", "export", "setup-android", "version" ]
}

export extern "ai_agent scan" [
    --enhanced                          # Use enhanced scanning
    --workers: int                      # Number of parallel workers
]

export extern "ai_agent search" [
    query: string                       # The search query
]
`
	os.MkdirAll(na.configDir, 0755)
	scriptPath := filepath.Join(na.configDir, "ai_agent_completions.nu")
	return os.WriteFile(scriptPath, []byte(script), 0644)
}

// InstallationInstructions returns instructions for the user
func (na *NushellAdapter) InstallationInstructions() string {
	return fmt.Sprintf(`To enable Nushell completions, add the following to your env.nu or config.nu:
    
    use %s`, filepath.Join(na.configDir, "ai_agent_completions.nu"))
}
