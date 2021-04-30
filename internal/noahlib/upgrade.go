package noahlib

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"

	"airdb.io/airdb/sailor/fileutil"
	"github.com/minio/selfupdate"
	"github.com/pkg/errors"
)

func DoSelfUpdate() {
	/*
	dl := "https://github.com/airdb/noah/releases/latest/download/noah"
	if runtime.GOOS == "darwin" {
		dl = dl + "-" + runtime.GOOS
	}
	 */
	dl := "http://sg.airdb.host/release/noah-" + runtime.GOOS

	fmt.Printf("It will take about 1 minute for downloading.\nDownload url: %s\n", dl)

	client := &http.Client{}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, dl, nil)
	if err != nil {
		log.Println(err)

		return
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)

		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Println("download failed!")

		return
	}

	defer resp.Body.Close()

	err = selfupdate.Apply(resp.Body, selfupdate.Options{})
	if err != nil {
		log.Println("update failed!")
	} else {
		log.Println("update successfully!")
	}
}

func Downloader() {
	mod := "plugin_greeter.so"
	dl := "https://github.com/airdb/noah/releases/latest/download/" + mod

	resp, err := doRequest(dl)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	content, _ := ioutil.ReadAll(resp.Body)

	err = fileutil.WriteFile(GetPluginPath()+mod, string(content))
	log.Println(err)
}

func doRequest(dl string) (*http.Response, error) {
	client := &http.Client{}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, dl, nil)
	if err != nil {
		log.Println(err)

		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)

		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("http status code is %v", resp.StatusCode)
	}

	return resp, nil
}

func InstallProcess() {
	tmpPath := "/sbin/noah"
	executable := "/tmp/noah_latest"
	defer os.Remove(tmpPath)
	err := exec.CommandContext(context.Background(), "/usr/bin/install", tmpPath, executable).Run()
	if err != nil {
		log.Println(err)
	}
}

func SendReloadSignal() error {
	ppid := strconv.Itoa(os.Getppid())
	err := exec.CommandContext(context.Background(), "/bin/kill", "-HUP", ppid).Run()

	return err
}