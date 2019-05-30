package config

import (
	"os"
)

// ServerConfig 服务端配置
type ServerConfig struct {
}

const (
	defaultListenerPort     = ":7777"
	defaultExternalHostname = "gringotts-server"
)

//GetDefaultListenerPort 获取默认监听端口
func GetDefaultListenerPort() string {
	return defaultListenerPort
}

//GetDefaultExternalAddress 获取默认的连接地址
func GetDefaultExternalAddress() string {
	// 获取主机名
	name, err := os.Hostname()
	if err != nil {
		name = defaultExternalHostname
	}
	//主机名+端口
	return name + GetDefaultListenerPort()
}
