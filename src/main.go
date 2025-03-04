package main

import (
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"

	"minitwit/datalayer"
	"minitwit/routes"
	"minitwit/template"
)

var SECRET_KEY = []byte("development key")

func main() {
	// Create app as an instance of Echo
	app := echo.New()

	// Add template-renderer to app
	app.Renderer = template.NewTemplateRenderer()

	db, err := datalayer.InitDB()
	if err != nil {
		fmt.Printf("initDB returned error: %v\n", err)
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	datalayer.PopulateDb(db, "./queries/generate_data.sql")

	app.Use(session.Middleware(sessions.NewCookieStore(SECRET_KEY)))

	app.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root: "static", // static folder
	}))

	routes.SetupRoutes(app, db)

	app.Logger.Fatal(app.Start(":8000"))
}
