# Noah 系统架构

本文档描述 Noah（noahand）的整体架构：进程模型、心跳机制、本地管理接口、路径与分发约定，以及部署方式。

## 1. 架构总览

Noah 采用 Master/Worker 架构：

- **Eagle（Master 进程）**：负责启动、监控和重启 Zeus，处理系统信号，运行心跳上报。
- **Zeus（Worker 进程）**：负责加载插件（Plugin）并管理模块（Module），执行具体任务。
- **Plugin（Zeus 插件）**：Zeus 的内置功能扩展，以 Go plugin（`.so`）形式动态加载，位于 `/opt/noah/plugins/`。
- **Module（模块）**：独立的第三方 agent 进程，由 Zeus 插件启动和管理，位于 `/opt/noah/modules/`。

```
                        ┌─────────────────────────┐
   heartbeat (POST)     │   Eagle (master 进程)    │
 ──────────────────────►│  - fork/监控/重启 Zeus    │
   中心服务              │  - 信号处理与转发         │
                        │  - 随机间隔心跳上报        │
                        └───────────┬─────────────┘
                                    │ fork (env: master=0)
                        ┌───────────▼─────────────┐
                        │   Zeus (worker 进程)     │
                        │  - 加载 .so 插件          │
                        │  - 管理 module           │
                        └───────────┬─────────────┘
              ┌─────────────┬───────┴──────┬──────────────┐
              ▼             ▼              ▼              ▼
         plugin 01     plugin 02      plugin N       modules/
         (.so)         (.so)          (.so)          (独立 agent)
```

另外，主进程内还运行一个本地管理 HTTP 服务（见第 4 节）。

## 2. 进程模型（Eagle / Zeus）

实现位于 `internal/noahlib/supervisor.go`：

- Eagle 通过 `os.StartProcess` 以相同二进制、追加环境变量 `master=0` 的方式 fork 出 Zeus；子进程检测到 `master=0` 后以 worker 身份运行。
- 进程名通过改写 argv 设置，便于在 `ps` 中区分：
  - master：`noah:master process - eagle`
  - worker：`noah:worker process - zeus`
- **崩溃重启**：Eagle `Wait()` 等待 Zeus 退出，若非正常退出则自动重新拉起。
- **信号处理**：
  - `SIGTERM` / `SIGINT`：Eagle 将信号转发给 Zeus 后退出。
  - `SIGHUP`：Eagle 拉起新的 Zeus，并向旧 Zeus 发送 `SIGHUP`，实现平滑重启。
- Worker 侧通过 `RegisterExitSignals` / `RegisterGraceSignals` 注册退出与优雅关闭回调，在 `WaitWorkerSignals` 中统一执行。

主流程入口为 `internal/coremain/run.go` 的 `Run()`：校验 root 权限 → 初始化配置 → 启动 supervisor → 启动管理 HTTP 服务（goroutine）→ 每分钟执行一次插件管理循环 → 阻塞等待信号。

## 3. 心跳机制

实现位于 `internal/noahlib/heartbeat.go`：

- Eagle 启动后开启心跳 goroutine，每轮在 `[0, MaxSleepInterval)` 秒内随机休眠后上报一次，避免大规模 agent 同时上报造成洪峰。
- `MaxSleepInterval` 默认 300 秒（`pkg/configkit/vars.go`）；agent 数量超过 2 万时建议调整为 900 秒。
- 上报内容（JSON POST 到 `{domain}/apis/infra/v1/host/heartbeat`）：

| 字段 | 说明 |
| --- | --- |
| `ip` | 本机非回环 IPv4 地址 |
| `hostname` | 主机名 |
| `arch` | CPU 架构（GOARCH） |
| `os` | 操作系统（GOOS） |
| `kernel` | 内核信息 |
| `system_info` | 系统信息 |
| `environment` | 环境标识 |

- 本地 sqlite（`internal/sqlite`，gorm）中定义了心跳记录表，用于本地留存。

## 4. 本地管理接口（Admin API）

实现位于 `internal/admin/`，基于 chi 路由，监听 `configkit.AdminAddr`（默认 `127.0.0.1:403`，仅本机可访问）：

