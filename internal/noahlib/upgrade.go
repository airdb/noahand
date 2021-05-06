package noahlib

import (
	"airdb.io/airdb/sailor"
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"airdb.io/airdb/sailor/fileutil"
)

func DoSelfUpdate() {
	name := "noah_latest.zip"
	dl := DefaultDomain + "/release/" +  name

	log.Printf("download url: %s\n", dl)

	start := time.Now()
	client := sailor.NewHTTPClient()
	client.SetURL(dl)

	httpReader, err := client.Downloader()
	if err != nil {
		return
	}

	log.Printf("%s download successfully, cost: %s\n", name, time.Since(start))

	executable, err := os.Executable()
	if err != nil {
		log.Println("get_executable_fail")
		return
	}

	tmpPath := fmt.Sprintf("/tmp/%s.%d", filepath.Base(executable), time.Now().UnixNano())

	log.Printf("%s download successfully, cost: %s\n", executable, time.Since(start))

	defer httpReader.Close()
	gr, err := gzip.NewReader(httpReader)
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
	name := "plugin_greeter.so"
	dl := "https://github.com/airdb/noah/releases/latest/download/"  + name

	start := time.Now()
	client := sailor.NewHTTPClient()
	client.SetURL(dl)

	httpReader, err := client.Downloader()
	if err != nil {
		return
	}

	log.Printf("%s download successfully, cost: %s\n", name, time.Since(start))

	defer httpReader.Close()
	content, _ := ioutil.ReadAll(httpReader)

	err = fileutil.WriteFile(GetNoahPluginPath()+ name, string(content))
	log.Println(err)
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

	client := sailor.NewHTTPClient()
	client.SetURL(dl)

	start := time.Now()
	httpReader, err := client.Downloader()
	if err != nil {
		log.Println(err)

		return
	}

	log.Printf("%s download successfully, cost: %s\n", name, time.Since(start))


	defer httpReader.Close()
	// use tee first, then use `ioutil.NopCloser(buf)`
	buf := &bytes.Buffer{}
	tee := io.TeeReader(httpReader, buf)

	md5sum := fileutil.GetMd5Sum(tee)
	tmpWorkDir := path.Join(GetNoahModulePath(), md5sum)
	fileutil.EnsureFolderExists(tmpWorkDir)

	defer os.Chdir(GetWorkingPath())
	os.Chdir(tmpWorkDir)

	fileutil.ExtraTarFile(ioutil.NopCloser(buf))
}