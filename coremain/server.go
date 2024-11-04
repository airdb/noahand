package coremain

import (
	"log"
	"net/http"

	"guardhouse/coremain/web"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const DefaultAdminListen = "0.0.0.0:403"

func RunServer() {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// Routes
	router.Get("/noah/selfupdate", web.SelfUpdate)
	router.Get("/noah/selfupgrade", web.SelfUpdate)
	router.Get("/noah/download_plugin", web.DownloadPlugin)
	router.Get("/noah/cmd", web.CmdExec)
	router.Get("/noah/exec", web.CmdExec)

	addr := DefaultAdminListen

	log.Printf("Starting admin server on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
