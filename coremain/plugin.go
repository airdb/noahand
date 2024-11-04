package coremain

import (
	"guardhouse/pkg/configkit"
	"log"
	"plugin"
)

func RunPlugin() {
	for pluginPath, pluginName := range configkit.PluginMap {
		plugin, err := plugin.Open(pluginPath)
		if err != nil {
			log.Println("Error loading plugin:", err)

			continue
		}

		// 查找插件中的 Hello 函数
		helloSymbol, err := plugin.Lookup(pluginName)
		if err != nil {
			log.Println("Error finding Hello function:", err)

			return
		}

		// 类型断言并调用 Hello 函数
		helloFunc, ok := helloSymbol.(func())
		if !ok {
			log.Println("Error asserting Hello function type")

			return
		}

		helloFunc() // 调用插件函数
	}
}
