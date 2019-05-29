// Package config 配置
package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

const (
	// PermissionMode 代表文件的默认权限
	PermissionMode = 0750
)

// AgentConfig Agent 的配置
type AgentConfig struct {
	sync.RWMutex
	workPath            string
	serverAddress       string
	executerDirName     string
	downloadTempDirName string
	isRegistered        bool
	agentInfo           *AgentRunningInfo
}

// AgentRunningInfo 保存了 agentID 和 配置版本
type AgentRunningInfo struct {
	agentID       string
	configVersion int64
}

//GetDefaultWorkPath 获取默认的工作目录
func GetDefaultWorkPath() string {
	if runtime.GOOS == "windows" {
		return `c:\gringotts-agent`
	}
	return "/usr/local/gringotts"
}

//GetDefaultServerAddress 获取默认的服务器地址
func GetDefaultServerAddress() string {
	return "127.0.0.1:7777"
}

// NewConfig 新建配置
func NewConfig(workPath string) (*AgentConfig, error) {
	c := &AgentConfig{
		workPath:            GetDefaultWorkPath(),
		serverAddress:       GetDefaultServerAddress(),
		executerDirName:     "executer",
		downloadTempDirName: "tmp",
	}

	if workPath != "" {
		if err := c.setWorkPath(workPath); err != nil {
			return nil, fmt.Errorf("can not set work path to %s: %s", workPath, err)
		}
	}

	agentInfo, err := c.getAgentIDFormWorkdir()
	if err != nil {
		log.Printf("read agent info faild so get state unregistered")
		c.isRegistered = false
		c.agentInfo = nil
	} else {
		c.isRegistered = true
		c.agentInfo = agentInfo
	}
	return c, nil
}

// setWorkPath 设置工作目录
func (c *AgentConfig) setWorkPath(path string) error {

	// 判断 path 是否存在
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Printf("dir %s not exist, to create it", path)
		// 新建目录
		if err := os.MkdirAll(path, PermissionMode); err != nil {
			return fmt.Errorf("can not make dir %s: %s", path, err)
		}
	}

	c.Lock()
	c.workPath = path
	c.Unlock()

	//create executerDir
	if _, err := os.Stat(c.GetExecuterPath()); os.IsNotExist(err) {
		log.Printf("dir %s not exist, to create it", c.GetExecuterPath())
		// 新建目录
		if err := os.MkdirAll(c.GetExecuterPath(), PermissionMode); err != nil {
			return fmt.Errorf("can not make dir %s: %s", c.GetExecuterPath(), err)
		}
	}
	//create downloadTempDir
	if _, err := os.Stat(c.GetDownloadTempPath()); os.IsNotExist(err) {
		log.Printf("dir %s not exist, to create it", c.GetDownloadTempPath())
		// 新建目录
		if err := os.MkdirAll(c.GetDownloadTempPath(), PermissionMode); err != nil {
			return fmt.Errorf("can not make dir %s: %s", c.GetDownloadTempPath(), err)
		}
	}

	return nil
}

func (c *AgentConfig) getAgentIDFormWorkdir() (*AgentRunningInfo, error) {
	c.RLock()
	defer c.RUnlock()
	path := c.getAgentRunningInfoFilePath()

	agentInfo := new(AgentRunningInfo)

	b, err := ioutil.ReadFile(filepath.Clean(path))
	if err != nil {
		// 如果文件不存在
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("agent info can not find: %v", err)
		}
		return nil, fmt.Errorf("read agent info faild: %s", err)
	}
	if err := json.Unmarshal(b, agentInfo); err != nil {
		return nil, fmt.Errorf("decode agent info file %s fail: %s", path, err)
	}
	return agentInfo, nil
}

func (c *AgentConfig) getAgentRunningInfoFilePath() string {
	c.RLock()
	defer c.RUnlock()
	runningInfoPaht := c.GetWorkPath() + string(os.PathSeparator) + "runinfo"
	return runningInfoPaht + string(os.PathSeparator) + "agent.info"
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

// GetAgentID 获取 AgentID
func (c *AgentConfig) GetAgentID() string {
	c.RLock()
	defer c.RUnlock()
	return c.agentInfo.agentID
}

// GetConfigVersion 获取配置版本
func (c *AgentConfig) GetConfigVersion() int64 {
	c.RLock()
	defer c.RUnlock()
	return c.agentInfo.configVersion
}

// SetConfigVersion 设置配置版本
func (c *AgentConfig) SetConfigVersion(v int64) error {

	c.Lock()
	c.agentInfo.configVersion = v
	c.Unlock()

	b, err := json.Marshal(c.agentInfo)
	if err != nil {
		return fmt.Errorf("encode agent config %+v fail: %s", c.agentInfo, err)
	}
	if err := ioutil.WriteFile(c.getAgentRunningInfoFilePath(), b, PermissionMode); err != nil {
		return fmt.Errorf("write agent config file %s fail: %s", c.getAgentRunningInfoFilePath(), err)
	}
	return nil
}

// IsRegistered Agent 是否成功注册
func (c *AgentConfig) IsRegistered() bool {
	c.RLock()
	defer c.RUnlock()
	return c.isRegistered
}
