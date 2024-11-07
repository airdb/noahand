package main

import (
	"guardhouse/pkg/oskit"
	"log"
	"os"
	"path"
)

func Plugin04() {
	log.SetFlags(log.Lshortfile)
	log.Println("Hello from Plugin 04!")

	// Download a file
	surl := "http://oracle.airdb.host:8000/noah_latest.tgz"
	dpath := "/tmp/noah_latest.tgz"

	err := oskit.ResumableDownload(surl, dpath)
	if err != nil {
		log.Println("Download failed")
		return
	}

	log.Println("Download successful")

	// get file md5
	md5, err := oskit.GetFileMD5(dpath)
	if err != nil {
		log.Println("Failed to get file md5")
		return
	}

	unzipPath := path.Join("/tmp/noah", md5)

	os.MkdirAll(unzipPath, 0o755)

	oskit.ExtractTarball(dpath, unzipPath)
	log.Println("Unzipped to", unzipPath)
}
