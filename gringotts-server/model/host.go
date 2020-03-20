package model

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jinlingan/gringotts/pkg/util"

	"github.com/jinlingan/gringotts/pkg/log"

	"github.com/jinlingan/gringotts/pkg/message"
)

// Host 主机
type Host struct {
	AgentID              string       `protobuf:"bytes,1,opt,name=agentID,proto3" json:"agentID,omitempty"`
	HostName             string       `protobuf:"bytes,2,opt,name=hostName,proto3" json:"hostName,omitempty"`
	HostUUID             string       `protobuf:"bytes,3,opt,name=hostUUID,proto3" json:"hostUUID,omitempty"`
	Os                   string       `protobuf:"bytes,4,opt,name=os,proto3" json:"os,omitempty"`
	Platform             string       `protobuf:"bytes,5,opt,name=platform,proto3" json:"platform,omitempty"`
	PlatformFamily       string       `protobuf:"bytes,6,opt,name=platformFamily,proto3" json:"platformFamily,omitempty"`
	PlatformVersion      string       `protobuf:"bytes,7,opt,name=platformVersion,proto3" json:"platformVersion,omitempty"`
	KernelVersion        string       `protobuf:"bytes,8,opt,name=kernelVersion,proto3" json:"kernelVersion,omitempty"`
	VirtualizationSystem string       `protobuf:"bytes,9,opt,name=virtualizationSystem,proto3" json:"virtualizationSystem,omitempty"`
	VirtualizationRole   string       `protobuf:"bytes,10,opt,name=virtualizationRole,proto3" json:"virtualizationRole,omitempty"`
	Interfaces           []*Interface `protobuf:"bytes,111,rep,name=Interfaces,proto3" json:"Interfaces,omitempty"`
	InterfacesJSON       string
	CreateTime           int64
	UpdateTime           int64
	LastHeartBeatTime    int64
}

type Interface struct {
	Name         string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	HardwareAddr string   `protobuf:"bytes,2,opt,name=hardwareAddr,proto3" json:"hardwareAddr,omitempty"`
	IpAddrs      []string `protobuf:"bytes,3,rep,name=ipAddrs,proto3" json:"ipAddrs,omitempty"`
}

type HostService interface {
	Create(context.Context, *Host) error
	//IsExist(context.Context, agentID string) error
	Find(ctx context.Context, agentID string) (*Host, error)
}

func NewHostFromGRPC(req *message.RegisterRequest) *Host {
	h := &Host{
		AgentID:              req.AgentID,
		HostName:             req.HostName,
		HostUUID:             req.HostUUID,
		Os:                   req.Os,
		Platform:             req.Platform,
		PlatformFamily:       req.PlatformFamily,
		PlatformVersion:      req.PlatformVersion,
		KernelVersion:        req.KernelVersion,
		VirtualizationSystem: req.VirtualizationSystem,
		VirtualizationRole:   req.VirtualizationRole,
	}

	for _, i := range req.Interfaces {
		h.Interfaces = append(h.Interfaces, &Interface{
			IpAddrs:      i.IpAddrs,
			Name:         i.Name,
			HardwareAddr: i.HardwareAddr,
		})
	}

	intJSONBytes, err := json.Marshal(h.Interfaces)
	if err != nil {
		log.Errorf("使用 GRPC 接收到的接口信息转换为 JSON 时发生异常：%s", err)
		return h
	}

	h.InterfacesJSON = string(intJSONBytes)
	return h
}

//CheckHostChanceAcceptable 判断变更是否在可接受范围
func CheckHostChanceAcceptable(ori *Host, now *Host) (allSame bool, acceptable bool, msg string) {
	allSame = true
	changeCount := 0
	ok, cmsg := checkHostName(ori, now)
	if !ok {
		allSame = false
		changeCount++
		msg += cmsg
	}

	ok, cmsg = checkHostUUID(ori, now)
	if !ok {
		allSame = false
		changeCount++
		msg += cmsg
	}

	ok, cmsg = checkVirtualization(ori, now)
	if !ok {
		allSame = false
		changeCount++
		msg += cmsg
	}

	ok, cmsg = checkOS(ori, now)
	if !ok {
		allSame = false
		changeCount++
		msg += cmsg
	}
	//ok,cmsg =checkOSVersion(ori, now)
	ok, cmsg = checkInterfaces(ori, now)
	if !ok {
		allSame = false
		changeCount++
		msg += cmsg
	}

	if changeCount <= 1 {

		return allSame, true, msg
	}
	return allSame, false, msg
}

