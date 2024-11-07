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
		log.Println("Error loading plugin:", err)

		return
	}

	// 查找插件中的 Hello 函数
	helloSymbol, err := plugin.Lookup(pluginName)
	if err != nil {
		log.Println("Error finding Hello function:", err)

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
