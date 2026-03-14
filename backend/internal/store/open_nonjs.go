//go:build !js

package store

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

const defaultSQLitePath = "app.db"

func openDB() (*sql.DB, error) {
	return sql.Open("sqlite", sqlitePath())
}

func SourceDescription() string {
	return fmt.Sprintf("local SQLite file %q", sqlitePath())
}

func sqlitePath() string {
	if value := os.Getenv("DB_PATH"); value != "" {
		return value
	}

	return defaultSQLitePath
}
