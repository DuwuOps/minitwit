package datalayer

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

const (
	DbHost      = "database"
	DbUserName  = "admin"
	DbPass      = "postgres"
	DbName      = "minitwit"
	DbPort      = 5432
	QueriesFile = "queries/schema.sql"
	SSLMode		= "disable"
)

func connectDB() (*sql.DB, error) {
	//Returns a new connection to the database.
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s", DbUserName, DbPass, DbHost, DbPort, DbName, SSLMode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Printf("sql.Open returned error: %v\n", err)
		return nil, err
	}
	
	err = db.Ping()
	if err != nil {
		fmt.Printf("db.Ping returned error: %v\n", err)
		return nil, err
	}
	
	return nil, fmt.Errorf("failed to connect to database: %w", err)
}

// Creates the database tables from quer
func CreateTables(db *sql.DB) error {
	// Read queries-file
	sqlFile, err := os.ReadFile(QueriesFile)
	if err != nil {
		fmt.Printf("os.ReadFile returned error: %v\n", err)
		db.Close()
		return fmt.Errorf("failed to read schema file: %w", err)
	}
	
	// Execture contents of queries-file
	_, err = db.Exec(string(sqlFile))
	if err != nil {
		fmt.Printf("db.Exec returned error: %v\n", err)
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
		fmt.Printf("connectDB returned error: %v\n", err)
		return nil, err
	}

	// TODO: Fix that this would drop and create the tables at every reset (pretty bad)
	err = CreateTables(db)
	if err != nil {
		fmt.Printf("CreateTables returned error: %v\n", err)
		return nil, err
	}

	fmt.Printf("Connecting to existing Minitwit Database!")
	return db, nil
}