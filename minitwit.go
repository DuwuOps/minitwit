package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/labstack/echo"
	_ "github.com/mattn/go-sqlite3"
)

// configuration
var DATABASE = "./tmp/minitwit.db"

// TODO: Choose new web framework
// create our little application :)
// app = Flask(__name__)

func connectDB() (*sql.DB, error) {
	//Returns a new connection to the database.
	db, err := sql.Open("sqlite3", DATABASE)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	return db, nil
}

func initDB() (*sql.DB, error) {
	// Create tmp directory
	dir := filepath.Dir(DATABASE)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	// Establish connection to database
	db, err := connectDB()
	if err != nil {
		return nil, err
	}

	// Creates the database tables (and file if it does not exist yet).
	sqlFile, err := os.ReadFile("./schema.sql")
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to read schema file: %w", err)
	}
	_, err = db.Exec(string(sqlFile))
	if err != nil {
		db.Close()
		log.Printf("%q: %s\n", err, sqlFile)
		return nil, fmt.Errorf("failed to execute schema: %w", err)
	}

	return db, nil
}

// Note: Method signature and return type have been modified
func queryDB(db *sql.DB, query string, singleResult bool) (*sql.Rows, error) {
	//Queries the database and returns a list of QueryResults.
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	if singleResult {
		if !rows.Next() {
			rows.Close()
			return nil, fmt.Errorf("no rows found")
		}
	}

	return rows, nil
}


// ==========================
// Start: Routes
func setupRoutes(app *echo.Echo) {
	app.GET("/", Timeline)

    app.GET("/public", PublicTimeline)
	app.GET("/:username", UserTimeline)

	app.GET("/:username/follow", FollowUser)
    app.GET("/:username/unfollow", UnfollowUser)

	app.POST("/add_message", AddMessage)

    app.GET("/login", Login)
    app.POST("/login", Login)

    app.GET("/register", Register)
    app.POST("/register", Register)

    app.GET("/logout", Logout)	
}
// End: Routes
// ==========================

// ==========================
// Start: Route-Handlers
func Timeline(c echo.Context) error {
	return errors.New("Not implemented yet") //TODO
}

func PublicTimeline(c echo.Context) error {
	return errors.New("Not implemented yet") //TODO
}

func UserTimeline(c echo.Context) error {
	return errors.New("Not implemented yet") //TODO
}

func FollowUser(c echo.Context) error {
	return errors.New("Not implemented yet") //TODO
}

func UnfollowUser(c echo.Context) error {
	return errors.New("Not implemented yet") //TODO
}

func AddMessage(c echo.Context) error {
	return errors.New("Not implemented yet") //TODO
}

func Login(c echo.Context) error {
    return errors.New("Not implemented yet") //TODO
}

func Register(c echo.Context) error {
    return errors.New("Not implemented yet") //TODO
}

func Logout(c echo.Context) error {
    return errors.New("Not implemented yet") //TODO
}
// End: Route-Handlers
// ==========================


// Example
func main() {
	db, err := initDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	fmt.Println("Database initialized successfully")

	rows, err := queryDB(db, "select * from user", false)
	if err != nil {
		log.Fatalf("Failed to query database: %v", err)
	}
	defer rows.Close()

	fmt.Println("Query executed successfully")
}
