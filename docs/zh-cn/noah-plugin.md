# Noah Plugin 开发规范

Zeus Plugin 是 Noah 系统的内置功能扩展，部署于 `/opt/noah/plugins/`，支持热插拔。插件以 Go plugin（`.so` 文件）形式存在，通过 `go build -buildmode=plugin` 编译生成，由 Zeus 在运行时动态加载。

## 1. 接口规范

插件必须实现 `PluginInterface` 接口（定义于 `internal/coremain/plugin.go`）：

```go
type PluginInterface interface {
	GetInfo()

	Init()  // 插件加载时调用
	Start() // 启动插件核心逻辑
	Stop()  // 停止插件并清理资源
}
```

同时，插件必须导出一个名为 `Plugin` 的变量，Noah 通过 `plugin.Open` 加载 `.so` 文件后，用 `Lookup("Plugin")` 获取该符号并断言为 `PluginInterface`。

> 注意：向接口新增方法前，必须确保生产环境所有存量插件都已实现新方法，否则类型断言会失败。

## 2. 最小示例

参考 `plugins/01_plugin/main.go`：

```go
package main

import "log"

// Plugin01 实现 PluginInterface 接口
type Plugin01 struct{}

// 导出变量必须命名为 Plugin，且使用具体类型
// 错误写法：var Plugin shared.PluginInterface = Plugin01{}
var Plugin = Plugin01{}

func (p Plugin01) GetInfo() {
	log.Println("Plugin01 version: v0.1.0")
}

func (p Plugin01) Init() {
	log.Println("Init from Plugin 01!")
}

func (p Plugin01) Start() {
	log.Println("Hello from Plugin 01!")
}

func (p Plugin01) Stop() {
	log.Println("Stop from Plugin 01!")
}
```

## 3. 目录结构

```
plugins/
└── plugin_name/
    ├── main.go
    ├── config/
    ├── lib/
    └── README.md
```

## 4. 构建

Go plugin 依赖 CGO：

```bash
CGO_ENABLED=1 go build -buildmode=plugin -o output/plugins/01_plugin.so plugins/01_plugin/main.go
```

仓库根目录提供了 task 任务一键构建全部示例插件：

```bash
task build-plugin
```

> 限制：Go plugin 仅支持 Linux / macOS，且插件与主程序必须使用相同的 Go 版本和依赖版本编译。

## 5. 发现与加载机制

- Noah 扫描插件目录（生产为 `/opt/noah/plugins/`，开发模式为 `./output/plugins`）下所有 `.so` 文件，登记到插件表中（`pkg/configkit` 的 `CheckPlugins` / `UpdatePluginMap`）。
- 主流程每分钟刷新一次插件表并执行插件（`internal/coremain` 的 `RunPlugin`）。
- 也可通过 `noah plugin` 命令单独以插件模式运行一次。

## 6. 限制规范

- 包大小：单个 plugin 包建议控制在 **50MB** 以内，超过则不允许上传。
