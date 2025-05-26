package main

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"minitwit/src/datalayer"
	"minitwit/src/handlers/repo_wrappers"
	"minitwit/src/metrics"
	"minitwit/src/routes"
	"minitwit/src/snapshots"
	"minitwit/src/template_rendering"
	"minitwit/src/utils"
)

var SECRET_KEY = []byte("development key")

func main() {
	// Set logging options
	utils.SetSlogDefaults()

	// Create app as an instance of Echo
	app := echo.New()

	// Add template-renderer to app
	app.Renderer = template_rendering.NewTemplateRenderer()

	db, err := datalayer.InitDB()
	if err != nil {
		utils.LogError("initDB returned an error", err)
		utils.LogFatal("Failed to initialize database", err)
	}
	defer db.Close()

	repo_wrappers.InitRepos(db)

	app.Use(session.Middleware(sessions.NewCookieStore(SECRET_KEY)))

	app.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root: "static", // static folder
	}))

	// Setup middleware for Prometheus
	app.Use(echoprometheus.NewMiddleware("minitwit"))

	metrics.Initialize()

	app.Use(metrics.PrometheusMiddleware()) // For metrics

	app.Use(middleware.SecureWithConfig(middleware.SecureConfig{
        XFrameOptions: 			"DENY", // prevents clickjacking
    }))

	app.Use(middleware.BodyLimit("2M")) // drop >2 MiB payloads early
	app.Use(middleware.RateLimiter(
		middleware.NewRateLimiterMemoryStore(100))) // 100 req/s per IP

	routes.SetupRoutes(app)

	snapshots.RecordSnapshots()

	srv := &http.Server{
		Addr:              ":8000",
		Handler:           app,
		ReadHeaderTimeout: utils.GetEnvDuration("READ_HEADER_TIMEOUT", "5s"),
		ReadTimeout:       utils.GetEnvDuration("READ_TIMEOUT", "10s"),
		WriteTimeout:      utils.GetEnvDuration("WRITE_TIMEOUT", "10s"),
		IdleTimeout:       utils.GetEnvDuration("IDLE_TIMEOUT", "60s"),
		MaxHeaderBytes:    utils.GetEnvInt("MAX_HEADER_BYTES", 1<<20),
	}

	app.Logger.Fatal(srv.ListenAndServe())
}
