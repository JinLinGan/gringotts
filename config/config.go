package config

import (
	"fmt"
	"os"
)

var (
	WorkDir       string = "/usr/local/gringotts"
	ServerAddress string = "127.0.0.1:7777"
)

func init() {
}

// SetWorkDir 设置工作目录
func SetWorkDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("path %s not exist", path)
	}
	WorkDir = path
	return nil
}

// GetWorkDir 获取工作目录名称
func GetWorkDir() string {
	return WorkDir
}
