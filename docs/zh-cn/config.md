# Noah 配置说明

配置是 Noah 最重要的组成部分之一。出于安全性和变更频率的考虑，Noah 的配置文件**不支持运行时动态修改**。

为了支持回滚和可追溯性，设计上配置的每次变更都会记录到本地 sqlite 数据库（`/opt/noah/noah.db`）中，便于查看历史配置。

> 注：配置历史记录目前为规划中的能力，当前代码（`internal/sqlite`）仅实现了心跳记录表。

## 1. 配置分层

Noah 的配置分为三层：

| 层级 | 说明 |
| --- | --- |
| Eagle（noah 主配置） | 整体配置文件，控制运行模式、运行时环境等 |
| Zeus Plugin | 插件配置（每个插件一份） |
| Module | 模块配置（每个模块一份） |

## 2. 配置加载

实现位于 `pkg/configkit/`：

- **文件加载**（`config.go` 的 `Init()`）：使用 viper 读取 YAML 配置文件，依次在当前目录 `.` 和家目录 `/opt/noah/` 中查找，找不到则 panic。
- **运行时初始化**（`vars.go` 的 `InitConfig()`）：根据运行方式设置全局变量。以 `go run` 启动（开发模式）时，家目录切换为 `./`，插件目录为 `./output/plugins`，远端域名切换为 `http://127.0.0.1:8080`。

## 3. 主配置（Config）

对应 `pkg/configkit/config.go` 中的 `Config` 结构：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `run_mode` | string | 运行模式，默认 `local` |
| `runtime` | string | 运行时环境标识（当前存放远端域名） |
| `build_info` | object | 构建信息（版本、构建时间等，由 ldflags 注入） |
| `admin_api_list` | []string | 管理接口暴露的 API 列表 |

示例 YAML：

```yaml
run_mode: local
runtime: "https://aid.run"
admin_api_list:
  - /internal/noah/host
  - /internal/noah/selfupdate
  - /internal/noah/download_plugin
  - /internal/noah/exec
```

常用的全局默认值（`pkg/configkit/vars.go`）：

| 变量 | 默认值 | 说明 |
| --- | --- | --- |
| `HomeDir` | `/opt/noah/` | 家目录 |
| `AdminAddr` | `127.0.0.1:403` | 管理接口监听地址 |
| `MaxSleepInterval` | `300`（秒） | 心跳随机间隔上限；agent 超过 2 万台建议设为 900 |
| `DefaultDomain` | `https://aid.run` | 中心服务域名 |
| `TmpDir` | `/tmp/noah/` | 临时目录 |

## 4. Plugin 配置（PluginConfig）

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `plugin_name` | string | 插件名 |
| `plugin_path` | string | 插件路径（`.so` 文件） |
| `plugin_version` | string | 插件版本 |
| `plugin_priority` | int | 加载优先级 |
| `plugin_status` | string | 插件状态，见下方枚举 |

插件状态枚举：`enabled` / `disabled` / `unknown` / `failed` / `outdated`

## 5. Module 配置（ModuleConfig）

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `module_name` | string | 模块名 |
| `module_path` | string | 模块路径 |
| `module_version` | string | 模块版本 |
| `module_priority` | int | 优先级 |
| `module_status` | string | 模块状态，见下方枚举 |

模块状态枚举：`enabled` / `disabled` / `unknown` / `failed` / `running` / `outdated`
