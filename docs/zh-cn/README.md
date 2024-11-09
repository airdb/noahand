# Guardhouse

Guardhouse 是一套主机 agent 管理的解决方案，Noah 是具体的项目实现。

Noah (Agent Management) 系统开发与部署文档

## 概述
项目名称: Noah（Guardhouse）
系统架构: Master/Worker 模式，Master 命名为 Eagle，Worker 命名为 Zeus。
核心目标: 提供主机管理、插件扩展、任务分发、监控和健康检查等功能。

Noah 系统采用 Master/Worker 架构，Eagle 负责启动和管理 Zeus 进程，Zeus 作为 Worker 进程负责执行任务并管理插件及 module。module 作为独立 agent 通过 Zeus 插件启动，提供扩展功能。


## 系统目录与路径约定
Noah 家目录: /opt/noah
主进程路径: /opt/noah/noah
Zeus 插件路径: /opt/noah/plugins/
module 路径: /opt/noah/modules/
临时文件路径: /tmp/noah/
数据库文件路径: /opt/noah/noah.db


## 下载路径规范
Noah 主文件: /download/noah/noah.tar.gz
Plugin 文件: /download/noah/plugin.tar.gz（通常与 Noah 主文件放置在一起）
Module 文件: /download/noah/module_name.tar.gz

MD5 校验文件
每个 tarball 文件会附带一个对应的 MD5 校验文件，文件名格式为：<文件名>.md5。
示例：
noah.tar.gz 的 MD5 校验文件为 noah.tar.gz.md5
plugin.tar.gz 的 MD5 校验文件为 plugin.tar.gz.md5
module_name.tar.gz 的 MD5 校验文件为 module_name.tar.gz.md5

在文件下载或传输后，可使用 md5sum 工具验证文件完整性，例如：

bash
Copy code
md5sum -c noah.tar.gz.md5


## 服务管理
Linux: 使用 systemd 管理 Noah 服务
macOS: 使用 launchctl 管理 Noah 服务


## 系统架构
Noah 系统包含以下组件：

### Eagle（Master 进程）

负责启动和监控 Zeus。
监听本地管理端口，用于任务下发和状态监控。
实现心跳机制，定期检查 Zeus 和 module 的运行状态。


### Zeus（Worker 进程）

负责加载插件并管理 module。
插件提供扩展功能，module 是独立 agent，通过 Zeus 插件启动。
插件（Zeus Plugins）

插件为 Zeus 的内置功能模块，位于 /opt/noah/plugins/，支持热插拔。
模块（Modules）

Module 是独立的 agent，位于 /opt/noah/modules/，通过 Zeus 插件启动和管理。

插件（Zeus Plugin）规范
开发规范参考 ./noah-plugin.md


模块（Modules）规范
开发规范参考 ./noah-module.md


## 部署流程

安装与配置
下载 Noah 主文件: /download/noah/noah.tar.gz

下载插件与模块:

插件包下载路径：/download/noah/plugin.tar.gz
各 module 包下载路径：/download/noah/module_name.tar.gz
MD5 校验:

通过 .md5 文件验证下载的 tarball 完整性：
bash
Copy code
md5sum -c /download/noah/noah.tar.gz.md5
解压文件并放置到指定目录:

Noah 主文件解压至 /opt/noah/。
插件解压至 /opt/noah/plugins/。
Module 解压至 /opt/noah/modules/。

## 启动服务
Linux 系统: 使用 systemd 启动 Noah。
macOS 系统: 使用 launchctl 启动 Noah。

## CI/CD 与版本管理
持续集成: 插件和 module 更新需通过 CI/CD 管道构建和测试。
版本控制: 版本需记录更新日志，使用语义化版本。
监控与调试
心跳监控: Eagle 负责监控 Zeus。
日志系统: 保存 Eagle、Zeus 和各 module 的日志。
调试工具: Noah 支持调试模式，便于问题排查。
