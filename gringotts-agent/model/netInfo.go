package model

// NetInfos 网卡信息的数组
type NetInfos map[string]NetInfo

// NetInfo 网卡信息
type NetInfo struct {
	IPAddress  string
	MacAddress string
}
