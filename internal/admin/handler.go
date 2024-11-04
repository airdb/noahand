package admin

import (
	"encoding/json"
	"guardhouse/internal/noahlib"
	"guardhouse/internal/version"
	"net/http"
)

func DefaultRoot(w http.ResponseWriter, _ *http.Request) {
	response := map[string]any{
		"deploy_info": version.GetBuildInfo(),
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func SelfUpdate(w http.ResponseWriter, _ *http.Request) {
	noahlib.DoSelfUpdate()
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("upgrade successfully"))
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func DownloadPlugin(w http.ResponseWriter, _ *http.Request) {
	noahlib.Downloader()
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("upgrade plugin successfully"))
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

type CmdReq struct {
	Cmd  string   `json:"cmd"`
	Args []string `json:"args"`
}

func CmdExec(w http.ResponseWriter, r *http.Request) {
	var req CmdReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusForbidden)
		_, err := w.Write([]byte("cmd not exists"))
		if err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
			return
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("exec command successfully"))
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}
