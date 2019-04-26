// Package config 配置
package config

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

const (
	dirPermission = 0750
)

// AgentConfig Agent 的配置
type AgentConfig struct {
	workPath            string
	serverAddress       string
	executerDirName     string
	downloadTempDirName string
}

// NewConfig 新建配置
func NewConfig() *AgentConfig {

	c := &AgentConfig{
		workPath:            "/usr/local/gringotts",
		serverAddress:       "127.0.0.1:7777",
		executerDirName:     "executer",
		downloadTempDirName: "tmp",
	}

	if runtime.GOOS == "windows" {
		c.workPath = `c:\gringotts-agent`
	}
	return c
}

// SetWorkPath 设置工作目录
func (c *AgentConfig) SetWorkPath(path string) error {

	// 判断 path 是否存在
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Printf("dir %s not exist, to create it", path)
		// 新建目录
		if err := os.MkdirAll(path, dirPermission); err != nil {
			return fmt.Errorf("can not make dir %s: %s", path, err)
		}
	}
	c.workPath = path

	//create executerDir
	if _, err := os.Stat(c.GetExecuterPath()); os.IsNotExist(err) {
		log.Printf("dir %s not exist, to create it", c.GetExecuterPath())
		// 新建目录
		if err := os.MkdirAll(c.GetExecuterPath(), dirPermission); err != nil {
			return fmt.Errorf("can not make dir %s: %s", c.GetExecuterPath(), err)
		}
	}
	//create downloadTempDir
	if _, err := os.Stat(c.GetDownloadTempPath()); os.IsNotExist(err) {
		log.Printf("dir %s not exist, to create it", c.GetDownloadTempPath())
		// 新建目录
		if err := os.MkdirAll(c.GetDownloadTempPath(), dirPermission); err != nil {
			return fmt.Errorf("can not make dir %s: %s", c.GetDownloadTempPath(), err)
		}
	}

	return nil
}

// GetWorkPath 获取工作目录名称
func (c *AgentConfig) GetWorkPath() string {
	return c.workPath
}

// GetExecuterPath 获取执行器目录
func (c *AgentConfig) GetExecuterPath() string {
	return c.workPath + string(os.PathSeparator) + c.executerDirName

}

// GetDownloadTempPath 获取用于存放下载临时文件的路径
func (c *AgentConfig) GetDownloadTempPath() string {
	return c.workPath + string(os.PathSeparator) + c.downloadTempDirName
}

// GetServerAddress 获取服务器路径
func (c *AgentConfig) GetServerAddress() string {
	return c.serverAddress
}

// SetServerAddress 设置服务器路径
func (c *AgentConfig) SetServerAddress(s string) {
	c.serverAddress = s
}
