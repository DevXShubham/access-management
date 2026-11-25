package sqlite

import (
	"database/sql"
	"fmt"
	"iam/apperrors"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func Connect() *sql.DB {
	db, err := sql.Open("sqlite3", "database/repository.db")

	if err != nil {
		log.Fatalf("%v: %v", apperrors.ErrDBConnection, err)
		return nil
	}

	if err := IntializeSchema(db); err != nil {
		log.Fatalf("%v: %v", apperrors.ErrDBConnection, err)
		return nil
	}

	return db
}

func IntializeSchema(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS tenants (
            id TEXT NOT NULL PRIMARY KEY,
            name TEXT NOT NULL UNIQUE,
            description TEXT NOT NULL,
            createdAt TEXT NOT NULL
        )`,

		`CREATE INDEX IF NOT EXISTS idx_tenants_created_at_id
        ON tenants(createdAt DESC, id DESC)`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("initializing database schema: %w", err)
		}
	}

	return nil
}
