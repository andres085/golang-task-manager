package main

import (
	"database/sql"
	"log/slog"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := sql.Open("mysql", "myuser:mypassword@/task_manager?parseTime=true")
	if err != nil {
		logger.Error("Failed to open database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		logger.Error("Failed to connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}

	script, err := os.ReadFile("./cmd/migrate/migrate.sql")
	if err != nil {
		logger.Error("Failed to read migration file", slog.String("error", err.Error()))
		os.Exit(1)
	}

	queries := strings.Split(string(script), ";")

	for _, query := range queries {
		query = strings.TrimSpace(query)
		if query == "" {
			continue
		}

		logger.Info("Executing query", slog.String("query", query))
		_, err = db.Exec(query)
		if err != nil {
			logger.Error("Failed to execute query", slog.String("query", query), slog.String("error", err.Error()))
			os.Exit(1)
		}
	}

	logger.Info("Migration completed successfully!")
}
