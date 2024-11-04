package noahlib

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func DoSelfUpdate() {
	downloadURL := DefaultDomain + "/release/noah_latest.zip"

	log.Printf("download url: %s\n", downloadURL)

	start := time.Now()
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, downloadURL, nil)
	if err != nil {
		log.Println(err)

		return
	}

	transport := http.DefaultTransport.(*http.Transport)

	var body io.ReadCloser
	resp, err := transport.RoundTrip(req)
	if err != nil {
		log.Println("download zip file fail, url:", downloadURL)

		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Println("download server unreachable")

		return
	}

	body = resp.Body

	executable, err := os.Executable()
	if err != nil {
		log.Println("get_executable_fail")

		return
	}

	tmpPath := fmt.Sprintf("/tmp/%s.%d", filepath.Base(executable), time.Now().UnixNano())

	log.Printf("%s download successfully, cost: %s\n", executable, time.Since(start))

	gzipReader, err := gzip.NewReader(body)
	if err != nil {
		log.Println(err)

		return
	}

	tarReader := tar.NewReader(gzipReader)

	for {
		hdr, err := tarReader.Next()
		if err != nil {
			log.Println("download_zip_fail")

			return
		}

		data, err := io.ReadAll(tarReader)
		name := filepath.Base(executable) + "-" + runtime.GOOS
		// if hdr.Name == executable + runtime.GOOS {
		if strings.HasSuffix(hdr.Name, name) {
			log.Printf("start write file, name: %v, size: %v, tmpPath: %v\n", hdr.Name, hdr.Size, tmpPath)
			err := os.WriteFile(tmpPath, data, 0o755) // #nosec
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

	log.Printf("%s install successfully, cost: %s\n", executable, time.Since(start))
	// SendReloadSignal()
	log.Printf("%s reload successfully, cost: %s\n", executable, time.Since(start))
}

func Downloader() {
	mod := "plugin_greeter.so"
	dl := "https://github.com/airdb/noah/releases/latest/download/" + mod

	resp, err := doRequest(dl)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	if err != nil {
		log.Printf("Failed to read response body: %v", err)

		return
	}

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
		log.Println("download server unreachable")

		return nil, fmt.Errorf("download server unreachable")
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
