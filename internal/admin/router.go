package admin

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const DefaultAdminListen = "0.0.0.0:403"

func RunServer() {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// Routes
	// Serve static files
	fileServer := http.FileServer(http.Dir("./html"))
	router.Handle("/html/*", http.StripPrefix("/html/", fileServer))

	router.Get("/", DefaultHandler)
	router.Get("/ping", DefaultHandler)

	router.Get("/internal", APIListHandler)
	router.Get("/internal/", APIListHandler)

	// Register metrics handler
	router.Handle("/internal/noah/metrics", promhttp.Handler())

	router.Get("/internal/noah/config", RuntimeConfigHandler)
	router.Get("/internal/noah/host", RuntimeHandler)
	router.Get("/internal/noah/selfupdate", SelfUpdate)
	router.Get("/internal/noah/selfupgrade", SelfUpdate)
	router.Get("/internal/noah/download_plugin", DownloadPlugin)
	router.Get("/internal/noah/cmd", CmdExec)
	router.Get("/internal/noah/exec", CmdExec)

	addr := DefaultAdminListen

	log.Printf("Starting admin server on %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