| 路由 | 说明 |
| --- | --- |
| `GET /` | 默认页 |
| `GET /ping`、`GET /health` | 健康检查 |
| `GET /internal` | 列出可用的内部 API |
| `GET /internal/noah/metrics` | Prometheus 指标 |
| `GET /internal/noah/config` | 运行时配置 |
| `GET /internal/noah/host` | 主机信息 |
| `GET /internal/noah/selfupdate`、`/selfupgrade` | 自更新 |
| `GET /internal/noah/download_plugin` | 下载插件 |
| `GET /internal/noah/cmd`、`/exec` | 执行命令 |
| `GET /internal/noah/reset_passwd` | 重置密码 |
| `GET /html/*` | 静态文件（`./html` 目录） |

> 注意：管理接口目前为纯 HTTP，无认证机制，依赖仅绑定 `127.0.0.1` 提供隔离。其中命令执行、重置密码等端点属于高危操作，请勿将该端口暴露到本机之外。

## 5. 目录与路径约定

| 路径 | 说明 |
| --- | --- |
| `/opt/noah/` | Noah 家目录 |
| `/opt/noah/noah` | 主进程二进制 |
| `/opt/noah/plugins/` | Zeus 插件目录（`.so`） |
| `/opt/noah/modules/` | Module 目录 |
| `/tmp/noah/` | 临时文件目录 |
| `/opt/noah/noah.db` | sqlite 数据库文件 |

开发模式（`go run` 启动，见 `configkit.IsRunGoBuild`）下路径切换为：家目录 `./`、插件 `./output/plugins`、模块 `./output/modules`，域名切换为 `http://127.0.0.1:8080`。

## 6. 下载与分发规范

- Noah 主文件：`/download/noah/noah.tar.gz`
- Plugin 包：`/download/noah/plugin.tar.gz`
- Module 包：`/download/noah/module_name.tar.gz`

每个 tarball 附带同名 `.md5` 校验文件（如 `noah.tar.gz.md5`），下载后校验完整性：

```bash
md5sum -c noah.tar.gz.md5
```

`pkg/oskit` 实现了完整的「下载 → MD5 校验 → 解压 → 安装」流程（支持断点续传，支持 tar / tar.gz / zip）。

## 7. 服务管理与部署

- **Linux**：使用 systemd 管理。`noah init` 会生成 unit 文件到 `/etc/systemd/system/noah.service`（模板见 `pkg/service/systemd.go`）：

  ```bash
  noah init
  sudo systemctl daemon-reload
  sudo systemctl start noah
  sudo systemctl enable noah
  sudo journalctl -u noah -f
  ```

- **macOS**：使用 launchctl 管理，plist 路径 `/Library/LaunchDaemons/com.example.noah.plist`（模板见 `pkg/service/launchctl.go`）。

部署流程：下载主文件 / 插件包 / 模块包 → MD5 校验 → 分别解压到 `/opt/noah/`、`/opt/noah/plugins/`、`/opt/noah/modules/` → 通过 systemd / launchctl 启动。

版本管理：使用语义化版本，构建时通过 ldflags 注入版本信息（`pkg/version`，含 Repo / Version / Build / BuildTime / GoVersion / Arch），`noah --version` 可查看。

## 8. 代码结构导览

| 目录 | 说明 |
| --- | --- |
| `cmd/` | cobra 命令行入口：`run` / `init` / `plugin` / `update` |
| `internal/coremain/` | 核心编排：启动 supervisor、管理服务、插件加载器（`PluginInterface`） |
| `internal/admin/` | chi 本地管理 HTTP 服务（路由 + handler） |
| `internal/noahlib/` | Eagle/Zeus 监督进程、心跳上报、主机信息采集 |
| `internal/sqlite/` | gorm + sqlite 持久层（心跳表） |
| `pkg/configkit/` | 配置加载（viper）与全局路径 / 端口 / 域名常量、插件发现 |
| `pkg/monitorkit/` | Prometheus 指标定义 |
| `pkg/oskit/` | 下载、解压、安装、进程等 OS 工具函数 |
| `pkg/service/` | systemd / launchd 服务文件模板生成 |
| `pkg/version/` | 构建版本信息（ldflags 注入） |
| `plugins/` | 示例插件（01–04），演示 `PluginInterface` 的实现 |
