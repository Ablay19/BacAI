package cmd

import (
	"fmt"
	"os"

	"ai_agent_termux/config"
	"ai_agent_termux/database"
	"ai_agent_termux/exporter"

	"github.com/spf13/cobra"
)

var (
	vaultPath string
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export summaries to Obsidian/Logseq",
	Long:  `Automatically exports all document summaries as Markdown notes to your second brain.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.LoadConfig()
		db, err := database.NewDatabase(cfg)
		if err != nil {
			fmt.Printf("❌ Database error: %v\n", err)
			os.Exit(1)
		}
		defer db.Close()

		if vaultPath == "" {
			vaultPath = os.Getenv("OBSIDIAN_VAULT")
			if vaultPath == "" {
				fmt.Println("❌ Please specify --vault or set OBSIDIAN_VAULT")
				os.Exit(1)
			}
		}

		fmt.Printf("📝 Exporting to: %s\n", vaultPath)

		exp := exporter.NewObsidianExporter(vaultPath, db)
		if err := exp.ExportAll(); err != nil {
			fmt.Printf("❌ Export failed: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("✅ Export complete! Check your vault.")
	},
}

func init() {
	exportCmd.Flags().StringVarP(&vaultPath, "vault", "v", "", "Path to Obsidian vault")
	rootCmd.AddCommand(exportCmd)
}
