package admin

import (
	"encoding/json"
	"net/http"

	"guardhouse/internal/noahlib"
	"guardhouse/internal/version"
)

func DefaultRoot(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"deploy_info": version.GetBuildInfo(),
	}
	json.NewEncoder(w).Encode(response)
}

func SelfUpdate(w http.ResponseWriter, r *http.Request) {
	noahlib.DoSelfUpdate()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("upgrade successfully"))
}

func DownloadPlugin(w http.ResponseWriter, r *http.Request) {
	noahlib.Downloader()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("upgrade plugin successfully"))
}

type CmdReq struct {
	Cmd  string   `json:"cmd"`
	Args []string `json:"args"`
}

func CmdExec(w http.ResponseWriter, r *http.Request) {
	var req CmdReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("cmd not exists"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("exec command successfully"))
}
