package main

import (
	"database/sql"
	"embed"
	"fmt"
	"log/slog"
	"sort"
	"strings"

	_ "github.com/lib/pq"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

// runMigrations applies all *.up.sql files from the migrations/ directory in order.
// Files are named: 001_description.up.sql, 002_description.up.sql, ...
func runMigrations(db *sql.DB) {
	entries, err := migrationFiles.ReadDir("migrations")
	if err != nil {
		slog.Error("migrations: list failed", "err", err)
		return
	}

	var upFiles []string
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".up.sql") {
			upFiles = append(upFiles, e.Name())
		}
	}
	sort.Strings(upFiles)

	for _, f := range upFiles {
		data, err := migrationFiles.ReadFile("migrations/" + f)
		if err != nil {
			slog.Error("migrations: read failed", "file", f, "err", err)
			continue
		}
		if _, err := db.Exec(string(data)); err != nil {
			// Log but don't fatal — migration may already be applied.
			slog.Warn("migrations: exec warning", "file", f, "err", fmt.Sprintf("%v", err))
		} else {
			slog.Info("migrations: applied", "file", f)
		}
	}
}
