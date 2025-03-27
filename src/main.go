package main

import (
	"fmt"
	"log"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"minitwit/src/datalayer"
	"minitwit/src/routes"
	"minitwit/src/template_rendering"
	"minitwit/src/handlers/helpers"
	"minitwit/src/handlers"
	"minitwit/src/models"
)

var SECRET_KEY = []byte("development key")

func main() {
	// Create app as an instance of Echo
	app := echo.New()

	// Add template-renderer to app
	app.Renderer = template_rendering.NewTemplateRenderer()

	db, err := datalayer.InitDB()
	if err != nil {
		fmt.Printf("initDB returned error: %v\n", err)
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	sqliteDb, err := datalayer.InitSqliteDB()
	if err != nil {
		fmt.Printf("initDB returned error: %v\n", err)
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer sqliteDb.Close()
	userRepo := datalayer.NewRepository[models.User](db, "users")
	userSqliteRepo := datalayer.NewSqliteRepository[models.User](sqliteDb, "user")
	messageRepo := datalayer.NewRepository[models.Message](db, "message")
	messageSqliteRepo := datalayer.NewSqliteRepository[models.Message](sqliteDb, "message")
	followerRepo := datalayer.NewRepository[models.Follower](db, "follower")
	followerSqliteRepo := datalayer.NewSqliteRepository[models.Follower](sqliteDb, "follower")
	handlers.SetUserRepo(userRepo, userSqliteRepo)
	handlers.SetMessageRepo(messageRepo, messageSqliteRepo)
	handlers.SetFollowerRepo(followerRepo, followerSqliteRepo)
	app.Use(session.Middleware(sessions.NewCookieStore(SECRET_KEY)))

	app.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root: "static", // static folder
	}))

	helpers.CreateLatestFile()

	routes.SetupRoutes(app, sqliteDb)

	// Custom error handler to log and expose internal errors
	app.HTTPErrorHandler = func(err error, c echo.Context) {
		log.Printf("❌ SERVER ERROR: %v", err)  // Log error
		c.JSON(500, map[string]string{"error": fmt.Sprintf("Server error: %v", err)})
	}

	app.Logger.Fatal(app.Start(":8000"))
}
