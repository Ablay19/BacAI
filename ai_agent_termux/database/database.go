package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"ai_agent_termux/config"

	_ "github.com/mattn/go-sqlite3"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type Database struct {
	db     *sql.DB
	config *config.Config
}

// Document represents a document in the database
type Document struct {
	ID          int
	Path        string
	Name        string
	FileType    string
	FileSize    int64
	LastIndexed time.Time
	Summary     string
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

// IsFileProcessed checks if a file with the given path and hash has already been processed
func (d *Database) IsFileProcessed(path, hash string) (bool, error) {
	if d.db == nil {
		return false, fmt.Errorf("database not initialized")
	}

	var count int
	// We check both path and hash to ensure it's the exact same file
	err := d.db.QueryRow("SELECT COUNT(*) FROM documents WHERE file_path = ? AND content_hash = ?", path, hash).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// UpdateDocumentHash updates or inserts a document's hash and metadata
func (d *Database) UpdateDocumentHash(path, name, fileType string, size int64, hash string) error {
	if d.db == nil {
		return fmt.Errorf("database not initialized")
	}

	query := `
	INSERT INTO documents (file_path, file_name, file_type, file_size, content_hash, processed_at)
	VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
	ON CONFLICT(file_path) DO UPDATE SET
		content_hash = excluded.content_hash,
		processed_at = excluded.processed_at,
		file_size = excluded.file_size;
	`
	_, err := d.db.Exec(query, path, name, fileType, size, hash)
	return err
}

// SaveSummary updates the summary for a document
func (d *Database) SaveSummary(path, summary string) error {
	if d.db == nil {
		return fmt.Errorf("database not initialized")
	}

	_, err := d.db.Exec("UPDATE documents SET summary = ? WHERE file_path = ?", summary, path)
	return err
}

// GetAllDocuments returns all documents in the database
func (d *Database) GetAllDocuments() ([]Document, error) {
	rows, err := d.db.Query("SELECT id, file_path, file_name, file_type, file_size, processed_at, summary FROM documents")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var docs []Document
	for rows.Next() {
		var doc Document
		var processedAt string
		err := rows.Scan(&doc.ID, &doc.Path, &doc.Name, &doc.FileType, &doc.FileSize, &processedAt, &doc.Summary)
		if err != nil {
			continue
		}
		doc.LastIndexed, _ = time.Parse("2006-01-02 15:04:05", processedAt)
		docs = append(docs, doc)
	}
	return docs, nil
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
