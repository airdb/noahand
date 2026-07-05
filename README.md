# noahand

[English](README.md) | [中文版](README_CN.md)

**noahand** (Noah) is a lightweight Linux host agent for centralized cloud-server management. It runs as a daemon with an Nginx-style master/worker process model — the master process (**Eagle**) supervises and restarts the worker process (**Zeus**), which executes tasks and loads plugins.

## Features

- **Master/Worker supervision** — Eagle forks Zeus, restarts it on crash, and forwards `SIGTERM` / `SIGINT` / `SIGHUP` for graceful shutdown and reload
- **Heartbeat reporting** — periodically reports host info (IP, hostname, arch, OS, kernel) to a central server at randomized intervals to avoid thundering herds
- **Local admin API** — a chi-based HTTP interface on `127.0.0.1:403` for health checks, runtime config, host info, and self-update
- **Prometheus metrics** — exposed at `/internal/noah/metrics`
- **Plugin system** — extensible via Go plugins (`.so`), loaded dynamically at runtime
- **Service integration** — `noah init` generates a systemd unit file (Linux); launchd (macOS) is also supported

## Requirements

- Go 1.26+
- Linux (primary target; macOS supported for development)
- Root privileges to run the agent

## Quick Start

```bash
# Clone the repository
git clone https://github.com/airdb/noahand.git
cd noahand

# Build plugins and the noah binary (requires go-task)
task build-plugin build

# Run the agent (requires root)
sudo ./output/noah run
```

## CLI Commands

| Command | Description |
| --- | --- |
| `noah run` | Start the agent (Eagle/Zeus supervisor, admin API, plugin loop) |
| `noah init` | Generate a systemd service file |
| `noah plugin` | Run in plugin mode (load and execute plugins once) |
| `noah update` | Self-update (stub, not yet implemented) |

To install as a system service:

```bash
noah init
sudo systemctl daemon-reload
sudo systemctl start noah
sudo systemctl enable noah
sudo journalctl -u noah -f   # view logs
```

## Documentation

Detailed documentation (in Chinese) lives in [docs/zh-cn/](docs/zh-cn/README.md):

- [Architecture](docs/zh-cn/architecture.md) — Eagle/Zeus process model, heartbeat, admin API, path conventions, deployment
- [Configuration](docs/zh-cn/config.md) — config file design and fields
- [Plugin Development](docs/zh-cn/noah-plugin.md) — how to write a Zeus plugin (`.so`)
- [Module Development](docs/zh-cn/noah-module.md) — module packaging and installation spec

## Contributing

Issues and pull requests are welcome: [GitHub Issues](https://github.com/airdb/noahand/issues)

## License

[Apache License 2.0](LICENSE)
