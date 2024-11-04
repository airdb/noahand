package coremain

import (
	"log"
	"plugin"
)

func RunPlugin() {
	pluginMap := map[string]string{
		"output/01_plugin.so": "Plugin01",
		"output/02_plugin.so": "Plugin02",
	}

	for pluginPath, pluginName := range pluginMap {
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
