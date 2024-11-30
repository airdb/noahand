package admin

import (
	"encoding/json"
	"guardhouse/pkg/configkit"
	"guardhouse/pkg/version"
	"math/rand"
	"net/http"
	"os/exec"
	"time"

	"github.com/go-chi/render"
)

func DefaultHandler(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte("pong\n"))
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func HeathHandler(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte("ok\n"))
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func APIListHandler(w http.ResponseWriter, _ *http.Request) {
	msg := "api list:\n"

	for _, api := range configkit.GetConfig().AdminApiList {
		msg += api + "\n"
	}

	_, err := w.Write([]byte(msg))
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func RuntimeConfigHandler(w http.ResponseWriter, r *http.Request) {
	config := configkit.GetConfig()

	render.JSON(w, r, config)
}

func RuntimeHandler(w http.ResponseWriter, _ *http.Request) {
	response := map[string]any{
		"deploy_info": version.GetBuildInfo(),
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func SelfUpdate(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("upgrade successfully"))
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func DownloadPlugin(w http.ResponseWriter, _ *http.Request) {
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

func ResetPasswdExec(w http.ResponseWriter, _ *http.Request) {
	cmd := &CmdReq{}

	// Generate random password
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	raw := make([]byte, 8)
	for i := range raw {
		raw[i] = charset[seededRand.Intn(len(charset))]
	}

	password := string(raw)

	cmd.Cmd = "echo"
	cmd.Args = []string{password, "|", "passwd", "--stdin", "root"}

	command := exec.Command(cmd.Cmd, cmd.Args...)
	command.Run()

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("reset passwd successfully, password: " + password))
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}
