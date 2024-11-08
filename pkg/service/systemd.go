package service

import (
	"fmt"
	"guardhouse/pkg/configkit"
	"os"
	"path/filepath"
	"text/template"
)

const SystemdServiceTemplate = `
# sudo cp noah.service /etc/systemd/system/noah.service
# sudo systemctl daemon-reload 
# sudo systemctl enable noah
# sudo journalctl -u noah.service
# sudo journalctl -u noah.service -f

[Unit]
Description={{.Description}}
Documentation={{.Documentation}}
After=
Wants=

[Service]
User={{.User}}
Group={{.Group}}
Environment={{.Environment}}
WorkingDirectory={{.WorkingDirectory}}
ExecStart={{.ExecStart}}
ExecReload=/bin/kill -HUP $MAINPID
StartLimitInterval=0
LimitNOFILE=65536
Restart=always
#Restart=on-failure
RestartSec=20
KillSignal=SIGTERM
TimeoutStopSec=20
SendSIGKILL=yes
TimeoutStartSec=30
TimeoutStopSec=30

[Install]
WantedBy=multi-user.target
`

// ServiceConfig 用于存储服务文件配置
type ServiceConfig struct {
	Description      string
	Documentation    string
	ExecStart        string
	User             string
	Group            string
	WorkingDirectory string
	Environment      string
}

// GenerateServiceFile 生成 systemd 服务文件
func GenerateServiceFile(serviceConfig ServiceConfig, outputPath string) error {
	// 设置默认值
	if serviceConfig.Description == "" {
		serviceConfig.Description = "Noah Server Management System"
	}
	if serviceConfig.Documentation == "" {
		serviceConfig.Documentation = "https://aid.run"
	}
	if serviceConfig.User == "" {
		serviceConfig.User = "root"
	}
	if serviceConfig.Group == "" {
		serviceConfig.Group = "root"
	}
	if serviceConfig.Environment == "" {
		serviceConfig.Environment = "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin"
	}

	// 创建模板并解析内容
	tmpl, err := template.New(configkit.SystemdFilename).Parse(SystemdServiceTemplate)
	if err != nil {
		return fmt.Errorf("error parsing template: %v", err)
	}

	// 确保目录存在
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("error creating directory: %v", err)
	}

	// 创建服务文件
	file, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return fmt.Errorf("error creating service file: %v", err)
	}
	defer file.Close()

	// 写入文件
	if err := tmpl.Execute(file, serviceConfig); err != nil {
		return fmt.Errorf("error writing service file: %v", err)
	}

	return nil
}
