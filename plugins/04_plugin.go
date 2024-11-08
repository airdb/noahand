package main

import (
	"guardhouse/pkg/oskit"
	"log"
)

func Plugin04() {
	log.SetFlags(log.Lshortfile)
	log.Println("Hello from Plugin 04!")

	oskit.InstallTarball()
}
