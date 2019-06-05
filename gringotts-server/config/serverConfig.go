package config

import (
	"os"
	"runtime"
	"sync"
)

// ServerConfig 服务端配置
type ServerConfig struct {
	sync.RWMutex
	listenerPort    string
	externalAddress string
	workDir         string
}

const (
	defaultListenerPort    = ":7777"
	defaultExternalAddress = "gringotts-server"
	defaultWinWorkPath     = `c:\gringotts\gringotts-server`
	defaultLinuxWorkPath   = "/var/gringotts/gringotts-server"
)

//GetDefaultWorkPath 获取默认的工作目录
func GetDefaultWorkPath() string {
	if runtime.GOOS == "windows" {
		return defaultWinWorkPath
	}
	return defaultLinuxWorkPath
}

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
		workDir:         GetDefaultWorkPath(),
	}
}

// GetListenerPort 获取监听地址
func (s *ServerConfig) GetListenerPort() string {
	s.RLock()
	defer s.RUnlock()
	return s.listenerPort
}

// SetListenerPort 设置监听地址
func (s *ServerConfig) SetListenerPort(port string) {
	s.Lock()
	defer s.Unlock()
	s.listenerPort = port
}

// GetExternalAddress 获取监听地址
func (s *ServerConfig) GetExternalAddress() string {
	s.RLock()
	defer s.RUnlock()
	return s.externalAddress
}

// SetExternalAddress 设置外部监听地址
func (s *ServerConfig) SetExternalAddress(add string) {
	s.Lock()
	defer s.Unlock()
	s.externalAddress = add

}

// GetLogFilePath 获取日志文件存储路径
func (s *ServerConfig) GetLogFilePath() string {
	return s.workDir + string(os.PathSeparator) + "gringotts-server.log"
}

// SetWorkPath 设置工作目录
func (s *ServerConfig) SetWorkPath(path string) {
	s.Lock()
	defer s.Unlock()
	s.workDir = path
}
