package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"guardhouse/internal/version"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func DefaultRoot(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"deploy_info": version.GetBuildInfo(),
	}
	json.NewEncoder(w).Encode(response)
}

func Run() {
	fmt.Printf("Chi server starting")

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Serve static files
	fileServer := http.FileServer(http.Dir("./release"))
	r.Handle("/release/*", http.StripPrefix("/release/", fileServer))

	r.Get("/", DefaultRoot)
	r.Get("/host", DefaultRoot)

	defaultPort := ":80"
	http.ListenAndServe(defaultPort, r)
}
