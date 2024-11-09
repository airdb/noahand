package main

import (
	"log"
)

type Plugin04 struct{}

var Plugin = Plugin04{}

func (p Plugin04) GetInfo() {
	log.Println("Plugin04 version: v0.1.0")
}

func (p Plugin04) Init() {
	log.Println("Init from Plugin 04!")
}

func (p Plugin04) Start() {
	log.Println("Hello from Plugin 04!")
	// oskit.DownloadAndInstallTarball()
}

func (p Plugin04) Stop() {
	log.Println("Stop from Plugin 04!")
}
