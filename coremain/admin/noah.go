package admin

import (
	"encoding/json"
	"net/http"

	"github.com/airdb/noah/internal/noahlib"
	"github.com/airdb/sailor/osutil"
	"github.com/airdb/sailor/version"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func RunWeb() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", Version)
	r.Get("/admin/selfupdate", SelfUpdate)
	r.Get("/admin/selfupgrade", SelfUpdate)
	r.Get("/admin/download_plugin", DownloadPlugin)
	r.Get("/admin/cmd", CmdExec)

	http.ListenAndServe(":403", r)
}

func Version(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(version.GetBuildInfo().ToString()))
	w.WriteHeader(http.StatusOK)
}

func SelfUpdate(w http.ResponseWriter, r *http.Request) {
	noahlib.DoSelfUpdate()
	w.Write([]byte("upgrade successfully"))
	w.WriteHeader(http.StatusOK)
}

func DownloadPlugin(w http.ResponseWriter, r *http.Request) {
	noahlib.Downloader()
	w.Write([]byte("upgrade plugin successfully"))
	w.WriteHeader(http.StatusOK)
}

type CmdReq struct {
	Cmd  string   `form:"cmd"`
	Args []string `form:"args"`
}

func CmdExec(w http.ResponseWriter, r *http.Request) {
	var req CmdReq

	if r.ContentLength == 0 {
		w.Write([]byte("invalid request"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.Write([]byte("invalid request"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ret, err := osutil.ExecCommand(req.Cmd, req.Args)
	if err != nil {
		w.Write([]byte("exec command failed"))
		w.WriteHeader(http.StatusOK)

		return
	}

	w.Write([]byte(ret))
	w.WriteHeader(http.StatusOK)
}
