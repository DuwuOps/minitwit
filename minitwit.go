package main

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"io"
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
var PER_PAGE  = 30
var Db *sql.DB

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

func initDB() (*sql.DB, error) {
	// Create tmp directory
	dir := filepath.Dir(DATABASE)
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
	sqlFile, err := os.ReadFile("./schema.sql")
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
}

func populateDb(db *sql.DB, sqlFilePath string) error {
	query, err := os.ReadFile(sqlFilePath)
	if err != nil {
		fmt.Printf("os.ReadFile returned error: %v\n", err)
		return fmt.Errorf("failed to read file: %w", err)
	}
	queryString := string(query)
	_, err = db.Exec(queryString)
	if err != nil {
		fmt.Printf("sql.Open returned error: %v\n", err)
		db.Close()
		return fmt.Errorf("failed to execute query: %w", err)
	}
	return nil
}

// Note: Method signature and return type have been modified
func queryDbSingle(db *sql.DB, query string, args ...any)  *sql.Row {
	row := db.QueryRow(query, args...)
	return row
}

func queryDB(db *sql.DB, query string, args ...any) (*sql.Rows, error) {
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

func rowsToMapList(rows *sql.Rows) ([]map[string]interface{}, error) {
    var result []map[string]interface{}
    cols, _ := rows.Columns()

    for rows.Next() {
        // Create a slice of interface{}'s to represent each column,
        // and a second slice to contain pointers to each item in the columns slice.
        columns := make([]interface{}, len(cols))
        columnPointers := make([]interface{}, len(cols))
        for i, _ := range columns {
            columnPointers[i] = &columns[i]
        }
        
        // Scan the result into the column pointers...
        if err := rows.Scan(columnPointers...); err != nil {
            fmt.Printf("rows.Scan returned error: %v\n", err)
            return nil, err
        }

        // Create our map, and retrieve the value for each column from the pointers slice,
        // storing it in the map with the name of the column as the key.
        m := make(map[string]interface{})
        for i, colName := range cols {
            val := columnPointers[i].(*interface{})
            m[colName] = *val
        }
        
        // Outputs: map[columnName:value columnName2:value2 columnName3:value3 ...] 
        fmt.Print(m)
        result = append(result, m)
    }

    return result, nil
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

	app.Static("/static", "static")
}
// End: Routes
// ==========================

// ==========================
// Start: Route-Handlers

// Shows a users timeline or if no user is logged in it will
// redirect to the public timeline.  This timeline shows the user's
// messages as well as all the messages of followed users.
func Timeline(c echo.Context) error {
	log.Println("User entered Timeline via route \"/\"")
    log.Printf("We got a visitor from: %s", c.Request().RemoteAddr)
    loggedIn, _ := isUserLoggedIn(c)
    if !loggedIn {
        return c.Redirect(http.StatusFound, "/public")
    }

    sessionUserId, _ := getSessionUserID(c)
    rows, err := queryDB(Db, `select message.*, user.* from message, user
                          where message.flagged = 0 and message.author_id = user.user_id and (
                              user.user_id = ? or
                              user.user_id in (select whom_id from follower
                                                      where who_id = ?))
                          order by message.pub_date desc limit ?`, 
                          sessionUserId, sessionUserId, PER_PAGE,
                        )
    
    if err != nil {
        fmt.Printf("Timeline: queryDB returned error: %v\n", err)
        return err
    }

    msgs, err := rowsToMapList(rows)
    if err != nil {
		fmt.Printf("rowsToMapList returned error: %v\n", err)
		return err
	}

	user, err := getCurrentUser(c)
	if err != nil {
		fmt.Printf("No user found. getCurrentUser returned error: %v\n", err)
	}

    data := map[string]interface{}{
		"Messages": msgs,
		"User": user,
    }
    return c.Render(http.StatusOK, "timeline.html", data)
}

func PublicTimeline(c echo.Context) error {
	log.Println("User entered PublicTimeline via route \"/public\"")
	
    rows, err := queryDB(Db, `select message.*, user.* from message, user
                            where message.flagged = 0 and message.author_id = user.user_id
                            order by message.pub_date desc limit ?`, 
                            PER_PAGE,
                        )
	if err != nil {
		fmt.Printf("PublicTimeline: queryDB returned error: %v\n", err)
		return err
	}

    
	msgs, err := rowsToMapList(rows)
    if err != nil {
        fmt.Printf("rowsToMapList returned error: %v\n", err)
        return err
    }

    data := map[string]interface{}{
		"Messages": msgs,
    }
    return c.Render(http.StatusOK, "timeline.html", data)
}

// Display's a users tweets.
func UserTimeline(c echo.Context) error {
	username := c.Param("username")
	fmt.Printf("User entered UserTimeline via route \"/:username\" as \"/%v\"\n", username)

	row := queryDbSingle(Db, "select user_id from user where username = ?", username)
	var requestedUserId int
	err := row.Scan(&requestedUserId)
	if err != nil {
		fmt.Printf("row.Scan returned error: %v\n", err)
		c.String(http.StatusNotFound, "Not found")
	}

	followed := false
	loggedIn, _ := isUserLoggedIn(c)
	if loggedIn {
		sessionUserId, _ := getSessionUserID(c)
		follow_result := queryDbSingle(Db, `select 1 from follower where
             follower.who_id = ? and follower.whom_id = ?`,
			sessionUserId, requestedUserId)
		
		// The query should return a 1, if the user follows the user of the timeline.
		var result int
		err := follow_result.Scan(&result)
		followed = err != nil
	}

	rows, err := queryDB(Db, `select message.*, user.* from message, user where
                            user.user_id = message.author_id and user.user_id = ?
                            order by message.pub_date desc limit ?`,
							requestedUserId, PER_PAGE,
	)

	if err != nil {
		fmt.Printf("UserTimeline: queryDB returned error: %v\n", err)
		return err
	}

	msgs, err := rowsToMapList(rows)
    if err != nil {
		fmt.Printf("rowsToMapList returned error: %v\n", err)
		return err
	}

	user, err := getCurrentUser(c)
	if err != nil {
		fmt.Printf("No user found. getCurrentUser returned error: %v\n", err)
	}

	data := map[string]interface{}{
		"Messages":    msgs,
		"Followed":    followed,
		"ProfileUser": followed,
		"User": user,
	}
	return c.Render(http.StatusOK, "timeline.html", data)
}

func FollowUser(c echo.Context) error {
	log.Println("User entered FollowUser via route \"/:username/follow\"")
	return errors.New("Not implemented yet") //TODO
}

func UnfollowUser(c echo.Context) error {
	log.Println("User entered UnfollowUser via route \"/:username/unfollow\"")
	return errors.New("Not implemented yet") //TODO
}

func AddMessage(c echo.Context) error {
	log.Println("User entered AddMessage via route \"/add_message\"")
	return errors.New("Not implemented yet") //TODO
}

// Logs the user in.
func Login(c echo.Context) error {
	log.Println("User entered Login via route \"/login\"")
    loggedIn, _ := isUserLoggedIn(c)
    if loggedIn {
        return c.Redirect(http.StatusFound, "/")
    }

	var dbUser user

    var errorMessage string
    if c.Request().Method == http.MethodPost {
        username := c.FormValue("username")
        password := c.FormValue("password")

		dbUser.Username = username

        err := Db.QueryRow(`
            SELECT user_id, pw_hash FROM user
            WHERE username = ?
        `, username).Scan(&dbUser.UserID, &dbUser.PwHash)

        if errors.Is(err, sql.ErrNoRows) {
            errorMessage = "Invalid username"
        } else if err != nil {
            fmt.Printf("Db.QueryRow returned error: %v\n", err)
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

	flashes, _ := getFlashes(c)

    data := map[string]interface{}{
		"Error":   errorMessage,
		"Flashes": flashes,
		"User": dbUser,
    }
    return c.Render(http.StatusOK, "login.html", data)
}

func Register(c echo.Context) error {
	log.Println("User entered Register via route \"/register\"")
    loggedIn, _ := isUserLoggedIn(c)
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
            existingID, _ := getUserId(username)
            if existingID != 0 {
                errorMessage = "The username is already taken"
            } else {
                hash, err := generatePasswordHash(password)
                if err != nil {
                    fmt.Printf("generatePasswordHash returned error: %v\n", err)
                    return err
                }
                _, err = Db.Exec(`
                    INSERT INTO user (username, email, pw_hash)
                    VALUES (?, ?, ?)
                `, username, email, hash)
                if err != nil {
                    fmt.Printf("Db.Exec returned error: %v\n", err)
                    return err
                }

                addFlash(c, "You were successfully registered and can login now")
                return c.Redirect(http.StatusFound, "/login")
            }
        }
    }

	flashes, _ := getFlashes(c)

    data := map[string]interface{}{
		"Error":   errorMessage,
		"Flashes": flashes,
    }
    return c.Render(http.StatusOK, "register.html", data)
}

func Logout(c echo.Context) error {
	log.Println("User entered Logout via route \"/logout\"")
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
		fmt.Printf("session.Get returned error: %v\n", err)
		return err
	}
    sess.Values["user_id"] = userID
    sess.Save(c.Request(), c.Response())
	return nil
}

