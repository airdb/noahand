package noahlib

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
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

	log.Printf("download url: %s\n", dl)

	start := time.Now()
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, dl, nil)
	if err != nil {
		log.Println(err)

		return
	}

	transport := http.DefaultTransport.(*http.Transport)
	var body io.ReadCloser
	resp, err := transport.RoundTrip(req)
	if err != nil {
		log.Println("download zip file fail, url:", dl)

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
	gr, err := gzip.NewReader(body)
	if err != nil {
		log.Println(err)

		return
	}

	tr := tar.NewReader(gr)
	for {
		hdr, err := tr.Next()
		if err != nil {
			log.Println("download_zip_fail")

			return
		}

		data, err := ioutil.ReadAll(tr)
		name := filepath.Base(executable) + "-" + runtime.GOOS
		// if hdr.Name == executable + runtime.GOOS {
		if strings.HasSuffix(hdr.Name, name) {
			log.Printf("start write file, name: %v, size: %v, tmpPath: %v\n", hdr.Name, hdr.Size, tmpPath)
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
	content, _ := ioutil.ReadAll(resp.Body)

	err = fileutil.WriteFile(GetNoahPluginPath()+mod, string(content))
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

func DownloadZip() {
	name := "modules.zip"
	dl := DefaultDomain + "/release/" + name

	log.Printf("download url: %s\n", dl)

	start := time.Now()
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, dl, nil)
	if err != nil {
		log.Println(err)

		return
	}

	transport := http.DefaultTransport.(*http.Transport)

	resp, err := transport.RoundTrip(req)
	if err != nil {
		log.Println("download zip file fail, url:", dl)

		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Println("download server unreachable")

		return
	}

	tmpPath := fmt.Sprintf("/tmp/%s.%d", name, time.Now().UnixNano())

	log.Printf("%s download successfully, cost: %s\n", tmpPath, time.Since(start))

	outFile, err := os.Create(tmpPath)
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		log.Println(err)
	}

	defer os.Remove(tmpPath)
	md5sum := GetFileMD5(tmpPath)

	ExtraTar(tmpPath, md5sum)
}

func GetFileMD5(filename string) string {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}
	log.Printf("%x", h.Sum(nil))

	return hex.EncodeToString(h.Sum(nil))
}
func ExtraTar(filename, md5sum string) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	fileutil.EnsureFolderExists(path.Join(GetNoahModulePath(), md5sum))


	uncompressedStream, err := gzip.NewReader(f)
	if err != nil {
		log.Fatal("ExtractTarGz: NewReader failed")
	}

	tarReader := tar.NewReader(uncompressedStream)

	for true {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("ExtractTarGz: Next() failed: %s", err.Error())
		}

		switch header.Typeflag {
		case tar.TypeDir:
			filename := path.Join(GetNoahModulePath(), md5sum, header.Name)
			fileutil.EnsureFolderExists(filename)
			if err := os.Mkdir(filename, 0755); err != nil {
				log.Fatalf("ExtractTarGz: Mkdir() failed: %s", err.Error())
			}
		case tar.TypeReg:
			filename := path.Join(GetNoahModulePath(), md5sum, header.Name)
			outFile, err := os.Create(filename)
			if err != nil {
				log.Fatalf("ExtractTarGz: Create() failed: %s", err.Error())
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				log.Fatalf("ExtractTarGz: Copy() failed: %s", err.Error())
			}
			outFile.Close()

		default:
			log.Fatalf(
				"ExtractTarGz: uknown type: %s in %s",
				header.Typeflag,
				header.Name)
		}
	}
}
