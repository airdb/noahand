package coremain

import (
	"guardhouse/pkg/configkit"
	"log"
	"plugin"
)

func RunPlugin() {
	for pluginPath, pluginName := range configkit.PluginMap {
		RunPluginWithParams(pluginPath, pluginName)
	}
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
