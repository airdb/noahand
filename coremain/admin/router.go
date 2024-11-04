package admin

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const DefaultAdminListen = "0.0.0.0:403"

func RunServer() {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// Routes
	// Serve static files
	fileServer := http.FileServer(http.Dir("./release"))
	router.Handle("/release/*", http.StripPrefix("/release/", fileServer))

	router.Get("/", DefaultRoot)
	router.Get("/host", DefaultRoot)

	router.Get("/noah/selfupdate", SelfUpdate)
	router.Get("/noah/selfupgrade", SelfUpdate)
	router.Get("/noah/download_plugin", DownloadPlugin)
	router.Get("/noah/cmd", CmdExec)
	router.Get("/noah/exec", CmdExec)

	addr := DefaultAdminListen

	log.Printf("Starting admin server on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
