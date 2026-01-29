package cmd

import (
	"context"
	"time"

	"ai_agent_termux/config"
	"ai_agent_termux/database"
	"ai_agent_termux/gardener"
	"ai_agent_termux/goroutines"

	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
)

var (
	gardenInterval time.Duration
)

var gardenCmd = &cobra.Command{
	Use:   "garden",
	Short: "Start the autonomous data gardener",
	Long: `Starts G.I.D.A in autonomous mode. The gardener will run in the background, 
periodically pruning dead metadata, refreshing summaries, and optimizing the device's 
indexed knowledge base.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.LoadConfig()

		// Initialize Core Components
		db, err := database.NewDatabase(cfg)
		if err != nil {
			slog.Error("Failed to initialize database", "error", err)
			return
		}

		tm := goroutines.NewTaskManager(nil)

		// Initialize Gardener
		g := gardener.NewGardener(db, tm, gardenInterval)

		slog.Info("G.I.D.A Autonomous Mode: ENGAGED")

		ctx := context.Background()
		g.Start(ctx)
	},
}

func init() {
	gardenCmd.Flags().DurationVarP(&gardenInterval, "interval", "i", 1*time.Hour, "Interval between gardening sessions")
}
