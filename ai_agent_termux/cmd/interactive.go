package cmd

import (
	"ai_agent_termux/automation"
	"ai_agent_termux/config"
	"ai_agent_termux/goroutines"
	"ai_agent_termux/interactive"

	"github.com/spf13/cobra"
)

var interactiveCmd = &cobra.Command{
	Use:   "interactive",
	Short: "Launch the G.I.D.A Premium TUI",
	Long:  `Starts the premium, Claude-style Terminal User Interface for G.I.D.A.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.LoadConfig()
		tm := goroutines.NewTaskManager(nil)
		glp := automation.NewGoogleLensProcessor(cfg)

		cli := interactive.NewInteractiveCLI(&interactive.CLIConfig{
			TaskManager: tm,
			GoogleLens:  glp,
		})

		cli.Start()
	},
}

func init() {
	rootCmd.AddCommand(interactiveCmd)
}
