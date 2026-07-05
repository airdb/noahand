# Noah 文档

noahand 是一套主机 agent 管理解决方案，Noah 是具体的项目实现。

系统采用 Master/Worker 架构：Master 进程命名为 **Eagle**，Worker 进程命名为 **Zeus**。Eagle 负责启动和管理 Zeus；Zeus 负责加载插件（Plugin）并管理模块（Module）；Module 作为独立 agent 由 Zeus 插件启动，提供扩展功能。

核心能力：主机管理、插件扩展、任务分发、心跳监控与健康检查。

## 文档目录

| 文档 | 内容 |
| --- | --- |
| [architecture.md](architecture.md) | 系统架构：Eagle/Zeus 进程模型、心跳机制、本地管理接口、目录与分发约定、部署与服务管理、代码结构导览 |
| [config.md](config.md) | 配置设计：配置文件加载方式、字段说明、Plugin/Module 配置结构与状态 |
| [noah-plugin.md](noah-plugin.md) | Zeus 插件开发规范：接口定义、示例代码、构建与发现机制、限制规范 |
| [noah-module.md](noah-module.md) | Module 开发规范：目录结构、打包规范（deb 式）、下载安装流程、限制规范 |

## 快速链接

- 项目入口：[README（English）](../../README.md) / [README（中文）](../../README_CN.md)
- 构建与运行：见根目录 `Taskfile.yaml`（`task build-plugin`、`task build`、`task run`）
