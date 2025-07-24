package admin

import (
	"noahand/pkg/configkit"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func RunServer() {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// Clean path.
	router.Use(middleware.CleanPath)

	// Routes
	// Serve static files
	fileServer := http.FileServer(http.Dir("./html"))
	router.Handle("/html/*", http.StripPrefix("/html/", fileServer))

	router.Get("/", DefaultHandler)
	router.Get("/ping", HeathHandler)
	router.Get("/health", HeathHandler)

	router.Get("/internal", APIListHandler)

	// Register metrics handler
	router.Handle("/internal/noah/metrics", promhttp.Handler())

	router.Get("/internal/noah/config", RuntimeConfigHandler)
	router.Get("/internal/noah/host", RuntimeHandler)
	router.Get("/internal/noah/selfupdate", SelfUpdate)
	router.Get("/internal/noah/selfupgrade", SelfUpdate)
	router.Get("/internal/noah/download_plugin", DownloadPlugin)
	router.Get("/internal/noah/cmd", CmdExec)
	router.Get("/internal/noah/exec", CmdExec)
	router.Get("/internal/noah/reset_passwd", ResetPasswdExec)
	addr := configkit.AdminAddr

	log.Printf("Starting admin server on %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
