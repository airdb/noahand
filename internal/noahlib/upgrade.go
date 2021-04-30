package noahlib

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"airdb.io/airdb/sailor/fileutil"
	"github.com/pkg/errors"
)

func DoSelfUpdate() {
	dl := DefaultDomain + "/release/noah_latest.zip"

	fmt.Printf("It will take about 1 minute for downloading.\nDownload url: %s\n", dl)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, dl, nil)
	if err != nil {
		log.Println(err)

		return
	}

	transport := http.DefaultTransport.(*http.Transport)
	var body io.ReadCloser
	resp, err := transport.RoundTrip(req)
	if resp != nil {
		defer resp.Body.Close()
		body = resp.Body
	}

	gr, err := gzip.NewReader(body)
	if err != nil {
		log.Println(err)
		return
	}

	executable, err := os.Executable()
	if err != nil {
		log.Println("get_executable_fail")
		return
	}
	fmt.Println(executable, err)

	tmpPath := fmt.Sprintf("/tmp/%s.%d", filepath.Base(executable), time.Now().UnixNano())
	tr := tar.NewReader(gr)
	for {
		hdr, err := tr.Next()
		if err != nil {
			log.Println("download_zip_fail")

			return
		}

		log.Printf("name: %v, size: %v, tmpPath: %v\n", hdr.Name, hdr.Size, tmpPath)

		data, err := ioutil.ReadAll(tr)
		name := filepath.Base(executable) +"-"+ runtime.GOOS
		// if hdr.Name == executable + runtime.GOOS {
		if strings.HasSuffix(hdr.Name,name) {
			err := ioutil.WriteFile(tmpPath, data, 0755) // #nosec
			if err != nil {

				return
			}
			// file written, quit listing loop.
			break
		}
		if err != nil {
			return
		}
	}

	defer os.Remove(tmpPath)
	err = exec.CommandContext(context.Background(), "/usr/bin/install", tmpPath, executable).Run()
	if err != nil {
		log.Println(err)
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