package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

// TursoSync manages cloud synchronization
type TursoSync struct {
	localDB   *Database
	cloudDB   *sql.DB
	syncURL   string
	authToken string
}

// NewTursoSync creates a Turso cloud sync instance
func NewTursoSync(localDB *Database) (*TursoSync, error) {
	tursoURL := os.Getenv("TURSO_DATABASE_URL")
	authToken := os.Getenv("TURSO_AUTH_TOKEN")

	if tursoURL == "" || authToken == "" {
		return nil, fmt.Errorf("TURSO_DATABASE_URL and TURSO_AUTH_TOKEN must be set")
	}

	cloudDB, err := sql.Open("libsql", tursoURL+"?authToken="+authToken)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Turso: %v", err)
	}

	return &TursoSync{
		localDB:   localDB,
		cloudDB:   cloudDB,
		syncURL:   tursoURL,
		authToken: authToken,
	}, nil
}

// SyncToCloud pushes new/updated documents to Turso
func (ts *TursoSync) SyncToCloud() error {
	// Get all documents updated in last sync window
	rows, err := ts.localDB.db.Query(`
		SELECT path, summary, embeddings, file_type, last_indexed
		FROM documents
		WHERE last_indexed > datetime('now', '-1 hour')
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	synced := 0
	for rows.Next() {
		var path, summary, embeddings, fileType string
		var lastIndexed time.Time

		if err := rows.Scan(&path, &summary, &embeddings, &fileType, &lastIndexed); err != nil {
			continue
		}

		// Upsert to cloud
		_, err = ts.cloudDB.Exec(`
			INSERT OR REPLACE INTO documents (path, summary, embeddings, file_type, last_indexed)
			VALUES (?, ?, ?, ?, ?)
		`, path, summary, embeddings, fileType, lastIndexed)

		if err == nil {
			synced++
		}
	}

	return nil
}

// Close closes the cloud connection
func (ts *TursoSync) Close() error {
	return ts.cloudDB.Close()
}
