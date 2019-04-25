package config

import (
	"fmt"
	"log"
	"os"
)

var (
	workPath            string = "/usr/local/gringotts"
	serverAddress       string = "127.0.0.1:7777"
	executerDirName     string = "executer"
	downloadTempDirName string = "tmp"
)

func init() {
}

// SetWorkPath 设置工作目录
func SetWorkPath(path string) error {

	// 判断 path 是否存在
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Printf("dir %s not exist, to create it", path)
		// 新建目录
		if err := os.MkdirAll(path, 0775); err != nil {
			return fmt.Errorf("can not make dir %s: %s", path, err)
		}
	}
	workPath = path

	//create executerDir
	if _, err := os.Stat(GetExecuterPath()); os.IsNotExist(err) {
		log.Printf("dir %s not exist, to create it", GetExecuterPath())
		// 新建目录
		if err := os.MkdirAll(GetExecuterPath(), 0775); err != nil {
			return fmt.Errorf("can not make dir %s: %s", GetExecuterPath(), err)
		}
	}
	//create downloadTempDir
	if _, err := os.Stat(GetDownloadTempPath()); os.IsNotExist(err) {
		log.Printf("dir %s not exist, to create it", GetDownloadTempPath())
		// 新建目录
		if err := os.MkdirAll(GetDownloadTempPath(), 0775); err != nil {
			return fmt.Errorf("can not make dir %s: %s", GetDownloadTempPath(), err)
		}
	}

	return nil
}

// GetWorkPath 获取工作目录名称
func GetWorkPath() string {
	return workPath
}

func GetExecuterPath() string {
	return workPath + string(os.PathSeparator) + executerDirName

}
func GetDownloadTempPath() string {
	return workPath + string(os.PathSeparator) + downloadTempDirName
}

func GetServerAddress() string {
	return serverAddress
}

func SetServerAddress(s string) {
	serverAddress = s
}
