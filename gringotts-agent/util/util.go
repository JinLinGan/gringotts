package util

import (
	"net"

	"github.com/jinlingan/gringotts/gringotts-agent/model"
)

// GetNIC 获取主机网卡信息
func GetNIC() ([]*model.NICInfo, error) {
	var nics []*model.NICInfo
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return nil, err
		}
		name := i.Name
		mac := i.HardwareAddr.String()
		//忽略 mac 地址为空的网卡
		if mac == "" {
			continue
		}
		for _, a := range addrs {
			//取 ip 地址，忽略回环地址
			if ipNet, ok := a.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
				nics = append(nics, &model.NICInfo{
					Name:       name,
					MacAddress: mac,
					IPAddress:  ipNet.IP.String(),
				})
			}
		}
	}
	return nics, nil
}
