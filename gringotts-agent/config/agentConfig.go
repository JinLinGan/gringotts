// Package config 配置
package config

import (
	"os"
	"runtime"
	"sync"

	"github.com/jinlingan/gringotts/gringotts-agent/log"
	"github.com/pkg/errors"
)

const (
	// PermissionMode 代表文件的默认权限
	PermissionMode       = 0750
	defaultServerAddress = "gringotts-server:7777"
	defaultWinWorkPath   = `c:\gringotts-agent`
	defaultLinuxWorkPath = "/var/gringotts"
)

// AgentConfig Agent 的配置
type AgentConfig struct {
	sync.RWMutex
	workPath            string
	serverAddress       string
	executerDirName     string
	downloadTempDirName string
	logger              log.Logger
}

//GetDefaultWorkPath 获取默认的工作目录
func GetDefaultWorkPath() string {
	if runtime.GOOS == "windows" {
		return defaultWinWorkPath
	}
	return defaultLinuxWorkPath
}

//GetDefaultServerAddress 获取默认的服务器地址
func GetDefaultServerAddress() string {
	return defaultServerAddress
}

// NewConfig 新建配置
func NewConfig(workPath string) (*AgentConfig, error) {
	c := &AgentConfig{
		workPath:            GetDefaultWorkPath(),
		serverAddress:       GetDefaultServerAddress(),
		executerDirName:     "executer",
		downloadTempDirName: "tmp",
		logger:              log.NewStdoutLogger(),
	}

	if workPath != GetDefaultServerAddress() {
		if err := c.setWorkPath(workPath); err != nil {
			return nil, errors.Wrapf(err, "can not set work path to %s", workPath)
		}
	}
	return c, nil
}

// setWorkPath 设置工作目录
func (c *AgentConfig) setWorkPath(path string) error {

	// 判断 path 是否存在
	if _, err := os.Stat(path); os.IsNotExist(err) {
		c.logger.Infof("dir %s not exist, to create it", path)
		// 新建目录
		if err := os.MkdirAll(path, PermissionMode); err != nil {
			return errors.Wrapf(err, "can not make dir %s", path)
		}
	}

	c.Lock()
	c.workPath = path
	c.Unlock()

	//create executerDir
	if _, err := os.Stat(c.GetExecuterPath()); os.IsNotExist(err) {
		c.logger.Infof("dir %s not exist, to create it", c.GetExecuterPath())
		// 新建目录
		if err := os.MkdirAll(c.GetExecuterPath(), PermissionMode); err != nil {
			return errors.Wrapf(err, "can not make dir %s", c.GetExecuterPath())
		}
	}
	//create downloadTempDir
	if _, err := os.Stat(c.GetDownloadTempPath()); os.IsNotExist(err) {
		c.logger.Infof("dir %s not exist, to create it", c.GetDownloadTempPath())
		// 新建目录
		if err := os.MkdirAll(c.GetDownloadTempPath(), PermissionMode); err != nil {
			return errors.Wrapf(err, "can not make dir %s", c.GetDownloadTempPath())
		}
	}

	return nil
}

// GetAgentRunningInfoFilePath 获取 Agent ID 和 版本记录文件所在目录
func (c *AgentConfig) GetAgentRunningInfoFilePath() string {
	c.RLock()
	defer c.RUnlock()
	runningInfoPath := c.GetWorkPath() + string(os.PathSeparator) + "runinfo"
	return runningInfoPath + string(os.PathSeparator) + "agent.info"
}

// GetWorkPath 获取工作目录名称
func (c *AgentConfig) GetWorkPath() string {
	c.RLock()
	defer c.RUnlock()
	return c.workPath
}

// GetExecuterPath 获取执行器目录
func (c *AgentConfig) GetExecuterPath() string {
	c.RLock()
	defer c.RUnlock()
	return c.workPath + string(os.PathSeparator) + c.executerDirName

}

// GetDownloadTempPath 获取用于存放下载临时文件的路径
func (c *AgentConfig) GetDownloadTempPath() string {
	c.RLock()
	defer c.RUnlock()
	return c.workPath + string(os.PathSeparator) + c.downloadTempDirName
}

// GetServerAddress 获取服务器路径
func (c *AgentConfig) GetServerAddress() string {
	c.RLock()
	defer c.RUnlock()
	return c.serverAddress
}

// SetServerAddress 设置服务器路径
func (c *AgentConfig) SetServerAddress(s string) {
	c.Lock()
	defer c.Unlock()
	c.serverAddress = s
}

// SetLogger 设置 logger 用于替换原有的标准输出 logger
func (c *AgentConfig) SetLogger(new log.Logger) {
	c.Lock()
	defer c.Unlock()
	c.logger = new
}
