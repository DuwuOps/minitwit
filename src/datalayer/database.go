package datalayer

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var DATABASE = "tmp/minitwit.db"

func connectDB() (*sql.DB, error) {
	//Returns a new connection to the database.
	db, err := sql.Open("sqlite3", DATABASE)
	if err != nil {
		fmt.Printf("sql.Open returned error: %v\n", err)
		db.Close()
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	return db, nil
}

func InitDB() (*sql.DB, error) {
	// Create tmp directory
	dir := filepath.Dir(DATABASE)

	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			//Create Database if not exists!
			fmt.Printf("Creating new Minitwit Database!")

			err := os.MkdirAll(dir, 0755)
			if err != nil {
				fmt.Printf("os.MkdirAll returned error: %v\n", err)
				return nil, fmt.Errorf("failed to create database directory: %w", err)
			}

			// Establish connection to database
			db, err := connectDB()
			if err != nil {
				fmt.Printf("connectDB returned error: %v\n", err)
				return nil, err
			}

			// Creates the database tables (and file if it does not exist yet).
			sqlFile, err := os.ReadFile("queries/schema.sql")
			if err != nil {
				fmt.Printf("os.ReadFile returned error: %v\n", err)
				db.Close()
				return nil, fmt.Errorf("failed to read schema file: %w", err)
			}
			_, err = db.Exec(string(sqlFile))
			if err != nil {
				fmt.Printf("db.Exec returned error: %v\n", err)
				db.Close()
				log.Printf("%q: %s\n", err, sqlFile)
				return nil, fmt.Errorf("failed to execute schema: %w", err)
			}

			return db, nil

		} else {
			fmt.Printf("Unexpected error when connecting to Minitwit Database!")
		}
	}

	// Establish connection to database
	db, err := connectDB()
	if err != nil {
		fmt.Printf("connectDB returned error: %v\n", err)
		return nil, err
	}

	fmt.Printf("Connecting to existing Minitwit Database!")
	return db, nil
}

func QueryDbSingle(db *sql.DB, query string, args ...any) *sql.Row {
	row := db.QueryRow(query, args...)
	return row
}

func QueryDB(db *sql.DB, query string, args ...any) (*sql.Rows, error) {
	fmt.Printf("\nqueryDB called with following arguments:\n")
	fmt.Printf(" - query: %v\n", query)
	fmt.Printf(" - args: %v\n\n", args)

	//Queries the database and returns a list of QueryResults.
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		defer rows.Close()
		fmt.Printf("db.Query returned error: %v\n", err)
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return rows, nil
}

func GetUserId(username string, db *sql.DB) (int, error) {
	var id int
	err := db.QueryRow(`SELECT user_id FROM user WHERE username = ?`, username).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil // user not found
	}
	if err != nil {
		fmt.Printf("Db.QueryRow returned error: %v\n", err)
		return 0, err
	}
	return id, nil
}
