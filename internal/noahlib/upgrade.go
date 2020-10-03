package noahlib

import (
	"context"
	"fmt"
	"github.com/airdb/sailor"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"

	"github.com/minio/selfupdate"
)

func DoSelfUpdate() {
	dl := "https://github.com/airdb/noah/releases/latest/download/noah"
	if runtime.GOOS == "darwin" {
		dl = dl + "-" + runtime.GOOS
	}

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

	sailor.WriteFile(GetPluginPath() + mod, string(content))
}

func doRequest(dl string) (*http.Response, error){
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

	if resp.StatusCode  != http.StatusOK {
		return  nil, errors.Errorf("http status code is", resp.StatusCode)
	}

	return resp, nil
}
