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
Rules:
- 1. Download file must be tarball, and file extension must be .tgz or .tar.gz, zip or tar.
- 2. Download file must be a valid file, not empty.
- 3. Download file must have a md5 download file. (example, tarbal is xxx.tgz, md5 file name is xxx.tgz.md5, and must keep at the same directory)
*/

func IsEmptyRemoteFile(surl string) (flag bool) {
	// Use HEAD request to check file size and ensure content is available
	headReq, err := http.NewRequestWithContext(context.Background(), http.MethodHead, surl, nil)
	if err != nil {
		flag = true

		return
	}

	client := &http.Client{}
	headResp, err := client.Do(headReq)
	if err != nil {
		flag = true

		return
	}
	defer headResp.Body.Close()

	// Check Content-Length to ensure file is not empty
	if headResp.ContentLength <= 0 {
		log.Println("file is empty or server does not provide content length")
		flag = true

		return
	}

	// Default flag is false, which means file is not empty.
	return
}

func ResumableDownload(surl, dest string) error {
	tmpFile, err := os.OpenFile(dest, os.O_RDWR|os.O_CREATE, 0o755)
	if err != nil {
		return err
	}
	defer tmpFile.Close()

	stat, err := tmpFile.Stat()
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, surl, nil)
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

// Download https://xxx/noah_latest.tgz.md5 file， read md5 value from body.
func DownloadMd5File(surl string) (string, error) {
	resp, err := http.Get(surl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	md5, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(md5), nil
}

// Get download filesize by http content-length.
func GetDownloadFileSize(surl string) (int64, error) {
	resp, err := http.Head(surl)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	return resp.ContentLength, nil
}
