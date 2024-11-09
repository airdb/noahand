package coremain

import (
	"fmt"
	"guardhouse/pkg/configkit"
	"log"
	"plugin"
)

// Plugin must implement the PluginInterface interface.
// Before add new method, please add it ensure all plugins has all the method in prod env.
type PluginInterface interface {
	GetInfo()

	Init()  // Init is called when the plugin is loaded.
	Start() // Start the core logic of the plugin.
	Stop()  // stop the plugin, and clean up resources.
}

func RunPlugin() {

	log.Println("Running plugins...", configkit.PluginMap)

	// Keep the plugin map up to date.
	configkit.UpdatePluginMap()

	for pluginPath, pluginName := range configkit.PluginMap {
		// RunPluginWithParams(pluginPath, pluginName)
		log.Println(pluginPath, pluginName)
		RunPluginWithInterface(pluginPath)
	}
}

// RunPluginWithParams runs a plugin with the given path and name.
func RunPluginWithInterface(pluginPath string) {
	p, err := plugin.Open(pluginPath)
	if err != nil {
		log.Printf("Error loading plugin from path %s: %v", pluginPath, err)
		return
	}

	symPlugin, err := p.Lookup("Plugin")
	if err != nil {
		fmt.Println("Error finding Plugin symbol:", err)
		return
	}

	pluginInstance, ok := symPlugin.(PluginInterface)
	if !ok {
		fmt.Printf("Plugin at %s does not implement PluginInterface\n", pluginPath)
		return
	}

	pluginInstance.GetInfo()
	pluginInstance.Start()
	// fmt.Println("Loaded plugin: +v\n", info)
}

func RunPluginWithParams(pluginPath string, pluginName string) {
	plugin, err := plugin.Open(pluginPath)
	if err != nil {
		log.Printf("Error loading plugin from path %s: %v", pluginPath, err)
		return
	}

	// 查找插件中的 Hello 函数
	helloSymbol, err := plugin.Lookup(pluginName)
	if err != nil {
		log.Printf("Error finding function %s in plugin %s: %v", pluginName, pluginPath, err)
		return
	}

	log.Println(pluginPath, pluginName)

	// 类型断言并调用 Hello 函数
	helloFunc, ok := helloSymbol.(func())
	if !ok {
		log.Println("Error asserting Hello function type", ok)
		return
	}

	helloFunc() // 调用插件函数
}
