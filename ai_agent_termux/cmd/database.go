package cmd

import (
	"fmt"
	"os"

	"ai_agent_termux/config"
	"ai_agent_termux/database"

	"github.com/spf13/cobra"
)

// dbCmd represents the database command
var dbCmd = &cobra.Command{
	Use:   "database [command]",
	Short: "Manage database operations",
	Long: `Manage SQLite and Turso database operations.
	Database stores document metadata, embeddings, and search history.`,
}

func init() {
	// Add subcommands
	dbCmd.AddCommand(dbInfoCmd)
	dbCmd.AddCommand(dbHealthCmd)
	dbCmd.AddCommand(dbMigrateCmd)
}

// dbInfoCmd represents the database info command
var dbInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show database information",
	Long:  `Display current database configuration and statistics.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.LoadConfig()

		fmt.Println("Database Information:")
		fmt.Println("====================")

		if cfg.UseLocalSQLite {
			fmt.Printf("Type: Local SQLite\n")
			fmt.Printf("Path: %s\n", cfg.SQLiteDBPath)
		} else {
			fmt.Printf("Type: Turso Cloud\n")
			fmt.Printf("URL: %s\n", cfg.TursoURL)
		}

		// Try to connect and get stats
		db, err := database.NewDatabase(cfg)
		if err != nil {
			fmt.Printf("❌ Failed to connect to database: %v\n", err)
			return
		}
		defer db.Close()

		info, err := db.GetDatabaseInfo()
		if err != nil {
			fmt.Printf("❌ Failed to get database info: %v\n", err)
			return
		}

		fmt.Println("\nDatabase Statistics:")
		fmt.Printf("Documents: %v\n", info["documents_count"])
		fmt.Printf("Embeddings: %v\n", info["embeddings_count"])
		fmt.Printf("Search History: %v\n", info["search_history_count"])
		fmt.Printf("LLM Usage Records: %v\n", info["llm_usage_count"])
	},
}

// dbHealthCmd represents the database health command
var dbHealthCmd = &cobra.Command{
	Use:   "health",
	Short: "Check database health",
	Long:  `Perform health check on the database connection and functionality.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.LoadConfig()

		fmt.Println("Database Health Check:")
		fmt.Println("=====================")

		// Check configuration
		if !cfg.UseLocalSQLite && (cfg.TursoURL == "" || cfg.TursoAuthToken == "") {
			fmt.Println("❌ No database configuration found")
			fmt.Println("   Please configure either local SQLite or Turso database")
			os.Exit(1)
		}

		// Try to connect
		db, err := database.NewDatabase(cfg)
		if err != nil {
			fmt.Printf("❌ Database connection failed: %v\n", err)
			os.Exit(1)
		}
		defer db.Close()

		// Perform health check
		if err := db.HealthCheck(); err != nil {
			fmt.Printf("❌ Database health check failed: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("✅ Database is healthy and accessible")

		if db.IsLocal() {
			fmt.Println("✅ Using local SQLite database")
		} else {
			fmt.Println("✅ Using Turso cloud database")
		}
	},
}

// dbMigrateCmd represents the database migrate command
var dbMigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Initialize database schema",
	Long:  `Create database tables and initialize schema. This is safe to run multiple times.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.LoadConfig()

		fmt.Println("Database Migration:")
		fmt.Println("===================")

		// Check configuration
		if !cfg.UseLocalSQLite && (cfg.TursoURL == "" || cfg.TursoAuthToken == "") {
			fmt.Println("❌ No database configuration found")
			fmt.Println("   Please configure either local SQLite or Turso database")
			os.Exit(1)
		}

		// Try to connect (this automatically runs initTables)
		db, err := database.NewDatabase(cfg)
		if err != nil {
			fmt.Printf("❌ Database migration failed: %v\n", err)
			os.Exit(1)
		}
		defer db.Close()

		fmt.Println("✅ Database migration completed successfully")
		fmt.Println("   All required tables have been created or verified")
	},
}
