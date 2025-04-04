package datalayer

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var DATABASE = "tmp/minitwit.db"
var DATABASE_NAME = "minitwit.db"

func connectDB() (*sql.DB, error) {
	//Returns a new connection to the database.
	db, err := sql.Open("sqlite3", DATABASE)
	if err != nil {
		log.Printf("sql.Open returned error: %v\n", err)
		db.Close()
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	return db, nil
}

func InitDB() (*sql.DB, error) {
	// Create tmp directory
	dir := filepath.Dir(DATABASE)

	if _, err := os.Stat(dir + "/" + DATABASE_NAME); err != nil {
		if os.IsNotExist(err) {
			//Create Database if not exists!
			log.Printf("Creating new Minitwit Database!")

			err := os.MkdirAll(dir, 0755)
			if err != nil {
				log.Printf("os.MkdirAll returned error: %v\n", err)
				return nil, fmt.Errorf("failed to create database directory: %w", err)
			}

			// Establish connection to database
			db, err := connectDB()
			if err != nil {
				log.Printf("connectDB returned error: %v\n", err)
				return nil, err
			}

			// Creates the database tables (and file if it does not exist yet).
			sqlFile, err := os.ReadFile("queries/schema.sql")
			if err != nil {
				log.Printf("os.ReadFile returned error: %v\n", err)
				db.Close()
				return nil, fmt.Errorf("failed to read schema file: %w", err)
			}
			_, err = db.Exec(string(sqlFile))
			if err != nil {
				log.Printf("db.Exec returned error: %v\n", err)
				db.Close()
				log.Printf("%q: %s\n", err, sqlFile)
				return nil, fmt.Errorf("failed to execute schema: %w", err)
			}

			return db, nil

		} else {
			log.Printf("Unexpected error when connecting to Minitwit Database!")
		}
	}

	// Establish connection to database
	db, err := connectDB()
	if err != nil {
		log.Printf("connectDB returned error: %v\n", err)
		return nil, err
	}

	log.Printf("Connecting to existing Minitwit Database!")
	return db, nil
}