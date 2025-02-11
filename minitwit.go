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

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"

	"github.com/labstack/echo/v4/middleware"
)

// configuration
var (
    DATABASE  = "./tmp/minitwit.db"
    PER_PAGE  = 30
    SECRET_KEY = []byte("development key") // to parallel the Python "SECRET_KEY"
)
var Db *sql.DB

func main() {
	app := echo.New()

	app.Renderer = NewTemplateRenderer()

	// initDB
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

func connectDB() (*sql.DB, error) {
	//Returns a new connection to the database.

	dir := filepath.Dir(DATABASE)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	db, err := sql.Open("sqlite3", DATABASE)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	return db, nil
}

func initDB() (*sql.DB, error) {
	//Creates the database tables.
	db, err := connectDB()
	if err != nil {
		return nil, err
	}

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

    app.GET("/login", Login)
    app.POST("/login", Login)

    app.GET("/register", Register)
    app.POST("/register", Register)

    app.GET("/logout", Logout)

    app.POST("/add_message", AddMessage)

    app.GET("/:username", UserTimeline)

    app.GET("/:username/follow", FollowUser)
    app.GET("/:username/unfollow", UnfollowUser)
}
// End: Routes
// ==========================

// ==========================
// Start: Handlers
func Timeline(c echo.Context) error {
	return errors.New("Not implemented yet") //TODO
}

func PublicTimeline(c echo.Context) error {
	return errors.New("Not implemented yet") //TODO
}

func Login(c echo.Context) error {
    _, loggedIn := getSessionUserID(c)
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

    data := map[string]interface{}{
        "error": errorMessage,
        "flashes": getFlashes(c),
    }
    return c.Render(http.StatusOK, "login.html", data)
}

func Register(c echo.Context) error {
    _, loggedIn := getSessionUserID(c)
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

    data := map[string]interface{}{
        "error":   errorMessage,
        "flashes": getFlashes(c),
    }
    return c.Render(http.StatusOK, "register.html", data)
}

func Logout(c echo.Context) error {
    addFlash(c, "You were logged out")
    clearSessionUserID(c)
    return c.Redirect(http.StatusFound, "/public")
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
// End: Routes
// ==========================

// ==========================
// Start: Helpers
func generatePasswordHash(password string) (string, error) {
    hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(hashBytes), err
}

func checkPasswordHash(hashedPassword, plainPassword string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
    return err == nil
}

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

func addFlash(c echo.Context, message string) {
    sess, _ := session.Get("session", c)
    flashes, ok := sess.Values["flashes"].([]string)
    if !ok {
        flashes = []string{}
    }
    flashes = append(flashes, message)
    sess.Values["flashes"] = flashes
    sess.Save(c.Request(), c.Response())
}

func getFlashes(c echo.Context) []string {
    sess, _ := session.Get("session", c)
    flashes, ok := sess.Values["flashes"].([]string)
    if !ok {
        return []string{}
    }
    sess.Values["flashes"] = []string{}
    sess.Save(c.Request(), c.Response())
    return flashes
}

func getSessionUserID(c echo.Context) (int, bool) {
    sess, _ := session.Get("session", c)
    userID, ok := sess.Values["user_id"].(int)
    return userID, ok
}

func setSessionUserID(c echo.Context, userID int) {
    sess, _ := session.Get("session", c)
    sess.Values["user_id"] = userID
    sess.Save(c.Request(), c.Response())
}

func clearSessionUserID(c echo.Context) {
    sess, _ := session.Get("session", c)
    delete(sess.Values, "user_id")
    sess.Save(c.Request(), c.Response())
}
// End: Helpers
// ==========================