func checkHostName(ori *Host, now *Host) (same bool, msg string) {
	if ori.HostName == now.HostName {
		return true, ""
	}

	return false, fmt.Sprintf("主机名不一致：旧值为 %q 新值为 %q", ori.HostName, now.HostName)

}

func checkHostUUID(ori *Host, now *Host) (same bool, msg string) {

	if ori.HostUUID == now.HostUUID {
		return true, ""
	}

	return false, fmt.Sprintf("主机名不一致：旧值为 %q 新值为 %q", ori.HostUUID, now.HostUUID)
}

func checkVirtualization(ori *Host, now *Host) (same bool, msg string) {

	if ori.VirtualizationSystem == now.VirtualizationSystem &&
		ori.VirtualizationRole == now.VirtualizationRole {
		return true, ""
	}

	return false, fmt.Sprintf(
		"虚拟化环境不一致：旧值为 %s-%s 新值为 %s-%s",
		ori.VirtualizationSystem,
		ori.VirtualizationRole,
		now.VirtualizationSystem,
		now.VirtualizationRole)
}

func checkOS(ori *Host, now *Host) (same bool, msg string) {
	if ori.Os == now.Os &&
		ori.PlatformFamily == now.PlatformFamily &&
		ori.Platform == now.Platform {
		return true, ""
	}

	return false, fmt.Sprintf(
		"操作系统类型不一致：旧值为 %s-%s-%s 新值为 %s-%s-%s",
		ori.Os,
		ori.PlatformFamily,
		ori.Platform,
		now.Os,
		now.PlatformFamily,
		now.Platform)
}

func checkOSVersion(ori *Host, now *Host) (same bool, msg string) {
	if ori.PlatformVersion == now.PlatformVersion &&
		ori.KernelVersion == now.KernelVersion {
		return true, ""
	}

	return false, fmt.Sprintf(
		"操作系统版本不一致：旧值为 %s-%s 新值为 %s-%s",
		ori.PlatformVersion,
		ori.KernelVersion,
		now.PlatformVersion,
		now.KernelVersion)
}

func checkInterfaces(ori *Host, now *Host) (same bool, msg string) {
	var interfaceNameMap map[string]struct{}
	same = true
	for _, nn := range now.Interfaces {
		found := false
		interfaceNameMap[nn.Name] = struct{}{}
		nnJSON, _ := json.Marshal(nn)
		for _, on := range ori.Interfaces {
			if on.Name == nn.Name {
				found = true

				// 对比 Mac 和 IP
				if on.HardwareAddr != nn.HardwareAddr || util.StringSliceEqual(&on.IpAddrs, &nn.IpAddrs) {
					onJSON, _ := json.Marshal(on)
					msg = msg + fmt.Sprintf("网卡 %s 信息不匹配，旧值 %s 新值 %s", on.Name, onJSON, nnJSON)
					same = false
				}
				break
			}
		}

		if !found {
			same = false
			msg = msg + fmt.Sprintf("发现新网卡 %s ，值 %s", nn.Name, nnJSON)
		}
	}

	if len(ori.Interfaces) == len(now.Interfaces) {
		return same, msg
	}

	for _, on := range ori.Interfaces {
		_, ok := interfaceNameMap[on.Name]
		if !ok {
			same = false
			onJSON, _ := json.Marshal(on)
			msg = msg + fmt.Sprintf("未找到旧网卡 %s ，值 %s", on.Name, onJSON)
		}
	}
	return same, msg
}
