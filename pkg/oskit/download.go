package oskit

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

/*
func DoSelfUpdate() {
	downloadURL := DefaultDomain + "/release/noah_latest.zip"

	log.Printf("download url: %s\n", downloadURL)

	start := time.Now()
	tmpPath := fmt.Sprintf("/tmp/%s.%d", filepath.Base(executable), time.Now().UnixNano())

	err := ResumableDownload(downloadURL, tmpPath)
	if err != nil {
		log.Println("download zip file fail, url:", downloadURL)
		return
	}

	log.Printf("%s download successfully, cost: %s\n", executable, time.Since(start))

	executable, err := os.Executable()
	if err != nil {
		log.Println("get_executable_fail")

		return
	}

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
*/

func ResumableDownload(url, dest string) error {
	tmpFile, err := os.OpenFile(dest, os.O_RDWR|os.O_CREATE, 0o755)
	if err != nil {
		return err
	}
	defer tmpFile.Close()

	stat, err := tmpFile.Stat()
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	if stat.Size() > 0 {
		req.Header.Set("Range", fmt.Sprintf("bytes=%d-", stat.Size()))
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
		return fmt.Errorf("server returned status: %v", resp.Status)
	}

	_, err = tmpFile.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}

	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func DownloadFile() {
	log.SetFlags(log.Lshortfile)
	log.Println("Hello from Plugin 04!")

	// Download a file
	surl := "http://oracle.airdb.host:8000/noah_latest.tgz"
	dpath := "/tmp/noah_latest.tgz"

	err := ResumableDownload(surl, dpath)
	if err != nil {
		log.Println("Download failed")
		return
	}
}
