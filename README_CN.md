# noahand

[English](README.md) | [中文版](README_CN.md)

**noahand**（Noah）是一个轻量级的 Linux 主机 agent，用于云服务器的集中管理。它以守护进程方式运行，采用类似 Nginx 的 Master/Worker 进程模型 —— Master 进程（**Eagle**）负责监督和拉起 Worker 进程（**Zeus**），Zeus 负责执行任务和加载插件。

## 功能特性

- **Master/Worker 进程监督** —— Eagle fork 出 Zeus，崩溃时自动重启，并转发 `SIGTERM` / `SIGINT` / `SIGHUP` 信号实现优雅退出和重载
- **心跳上报** —— 以随机间隔向中心服务上报主机信息（IP、主机名、架构、操作系统、内核版本），避免上报洪峰
- **本地管理接口** —— 基于 chi 的 HTTP 接口，监听 `127.0.0.1:403`，提供健康检查、运行时配置、主机信息、自更新等能力
- **Prometheus 指标** —— 通过 `/internal/noah/metrics` 暴露
- **插件体系** —— 通过 Go plugin（`.so`）扩展，运行时动态加载
- **服务化集成** —— `noah init` 生成 systemd unit 文件（Linux）；同时支持 macOS launchd

## 环境要求

- Go 1.26+
- Linux（主要目标平台；macOS 可用于开发）
- 运行 agent 需要 root 权限

## 快速开始

```bash
# 克隆仓库
git clone https://github.com/airdb/noahand.git
cd noahand

# 构建插件和 noah 二进制（需要 go-task）
task build-plugin build

# 运行 agent（需要 root）
sudo ./output/noah run
```

## 命令行

| 命令 | 说明 |
| --- | --- |
| `noah run` | 启动 agent（Eagle/Zeus 监督进程、管理接口、插件循环） |
| `noah init` | 生成 systemd 服务文件 |
| `noah plugin` | 插件模式（加载并执行一次插件） |
| `noah update` | 自更新（占位，尚未实现） |

安装为系统服务：

```bash
noah init
sudo systemctl daemon-reload
sudo systemctl start noah
sudo systemctl enable noah
sudo journalctl -u noah -f   # 查看日志
```

## 文档

详细文档见 [docs/zh-cn/](docs/zh-cn/README.md)：

- [系统架构](docs/zh-cn/architecture.md) —— Eagle/Zeus 进程模型、心跳机制、管理接口、路径约定、部署
- [配置说明](docs/zh-cn/config.md) —— 配置文件设计与字段
- [插件开发](docs/zh-cn/noah-plugin.md) —— 如何编写 Zeus 插件（`.so`）
- [模块开发](docs/zh-cn/noah-module.md) —— 模块打包与安装规范

## 参与贡献

欢迎提交 Issue 和 Pull Request：[GitHub Issues](https://github.com/airdb/noahand/issues)

## 许可证

[Apache License 2.0](LICENSE)
