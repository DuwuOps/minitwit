package datalayer

import (
	"database/sql"
	"fmt"
	"log"
	"minitwit/src/handlers/helpers"
	"os"
	"time"

	_ "github.com/lib/pq"
)

const (
	QueriesFile = "queries/schema.sql"
	MaxRetries 	= 10
	RetryDelay  = 2 * time.Second
)

func connectDB() (*sql.DB, error) {
	DbUserName := helpers.GetEnvVar("DB_USER", "admin")
	DbPass := helpers.GetEnvVar("DB_PASSWORD", "localhost")
	DbHost := helpers.GetEnvVar("DB_HOST", "database")
	DbPort := helpers.GetEnvVar("DB_PORT", "5433")
	DbName := helpers.GetEnvVar("DB_NAME", "minitwit")
	SSLMode := helpers.GetEnvVar("DB_SSL_MODE", "disable")

	//Returns a new connection to the database.
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", DbUserName, DbPass, DbHost, DbPort, DbName, SSLMode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("sql.Open returned error: %v\n", err)
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	
	// Try multiple (10) times (sometimes it takes the postgres database a second to start)
	for i := range MaxRetries {
		err = db.Ping()
		if err == nil {
			fmt.Printf("Successfully connected to database on attempt %d\n", i+1)
			return db, nil
		}
		
		fmt.Printf("Failed to connect to database (attempt %d/%d): %v\n", i+1, MaxRetries, err)
		time.Sleep(RetryDelay)
	}
	
	return nil, fmt.Errorf("failed to connect to database: %w", err)
}

// Creates the database tables from quer
func CreateTables(db *sql.DB) error {
	// Read queries-file
	sqlFile, err := os.ReadFile(QueriesFile)
	if err != nil {
		log.Printf("os.ReadFile returned error: %v\n", err)
		db.Close()
		return fmt.Errorf("failed to read schema file: %w", err)
	}
	
	// Execture contents of queries-file
	_, err = db.Exec(string(sqlFile))
	if err != nil {
		log.Printf("db.Exec returned error: %v\n", err)
		db.Close()
		log.Printf("%q: %s\n", err, sqlFile)
		return fmt.Errorf("failed to execute schema: %w", err)
	}

	return nil
}

func InitDB() (*sql.DB, error) {
	// Establish connection to database
	db, err := connectDB()
	if err != nil {
		log.Printf("connectDB returned error: %v\n", err)
		return nil, err
	}

	// TODO: Fix that this would drop and create the tables at every reset (pretty bad)
	err = CreateTables(db)
	if err != nil {
		log.Printf("CreateTables returned error: %v\n", err)
		return nil, err
	}

	log.Printf("Connecting to existing Minitwit Database!")
	return db, nil
}