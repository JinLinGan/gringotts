package config

import (
	"os"
)

// ServerConfig 服务端配置
type ServerConfig struct {
	//TODO:加锁
	listenerPort    string
	externalAddress string
}

const (
	defaultListenerPort    = ":7777"
	defaultExternalAddress = "gringotts-server"
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
		name = defaultExternalAddress
	}
	//主机名+端口
	return name + GetDefaultListenerPort()
}

// NewServerConfig 新建服务配置
func NewServerConfig() *ServerConfig {
	return &ServerConfig{
		listenerPort:    GetDefaultListenerPort(),
		externalAddress: GetDefaultExternalAddress(),
	}
}

// GetListenerPort 获取监听地址
func (s *ServerConfig) GetListenerPort() string {
	return s.listenerPort
}

// SetListenerPort 设置监听地址
func (s *ServerConfig) SetListenerPort(port string) {
	s.listenerPort = port

}

// GetExternalAddress 获取监听地址
func (s *ServerConfig) GetExternalAddress() string {
	return s.externalAddress
}

// SetExternalAddress 设置外部监听地址
func (s *ServerConfig) SetExternalAddress(add string) {
	s.externalAddress = add
}
