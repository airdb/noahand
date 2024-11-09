package main

import (
	"log"
)

// Define the Plugin01 structure and implement the PluginInterface interface
type Plugin01 struct{}

// Define the exported variable of the plugin, which must be of type PluginInterface
// but donot define like this.
// [wrong] var Plugin shared.PluginInterface = Plugin01{}
var Plugin = Plugin01{}

// Implement the GetInfo method
func (p Plugin01) GetInfo() {
	log.Println("Plugin01 version: v0.1.0")
}

func (p Plugin01) Init() {
	log.Println("Init from Plugin 01!")
}

// Implement the start method
func (p Plugin01) Start() {
	log.Println("Hello from Plugin 01!")
}

func (p Plugin01) Stop() {
	log.Println("Stop from Plugin 01!")
}
