// pluginA.go
package main

import (
	"guardhouse/internal/noahlib"
	"log"

	"github.com/go-resty/resty/v2"
)

// Main function of the plugin.
func Plugin02() {
	// Create a new resty client
	client := resty.New()

	// Make GET request to sg.airdb.host:8000
	resp, err := client.R().
		Get(noahlib.DefaultDomain + "/info")
	if err != nil {
		log.Printf("Error making request: %v\n", err)

		return
	}

	// Log the response
	log.Printf("Response Status: %v\n", resp.Status())
	log.Printf("Response Body: %s\n", resp.String())
}
