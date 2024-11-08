package oskit

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Calculate the MD5 checksum of a file
func GetFileMD5(filePath string) (string, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Create a new MD5 hash object
	hash := md5.New()

	// Copy the file content into the hash object
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	// Get the final hash as a byte slice and convert it to a hex string
	hashInBytes := hash.Sum(nil)
	hashInString := hex.EncodeToString(hashInBytes)

	return hashInString, nil
}

func ExtractTarball(srcFile string, destPath string) {
	var err error

	ext := filepath.Ext(srcFile)

	switch ext {
	case ".zip":
		err = unzip(srcFile, destPath)
		if err != nil {
			log.Println("Failed to extract zip file:", err)
		}
	case ".tar.gz", ".tgz", ".tar":
		err = Untar(srcFile, destPath)
		if err != nil {
			log.Println("Failed to extract tarball:", err)
		}
	default:
		log.Println("Unknown file format")
	}
}

// Extract a .tar or .tar.gz file
func Untar(src string, dest string) error {
	log.Println("Untar", src, "to", dest)

	// Open the source file
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer file.Close()

	var fileReader io.Reader = file

	// Check if it's a .tar.gz file, and create a gzip reader if so
	if filepath.Ext(src) == ".gz" || filepath.Ext(src) == ".tgz" {
		gzipReader, err := gzip.NewReader(file)
		if err != nil {
			return err
		}
		defer gzipReader.Close()
		fileReader = gzipReader
	}

	// Create a tar reader
	tarReader := tar.NewReader(fileReader)
	log.Println("Start extract tarball")

	// Iterate through the tar archive
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// Define the target location
		target := filepath.Join(dest, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			// Create directory if it doesn't exist
			if err := os.MkdirAll(target, os.FileMode(header.Mode)); err != nil {
				return err
			}
		case tar.TypeReg:
			// Ensure the directory exists before creating the file
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				return err
			}

			// Create the file and defer its closure
			outFile, err := os.Create(target)
			if err != nil {
				return err
			}
			defer outFile.Close()

			// Copy the file content
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return err
			}

		default:
			// Handle other file types if necessary (e.g., symlinks)
		}
	}

	return nil
}

func unzip(src string, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", fpath)
		}

		if f.FileInfo().IsDir() {
			_ = os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}
	return nil
}
