package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"ai_agent_termux/config"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type Database struct {
	db     *sql.DB
	config *config.Config
}

// NewDatabase creates a new database connection
func NewDatabase(cfg *config.Config) (*Database, error) {
	db := &Database{config: cfg}

	if err := db.connect(); err != nil {
		return nil, err
	}

	return db, nil
}

// connect establishes database connection based on configuration
func (d *Database) connect() error {
	var err error

	if d.config.UseLocalSQLite {
		// Ensure database directory exists
		dbDir := filepath.Dir(d.config.SQLiteDBPath)
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			return fmt.Errorf("failed to create database directory: %w", err)
		}

		// Connect to local SQLite
		d.db, err = sql.Open("sqlite3", d.config.SQLiteDBPath)
		if err != nil {
			return fmt.Errorf("failed to open SQLite database: %w", err)
		}

		fmt.Printf("Connected to local SQLite database: %s\n", d.config.SQLiteDBPath)

	} else if d.config.TursoURL != "" && d.config.TursoAuthToken != "" {
		// Connect to Turso (libsql)
		dsn := fmt.Sprintf("%s?authToken=%s", d.config.TursoURL, d.config.TursoAuthToken)
		d.db, err = sql.Open("libsql", dsn)
		if err != nil {
			return fmt.Errorf("failed to connect to Turso: %w", err)
		}

		fmt.Println("Connected to Turso database")

	} else {
		return fmt.Errorf("no database configuration found (either local SQLite or Turso must be configured)")
	}

	// Test connection
	if err := d.db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Initialize tables
	return d.initTables()
}

// initTables creates necessary tables if they don't exist
func (d *Database) initTables() error {
	// Documents table - stores document metadata
	documentsTable := `
	CREATE TABLE IF NOT EXISTS documents (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		file_path TEXT UNIQUE NOT NULL,
		file_name TEXT NOT NULL,
		file_type TEXT NOT NULL,
		file_size INTEGER,
		processed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		summary TEXT,
		content_hash TEXT,
		embedding_count INTEGER DEFAULT 0
	);`

	// Embeddings table - stores vector embeddings
	embeddingsTable := `
	CREATE TABLE IF NOT EXISTS embeddings (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		document_id INTEGER,
		chunk_index INTEGER,
		chunk_text TEXT,
		embedding_data BLOB,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (document_id) REFERENCES documents (id) ON DELETE CASCADE
	);`

	// Search history table - stores search queries and results
	searchHistoryTable := `
	CREATE TABLE IF NOT EXISTS search_history (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		query TEXT NOT NULL,
		results TEXT, -- JSON string of results
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	// LLM usage table - tracks LLM API usage
	llmUsageTable := `
	CREATE TABLE IF NOT EXISTS llm_usage (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		provider TEXT NOT NULL,
		model TEXT NOT NULL,
		tokens_used INTEGER,
		cost_cents INTEGER,
		purpose TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	tables := []string{documentsTable, embeddingsTable, searchHistoryTable, llmUsageTable}

	for _, tableSQL := range tables {
		if _, err := d.db.Exec(tableSQL); err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}
	}

	return nil
}

// Close closes the database connection
func (d *Database) Close() error {
	if d.db != nil {
		return d.db.Close()
	}
	return nil
}

// GetDB returns the underlying sql.DB connection
func (d *Database) GetDB() *sql.DB {
	return d.db
}

// IsLocal returns true if using local SQLite
func (d *Database) IsLocal() bool {
	return d.config.UseLocalSQLite
}

// GetDatabaseInfo returns information about the database
func (d *Database) GetDatabaseInfo() (map[string]interface{}, error) {
	info := make(map[string]interface{})

	if d.IsLocal() {
		info["type"] = "SQLite"
		info["path"] = d.config.SQLiteDBPath
		info["location"] = "local"
	} else {
		info["type"] = "Turso"
		info["url"] = d.config.TursoURL
		info["location"] = "remote"
	}

	// Get table counts
	tables := []string{"documents", "embeddings", "search_history", "llm_usage"}
	for _, table := range tables {
		var count int
		err := d.db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", table)).Scan(&count)
		if err != nil {
			info[table+"_count"] = 0
		} else {
			info[table+"_count"] = count
		}
	}

	return info, nil
}

// HealthCheck performs a database health check
func (d *Database) HealthCheck() error {
	if d.db == nil {
		return fmt.Errorf("database not initialized")
	}

	// Test basic connectivity
	if err := d.db.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	// Test basic query
	var result int
	if err := d.db.QueryRow("SELECT 1").Scan(&result); err != nil {
		return fmt.Errorf("database query failed: %w", err)
	}

	if result != 1 {
		return fmt.Errorf("database query returned unexpected result")
	}

	return nil
}
