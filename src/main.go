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
	userRepo := datalayer.NewRepository[models.User](db, "user")
	handlers.SetUserRepo(userRepo)
	app.Use(session.Middleware(sessions.NewCookieStore(SECRET_KEY)))

	app.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root: "static", // static folder
	}))

	helpers.CreateLatestFile()

	routes.SetupRoutes(app, db)

	app.Logger.Fatal(app.Start(":8000"))
}
