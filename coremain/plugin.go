package coremain

import (
	"fmt"
	"plugin"
)

func RunPlugin() {
	pluginMap := map[string]string{
		"output/01_plugin.so": "Plugin01",
		"output/02_plugin.so": "Plugin02",
	}

	for pluginPath, pluginName := range pluginMap {
		p, err := plugin.Open(pluginPath)
		if err != nil {
			fmt.Println("Error loading plugin:", err)
			continue
		}

		// 查找插件中的 Hello 函数
		helloSymbol, err := p.Lookup(pluginName)
		if err != nil {
			fmt.Println("Error finding Hello function:", err)
			return
		}

		// 类型断言并调用 Hello 函数
		helloFunc, ok := helloSymbol.(func())
		if !ok {
			fmt.Println("Error asserting Hello function type")
			return
		}

		helloFunc() // 调用插件函数
	}
}
