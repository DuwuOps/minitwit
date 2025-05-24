package datalayer

import (
	"database/sql"
	"fmt"
	"log/slog"
	"minitwit/src/handlers/helpers"
	"minitwit/src/utils"
	"os"
	"time"

	_ "github.com/lib/pq"
)

const (
	QueriesDirectory = "queries/"
	MaxRetries  = 10
	RetryDelay  = 2 * time.Second
)

func connectDB() (*sql.DB, error) {
	dbUserName := helpers.GetEnvVar("DB_USER", "admin")
	dbPass := helpers.GetEnvVar("DB_PASSWORD", "postgres")
	dbHost := helpers.GetEnvVar("DB_HOST", "localhost")
	dbPort := helpers.GetEnvVar("DB_PORT", "5433")
	dbName := helpers.GetEnvVar("DB_NAME", "minitwit")
	sslMode := helpers.GetEnvVar("DB_SSL_MODE", "disable")

	//Returns a new connection to the database.
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", dbUserName, dbPass, dbHost, dbPort, dbName, sslMode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		utils.LogError("sql.Open returned an error", err)
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Try multiple (10) times (sometimes it takes the postgres database a second to start)
	for i := range MaxRetries {
		err = db.Ping()
		if err == nil {
			slog.Info(fmt.Sprintf("Successfully connected to database on attempt %d", i+1))
			return db, nil
		}

		slog.Error("Failed to connect to database", slog.Any("error", err), slog.Any("current_attempt", i+1), slog.Any("max_attempts", MaxRetries))
		time.Sleep(RetryDelay)
	}

	return nil, fmt.Errorf("failed to connect to database: %w", err)
}

// Creates the database tables from query in {QueriesFile}
func createTablesIfNotExists(db *sql.DB) error {
	// Create table "users"
	err := createTableIfNotExists(db, "users")
	if err != nil {
		return err
	}

	// Create table "message"
	err = createTableIfNotExists(db, "message")
	if err != nil {
		return err
	}

	// Create table "follower"
	err = createTableIfNotExists(db, "follower")
	if err != nil {
		return err
	}

	// Create table "latest_processed"
	err = createTableIfNotExists(db, "latest_processed")
	if err != nil {
		return err
	}

	return nil
}

func createTableIfNotExists(db *sql.DB, tableName string) error {
	// Read queries-file
	QueriesFile := fmt.Sprintf("%sschema.%s.sql", QueriesDirectory, tableName)
	sqlFile, err := os.ReadFile(QueriesFile)
	if err != nil {
		utils.LogError("os.ReadFile returned an error", err)
		db.Close()
		return fmt.Errorf("failed to read schema file: %w", err)
	}

	// Execture contents of queries-file
	_, err = db.Exec(string(sqlFile))
	if err != nil {
		slog.Error("db.Exec returned an error", slog.Any("error", err), slog.Any("SQL-query", sqlFile))
		db.Close()
		return fmt.Errorf("failed to execute schema: %w", err)
	}

	return nil
}

func InitDB() (*sql.DB, error) {
	// Establish connection to database
	db, err := connectDB()
	if err != nil {
		utils.LogError("connectDB returned an error", err)
		return nil, err
	}

	err = createTablesIfNotExists(db)
	if err != nil {
		utils.LogError("InitDB : createTablesIfNotExists returned an error", err)
		return nil, err
	}

	slog.Info("Connecting to existing Minitwit Database!")
	return db, nil
}
