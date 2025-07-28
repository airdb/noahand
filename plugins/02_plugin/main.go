// pluginA.go
package main

import (
	"noahand/internal/noahlib"
	"log"

	"github.com/go-resty/resty/v2"
)

type Plugin02 struct{}

var Plugin = Plugin02{}

func (p Plugin02) GetInfo() {
	log.Println("Plugin02 version: v0.1.0")
}

func (p Plugin02) Init() {
	log.Println("Init from Plugin 02!")
}

func (p Plugin02) Start() {
	log.Println("Hello from Plugin 02!")
	p.Execute()
}

func (p Plugin02) Stop() {
	log.Println("Stop from Plugin 02!")
}

func (p Plugin02) Execute() {
	// Create a new resty client
	client := resty.New()

	// Make GET request to server
	resp, err := client.R().Get(noahlib.GetHostHeartbeatURL())
	if err != nil {
		log.Printf("Error making request: %v\n", err)

		return
	}

	// Log the response
	log.Printf("Response Status: %v\n", resp.Status())
	log.Printf("Response Body: %s\n", resp.String())
}