func getSessionUserID(c echo.Context) (int, error) {
	sess, err := session.Get("session", c)
	if err != nil {
		fmt.Printf("session.Get returned error: %v\n", err)
		return 0, err
	}
    id, _ := sess.Values["user_id"].(int)
	return id, nil
}

func clearSessionUserID(c echo.Context) error {
    sess, err := session.Get("session", c)
	if err != nil {
		fmt.Printf("session.Get returned error: %v\n", err)
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
        fmt.Printf("Db.QueryRow returned error: %v\n", err)
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
		fmt.Printf("session.Get returned error: %v\n", err)
		return err
	}
	flashes, ok := sess.Values["Flashes"].([]string)
    if !ok {
        flashes = []string{}
    }
    flashes = append(flashes, message)
	sess.Values["Flashes"] = flashes
    sess.Save(c.Request(), c.Response())
	return nil
}

// Takes a context
// Returns empties the flashes in the given context and returns the flashes in a list of strings
func getFlashes(c echo.Context) ([]string, error) {
	sess, err := session.Get("session", c)
	if err != nil {
		fmt.Printf("session.Get returned error: %v\n", err)
		return []string{}, err
	}
	flashes, ok := sess.Values["Flashes"].([]string)
    if !ok {
        return []string{}, nil
    }
	sess.Values["Flashes"] = []string{}
    sess.Save(c.Request(), c.Response())
    return flashes, nil
}

