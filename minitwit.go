package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/sessions"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

// configuration
var DATABASE = "./tmp/minitwit.db"
var SECRET_KEY = []byte("development key") // to parallel the Python "SECRET_KEY"
var Db *sql.DB

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

// Logs the user in.
func Login(c echo.Context) error {
    loggedIn, err := isUserLoggedIn(c)
	if err != nil {
		return err
	}
    if loggedIn {
        return c.Redirect(http.StatusFound, "/")
    }

    var errorMessage string
    if c.Request().Method == http.MethodPost {
        username := c.FormValue("username")
        password := c.FormValue("password")

        var dbUser struct {
            UserID int
            PwHash string
        }
        err := Db.QueryRow(`
            SELECT user_id, pw_hash FROM user
            WHERE username = ?
        `, username).Scan(&dbUser.UserID, &dbUser.PwHash)

        if errors.Is(err, sql.ErrNoRows) {
            errorMessage = "Invalid username"
        } else if err != nil {
            return err
        } else {
            if !checkPasswordHash(dbUser.PwHash, password) {
                errorMessage = "Invalid password"
            } else {
                addFlash(c, "You were logged in")
                setSessionUserID(c, dbUser.UserID)
                return c.Redirect(http.StatusFound, "/")
            }
        }
    }

	flashes, err := getFlashes(c)
	if err != nil {
		return err
	}

    data := map[string]interface{}{
        "error": errorMessage,
        "flashes": flashes,
    }
    return c.Render(http.StatusOK, "login.html", data)
}

func Register(c echo.Context) error {
    loggedIn, err := isUserLoggedIn(c)
	if err != nil {
		return err
	}
    if loggedIn {
        return c.Redirect(http.StatusFound, "/")
    }

    var errorMessage string
    if c.Request().Method == http.MethodPost {
        username := c.FormValue("username")
        email := c.FormValue("email")
        password := c.FormValue("password")
        password2 := c.FormValue("password2")

        switch {
        case username == "":
            errorMessage = "You have to enter a username"
        case email == "" || !strings.Contains(email, "@"):
            errorMessage = "You have to enter a valid email address"
        case password == "":
            errorMessage = "You have to enter a password"
        case password != password2:
            errorMessage = "The two passwords do not match"
        default:
            existingID, err := getUserId(username)
            if err != nil {
                return err
            }
            if existingID != 0 {
                errorMessage = "The username is already taken"
            } else {
                hash, err := generatePasswordHash(password)
                if err != nil {
                    return err
                }
                _, err = Db.Exec(`
                    INSERT INTO user (username, email, pw_hash)
                    VALUES (?, ?, ?)
                `, username, email, hash)
                if err != nil {
                    return err
                }

                addFlash(c, "You were successfully registered and can login now")
                return c.Redirect(http.StatusFound, "/login")
            }
        }
    }

	flashes, err := getFlashes(c)
	if err != nil {
		return err
	}

    data := map[string]interface{}{
        "error":   errorMessage,
        "flashes": flashes,
    }
    return c.Render(http.StatusOK, "register.html", data)
}

func Logout(c echo.Context) error {
    addFlash(c, "You were logged out")
    clearSessionUserID(c)
    return c.Redirect(http.StatusFound, "/public")
}
// End: Route-Handlers
// ==========================


// ==========================
// Start: Helpers
// Securely check that the given stored password hash matches the given password.
func checkPasswordHash(hashedPassword, plainPassword string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
    return err == nil
}

func setSessionUserID(c echo.Context, userID int) error {
	sess, err := session.Get("session", c)
	if err != nil {
		return err
	}
    sess.Values["user_id"] = userID
    sess.Save(c.Request(), c.Response())
	return nil
}

func clearSessionUserID(c echo.Context) error {
    sess, err := session.Get("session", c)
	if err != nil {
		return err
	}
    delete(sess.Values, "user_id")
    sess.Save(c.Request(), c.Response())
	return nil
}

// Take a context and returns whether the current user is logged in
func isUserLoggedIn(c echo.Context) (bool, error) {
	sess, err := session.Get("session", c)
    _, ok := sess.Values["user_id"].(int)
    return ok, err
}

// Takes a username and return the user's id
func getUserId(username string) (int, error) {
	var id int
    err := Db.QueryRow(`SELECT user_id FROM user WHERE username = ?`, username).Scan(&id)
    if errors.Is(err, sql.ErrNoRows) {
        return 0, nil // user not found
    } 
    if err != nil {
        return 0, err
    }
    return id, nil
}

// Takes a password-string and returns a hashed version of the password-string
func generatePasswordHash(password string) (string, error) {
    hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(hashBytes), err
}

// Takes a message to be flashed and a context
// Flashes a message to the next request
func addFlash(c echo.Context, message string) error {
	sess, err := session.Get("session", c)
	if err != nil {
		return err
	}
    flashes, ok := sess.Values["flashes"].([]string)
    if !ok {
        flashes = []string{}
    }
    flashes = append(flashes, message)
    sess.Values["flashes"] = flashes
    sess.Save(c.Request(), c.Response())
	return nil
}

// Takes a context
// Returns empties the flashes in the given context and returns the flashes in a list of strings
func getFlashes(c echo.Context) ([]string, error) {
	sess, err := session.Get("session", c)
	if err != nil {
		return []string{}, err
	}
	flashes, ok := sess.Values["flashes"].([]string)
    if !ok {
        return []string{}, nil
    }
    sess.Values["flashes"] = []string{}
    sess.Save(c.Request(), c.Response())
    return flashes, nil
}
// End: Helpers
// ==========================

// Main method
func main() {
	// Create app as an instance of Echo
	app := echo.New()

	// Add template-renderer to app
	// app.Renderer = TODO

	db, err := initDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()
	Db = db

	app.Use(session.Middleware(sessions.NewCookieStore(SECRET_KEY)))

	app.Use(middleware.StaticWithConfig(middleware.StaticConfig{
        Root: "static", // static folder
    }))

	setupRoutes(app)

	app.Logger.Fatal(app.Start(":8000"))
}
