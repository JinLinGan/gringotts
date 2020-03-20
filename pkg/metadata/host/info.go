package host

import (
	"strings"

	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/net"
)

// HostInfo 硬件及系统信息
type HostInfo struct {
	//CPU
	processor string
	// 主机名
	HostName string
	// 主机 UUID
	HostUUID string
	// 操作系统类型
	OS string
	// 操作系统家族
	PlatformFamily string
	//操作系统
	Platform string
	//操作系统版本
	PlatformVersion string
	// 内核版本
	KernelVersion string
	// 内核架构
	kernelArch string
	// 虚拟化系统
	VirtualizationSystem string
	// 虚拟化角色
	VirtualizationRole string

	//Mac地址列表（仅包含物理端口）
	Interfaces []InterfaceStat
}

type InterfaceStat struct {
	Name         string   `json:"Name"`         // e.g., "en0", "lo0", "eth0.100"
	HardwareAddr string   `json:"HardwareAddr"` // IEEE MAC-48, EUI-48 and EUI-64 form
	Addrs        []string `json:"IPAddrs"`
}

var InterfaceNamePrefixBlackList = [...]string{
	"cali", "veth", "flannel",
	"docker", "lo", "tunl",
	"ip6tnl0", "gif", "stf",
	"XHC", "llw", "utun",
	"bridge", "p2p", "awdl",
}

// GetHostInfo 获取主机信息
func GetHostInfo() *HostInfo {
	ret := &HostInfo{}

	// 获取主机信息
	gHInfo, err := host.Info()
	if err == nil {
		ret.HostName = gHInfo.Hostname
		ret.HostUUID = gHInfo.HostID
		ret.OS = gHInfo.OS
		ret.Platform = gHInfo.Platform
		ret.PlatformFamily = gHInfo.PlatformFamily
		ret.PlatformVersion = gHInfo.PlatformVersion
		ret.KernelVersion = gHInfo.KernelVersion
		ret.kernelArch = gHInfo.KernelArch
		ret.VirtualizationRole = gHInfo.VirtualizationRole
		ret.VirtualizationSystem = gHInfo.VirtualizationSystem
	}

	ret.Interfaces = getNetInfo()

	return ret
}

func getNetInfo() []InterfaceStat {
	var netInfo []InterfaceStat
	// 获取网卡信息
	gNInfo, err := net.Interfaces()

	if err == nil {
		for _, i := range gNInfo {
			if checkNetName(i.Name) {
				n := InterfaceStat{
					Name:         i.Name,
					HardwareAddr: i.HardwareAddr,
				}

				for _, a := range i.Addrs {
					n.Addrs = append(n.Addrs, a.Addr)
				}

				netInfo = append(netInfo, n)
			}
		}

	}
	return netInfo
}

func checkNetName(n string) bool {
	b := true

	for _, p := range InterfaceNamePrefixBlackList {
		if strings.HasPrefix(n, p) {
			b = false
			break
		}
	}
	return b
}
