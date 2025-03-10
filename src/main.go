package main

import (
	"log"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"minitwit/src/datalayer"
	"minitwit/src/handlers/helpers"
	"minitwit/src/handlers/repo_wrappers"
	"minitwit/src/metrics"
	"minitwit/src/routes"
	"minitwit/src/template_rendering"
)

var SECRET_KEY = []byte("development key")

func main() {
	// Create app as an instance of Echo
	app := echo.New()

	// Add template-renderer to app
	app.Renderer = template_rendering.NewTemplateRenderer()

	db, err := datalayer.InitDB()
	if err != nil {
		log.Printf("initDB returned error: %v\n", err)
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()
	
	repo_wrappers.InitRepos(db)
	
	app.Use(session.Middleware(sessions.NewCookieStore(SECRET_KEY)))

	app.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root: "static", // static folder
	}))

	// Setup middleware for Prometheus
	app.Use(echoprometheus.NewMiddleware("minitwit"))

	helpers.CreateLatestFile()

	routes.SetupRoutes(app)

	app.Logger.Fatal(app.Start(":8000"))
}