type user struct {
    UserID int
	Username string
	Email string
	PwHash string
}

// Takes a context
// Returns the current user
func getCurrentUser(c echo.Context) (*user, error) {
	id, err := getSessionUserID(c)
	var user user

	if err != nil {
		fmt.Printf("getSessionUserID returned error: %v\n", err)
		return nil, err
	}

	rows := queryDbSingle(Db, "select * from user where user_id = ?",
						id,
					)
	
	err = rows.Scan(&user.UserID, &user.Username, &user.Email, &user.PwHash)
	if err != nil {
		fmt.Printf("rows.Scan returned error: %v\n", err)
		return nil, err
	}
	fmt.Printf("Found user in database! %v\n", user)
	fmt.Printf("user.UserID: %v\n", user.UserID)
	fmt.Printf("user.Username: %v\n", user.Username)
	fmt.Printf("user.Email: %v\n", user.Email)
	fmt.Printf("user.PwHash: %v\n", user.PwHash)
	return &user, nil
}
// End: Helpers
// ==========================


// ==========================
// Begin: Template Rendering

// Implementation of echo.Renderer interface
type TemplateRenderer struct {
    templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
    tmpl := template.Must(t.templates.Clone())
    tmpl = template.Must(tmpl.ParseFiles(filepath.Join("templates", name)))
    return tmpl.ExecuteTemplate(w, name, data)
}

// Create and return a new instance of a TemplateRenderer
func NewTemplateRenderer() *TemplateRenderer {
	tmpl := template.Must(template.ParseGlob(filepath.Join("templates", "*.html")))
    return &TemplateRenderer{
        templates: tmpl,
    }
}

// End: Template Rendering
// ==========================


// Main method
func main() {
	// Create app as an instance of Echo
	app := echo.New()

	// Add template-renderer to app
	app.Renderer = NewTemplateRenderer()

	db, err := initDB()
	if err != nil {
		fmt.Printf("initDB returned error: %v\n", err)
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()
	
	populateDb(db, "./tmp/generate_data.sql")
	Db = db

	app.Use(session.Middleware(sessions.NewCookieStore(SECRET_KEY)))

	app.Use(middleware.StaticWithConfig(middleware.StaticConfig{
        Root: "static", // static folder
    }))

	setupRoutes(app)

	app.Logger.Fatal(app.Start(":8000"))
}
