package api

import (
	"fmt"

	"github.com/airdb/sailor/faas"
	"github.com/go-chi/chi"
)

func Run() {
	fmt.Printf("Gin start")

	r := chi.NewRouter()
	r.GET("/", "")

	faas.RunTencentChi(r)
}
