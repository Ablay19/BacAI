package cmd

import (
	"fmt"
	"os"

	"ai_agent_termux/config"
	"ai_agent_termux/utils"

	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
)

var (
	cfgFile    string
	debug      bool
	configPath string
	logPath    string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ai_agent",
	Short: "A.I.D.A - Android Intelligent Document Analyst",
	Long: `A.I.D.A. (Android Intelligent Document Analyst) is an autonomous AI agent 
for Termux environments that discovers, scans, extracts, analyzes, categorizes, 
summarizes, indexes, and semantically searches through documents.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ai_agent.yaml)")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable debug logging")
	rootCmd.PersistentFlags().StringVar(&configPath, "config-path", "", "path to configuration directory")
	rootCmd.PersistentFlags().StringVar(&logPath, "log-path", "", "path to log file")

	// Register subcommands
	rootCmd.AddCommand(scanCmd)
	rootCmd.AddCommand(preprocessCmd)
	rootCmd.AddCommand(summarizeCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(dbCmd)
	rootCmd.AddCommand(androidCmd)
	rootCmd.AddCommand(gardenCmd)
	rootCmd.AddCommand(versionCmd)

	// Add completion command
	rootCmd.AddCommand(&cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate completion script",
		Long: `To load completions:

Bash:

$ source <(ai_agent completion bash)

# To load completions for each session, execute once:
Linux:
  $ ai_agent completion bash > /etc/bash_completion.d/ai_agent
MacOS:
  $ ai_agent completion bash > /usr/local/etc/bash_completion.d/ai_agent

Zsh:

# If shell completion is not already enabled in your environment you will need
# to enable it.  You can execute the following once:

$ echo "autoload -U compinit; compinit" >> ~/.zshrc

# To load completions for each session, execute once:
$ ai_agent completion zsh > "${fpath[1]}/_ai_agent"

# You will need to start a new shell for this setup to take effect.

Fish:

$ ai_agent completion fish | source

# To load completions for each session, execute once:
$ ai_agent completion fish > ~/.config/fish/completions/ai_agent.fish

Nushell:

For Nushell users, you can use the bash completion:
$ ai_agent completion bash | save ~/.cache/ai_agent-completion.nu

Then in your Nushell config (~/.config/nushell/config.nu), add:
source ~/.cache/ai_agent-completion.nu

For improved experience, install carapace from https://github.com/rsteube/carapace-bin
and add to your config:
source ~/.config/carapace/completers.nu
`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		Run: func(cmd *cobra.Command, args []string) {
			switch args[0] {
			case "bash":
				cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				cmd.Root().GenFishCompletion(os.Stdout, true)
			case "powershell":
				cmd.Root().GenPowerShellCompletion(os.Stdout)
			}
		},
	})
}

func initConfig() {
	// Initialize configuration
	cfg := config.LoadConfig()

	// Override with command line flags if provided
	if configPath != "" {
		// TODO: Implement custom config path handling
		slog.Info("Using custom config path", "path", configPath)
	}

	// Initialize logger
	logFilePath := cfg.LogFilePath
	if logPath != "" {
		logFilePath = logPath
	}

	utils.InitLogger(logFilePath)

	if debug {
		slog.Info("Debug mode enabled")
	}

	slog.Info("Configuration loaded", "config", cfg)
}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of AI Agent",
	Long:  `All software has versions. This is AI Agent's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("AI Agent v0.1 -- HEAD")
	},
}
