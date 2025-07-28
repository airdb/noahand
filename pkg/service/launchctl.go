package service

import (
	"fmt"
	"noahand/pkg/configkit"
	"os"
	"path/filepath"
	"text/template"
)

const LaunchctlServiceTemplate = `
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<!--
This is a multi-line comment.
sudo cp com.example.noah.plist /Library/LaunchDaemons/com.example.noah.plist
sudo launchctl load /Library/LaunchDaemons/com.example.noah.plist
sudo launchctl start com.example.noah
sudo launchctl stop com.example.noah
It can span multiple lines.
-->

<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.example.noah</string>

    <key>ProgramArguments</key>
    <array>
        <string>/opt/noah/noah</string>
        <string>run</string>
    </array>

    <key>WorkingDirectory</key>
    <string>/opt/noah</string>

    <key>RunAtLoad</key>
    <true/>

    <key>KeepAlive</key>
    <true/>

    <key>StandardOutPath</key>
    <string>/var/log/noah.log</string>
    <key>StandardErrorPath</key>
    <string>/var/log/noah.err</string>
</dict>
</plist>
`

// ServiceConfig 用于存储服务文件配置
type LaunchctlServiceConfig struct {
	Description      string
	Documentation    string
	ExecStart        string
	User             string
	Group            string
	WorkingDirectory string
	Environment      string
}

// GenerateServiceFile 生成 systemd 服务文件
func GenerateLaunchctlServiceFile(serviceConfig LaunchctlServiceConfig, outputPath string) error {
	// 设置默认值
	if serviceConfig.Description == "" {
		serviceConfig.Description = "Noah Server Management System"
	}
	if serviceConfig.Documentation == "" {
		serviceConfig.Documentation = configkit.DefaultDomain
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
	tmpl, err := template.New(configkit.LaunchctlFilename).Parse(LaunchctlServiceTemplate)
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
