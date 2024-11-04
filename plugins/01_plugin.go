// pluginA.go
package main

import (
	"log"
)

// Main function of the plugin.
func Plugin01() {
	log.SetFlags(log.Lshortfile)
	log.Println("Hello from Plugin 01!")
}
