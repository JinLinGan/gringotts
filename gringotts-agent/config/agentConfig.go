// Package config 配置
package config

import (
	"os"
	"runtime"
	"sync"

	"github.com/jinlingan/gringotts/pkg/log"
	"github.com/pkg/errors"
)

const (
	// PermissionMode 代表文件的默认权限
	PermissionMode       = 0750
	defaultServerAddress = "server:6666"
	defaultWinWorkPath   = `c:\gringotts\gringotts-agent`
	defaultLinuxWorkPath = "/var/gringotts/gringotts-agent"
)

// AgentConfig Agent 的配置
type AgentConfig struct {
	sync.RWMutex
	workPath            string
	serverAddress       string
	executorDirName     string
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
		executorDirName:     "executor",
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

	err := checkAndMkdir(path)
	if err != nil {
		return err
	}

	c.Lock()
	c.workPath = path
	c.Unlock()

	//create executorDir
	execPath := c.GetExecutorPath()
	err = checkAndMkdir(execPath)
	if err != nil {
		return err
	}

	//create downloadTempDir
	downloadTmp := c.GetDownloadTempPath()
	err = checkAndMkdir(downloadTmp)
	if err != nil {
		return err
	}

	//create runningInfo
	runningInfo := c.GetWorkPath() + string(os.PathSeparator) + "runinfo"
	err = checkAndMkdir(runningInfo)
	if err != nil {
		return err
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

func checkAndMkdir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// 新建目录
		if err := os.MkdirAll(path, PermissionMode); err != nil {
			return errors.Wrapf(err, "can not make dir %s", path)
		}
	}
	return nil
}

// GetWorkPath 获取工作目录名称
func (c *AgentConfig) GetWorkPath() string {
	c.RLock()
	defer c.RUnlock()
	return c.workPath
}

// GetExecutorPath 获取执行器目录
func (c *AgentConfig) GetExecutorPath() string {
	c.RLock()
	defer c.RUnlock()
	return c.workPath + string(os.PathSeparator) + c.executorDirName

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
