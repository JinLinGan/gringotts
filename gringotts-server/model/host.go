package model

import "github.com/jinzhu/gorm"

// Host 主机
type Host struct {
	gorm.Model
	HostName      string
	HostInterface []*HostInterface
}

// HostInterface 主机网卡
type HostInterface struct {
	gorm.Model
	HostID        uint
	HWAddr        string
	InterfaceName string
	IPAddress     string
}
