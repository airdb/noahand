# Noah Plugin 开发规范

Zeus Plugin 是 Noah 系统的内置功能模块，位于 `/opt/noah/plugins/`，支持热插拔。本文档介绍了 Plugin 的开发规范和最佳实践。
plugin 是以 `.so` 文件形式存在的，可以通过 `go build -buildmode=plugin` 编译生成。

## Plugin 结构

插件（Zeus Plugin）规范
目录结构:

arduino
Copy code
plugins/
├── plugin_name/
│   ├── main.go
│   ├── config/
│   ├── lib/
│   └── README.md


接口规范: 插件需实现 Init、Run、Stop 等接口，以与 Zeus 无缝集成。


## 限制规范

包大小: 单个 module 包大小建议控制在 50MB 以内。 如果超过则不允许上传。