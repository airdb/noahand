// pluginA.go
package main

import (
	"fmt"
	"guardhouse/internal/noahlib"

	"github.com/go-resty/resty/v2"
)

func Plugin02() {
	// Create a new resty client
	client := resty.New()

	// Make GET request to sg.airdb.host:8000
	resp, err := client.R().
		Get(noahlib.DefaultDomain + "/info")

	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return
	}

	// Print the response
	fmt.Printf("Response Status: %v\n", resp.Status())
	fmt.Printf("Response Body: %s\n", resp.String())
}
