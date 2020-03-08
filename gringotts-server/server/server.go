package server

import (
	"context"
	"io"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"

	"github.com/jinlingan/gringotts/gringotts-server/model"

	"github.com/jinlingan/gringotts/gringotts-server/config"
	"github.com/jinlingan/gringotts/pkg/log"
	"github.com/jinlingan/gringotts/pkg/message"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// GringottsServer 服务器
type GringottsServer struct {
	sync.RWMutex
	grServer *grpc.Server
	config   *config.ServerConfig
	logger   log.Logger
	db       *gorm.DB
}

//NewServer 新建 Server 对象
func NewServer(cfg *config.ServerConfig, logger log.Logger) (*GringottsServer, error) {
	//TODO:移动到配置文件中
	dataSourceName := "gringotts:gringotts@tcp(mysql)/gringotts?parseTime=true"
	db, err := gorm.Open("mysql", dataSourceName)
	if err != nil {
		return nil, errors.Wrap(err, "open database error")
	}
	db = db.Debug()
	db.SingularTable(true)

	server := &GringottsServer{
		grServer: grpc.NewServer(),
		config:   cfg,
		logger:   logger,
		db:       db,
	}
	message.RegisterGringottsServer(server.grServer, server)
	return server, nil

}

//Serve 开始提供服务
func (s *GringottsServer) Serve() error {
	lsP := s.config.GetListenerPort()
	lis, err := net.Listen("tcp", lsP)
	if err != nil {
		return errors.Wrapf(err, "can not listen in port 0.0.0.0%s", lsP)
	}
	s.logger.Infof("gringotts server listen in port 0.0.0.0%s", lsP)
	return s.grServer.Serve(lis)
}

//HeartBeat 接收心跳
func (s *GringottsServer) HeartBeat(ctx context.Context,
	req *message.HeartBeatRequest) (*message.HeartBeatResponse, error) {
	s.logger.Debugf("get HeartBeat message from agent(id=%s,hostname=%s)", req.GetAgentId(), req.GetHostName())
	return s.newHeartBeatResponse(), nil
}

//DownloadFile 下载文件
func (s *GringottsServer) DownloadFile(f *message.File, fs message.Gringotts_DownloadFileServer) error {
	//TODO:改变文件路径
	rf, err := os.Open("/Users/jinlin/code/golang/src/github.com/jinlingan/gringotts/testfile/" + f.GetFileName())

	if err != nil {
		return err
	}
	buff := make([]byte, 500)
	for {
		n, err := rf.Read(buff)
		if err != nil && err != io.EOF {
			return status.Errorf(500, "read file error %s", err)
		}
		exit := false
		if err == io.EOF {
			exit = true
		}

		err = fs.Send(&message.FileChunk{
			Data: buff[:n],
		})
		if err != nil {
			return status.Errorf(500, "send file Chunk error %s", err)
		}

		if exit {
			break
		}
	}
	if rf.Close() != nil {
		return err
	}
	return nil
}

// Register agent 注册 - 暂未实现
func (s *GringottsServer) Register(ctx context.Context,
	req *message.RegisterRequest) (*message.RegisterResponse, error) {

	//TODO:返回正确的 ConfigVersion

	//hosts,notFindNetInfo, err := s.findHost(req.NetInfo)
	hosts, _, err := s.findHost(req.NetInfo)
	if err != nil {
		return nil, errors.Wrap(err, "look up hosts by interface info fail")
	}
	if len(hosts) > 1 {
		return nil, errors.Errorf("find multiple host: %s", spew.Sdump(hosts))
	}

	if len(hosts) == 1 {
		//如果请求中的 HostId == 0 ，说明客户端是新注册的，那找到主机就是不对的
		if req.HostId == "" {
			return nil, errors.Errorf("find registered host: %s", spew.Sdump(hosts))
		}
		//如果只返回一个主机并且 id 一样，注册成功
		if req.HostId == strconv.FormatUint(uint64(hosts[0].ID), 10) {
			//TODO:是否要更新主机信息？
			return &message.RegisterResponse{AgentId: req.HostId, ConfigVersion: "0"}, nil
		}
	}

	//没有找到主机的情况
	//如果没有发送 ID 说明是新注册
	if req.HostId == "" {
		//注册主机
		h := s.GetHostByRegisterReq(req)
		h, err := s.RegisterHost(h)
		if err != nil {
			return nil, errors.Wrapf(err, "register host %s fail", spew.Sdump(req))
		}
		return &message.RegisterResponse{
			AgentId:       strconv.FormatUint(uint64(h.ID), 10),
			ConfigVersion: "0",
		}, nil
	}

	//发送 ID 说明是续租，但是没有找到对应的主机说明有问题
	return nil, errors.Errorf("can not find host by network interface %s", spew.Sdump(req))

}

// GetHostByRegisterReq 从注册请求中获取主机对象
func (s *GringottsServer) GetHostByRegisterReq(req *message.RegisterRequest) *model.Host {
	h := &model.Host{HostName: req.GetHostName()}
	for _, value := range req.NetInfo {
		h.HostInterface = append(
			h.HostInterface,
			&model.HostInterface{
				HWAddr:        value.MacAddress,
				IPAddress:     value.IpAddress,
				InterfaceName: value.InterfaceName,
			})
	}
	return h
}

// RegisterHost 注册主机
func (s *GringottsServer) RegisterHost(host *model.Host) (*model.Host, error) {

	if r := s.db.Create(host); r.Error != nil {
		return nil, errors.Wrapf(r.Error, "register host %s fail", spew.Sdump(host))
	}
	return host, nil

}

// findHost 使用网卡信息查找主机，返回主机对象以及未找到匹配项的网卡信息
func (s *GringottsServer) findHost(
	net []*message.RegisterRequest_NetInfo,
) (
	[]*model.Host,
	[]*message.RegisterRequest_NetInfo,
	error,
) {

	notFind := make([]*message.RegisterRequest_NetInfo, len(net))
	var hosts []*model.Host

	for _, value := range net {
		h, err := s.findHostsByInterface(value)
		if err != nil {
			return nil, nil, err
		}
		if h != nil {
			hosts = appendHosts(hosts, h...)
		} else {
			notFind = append(notFind, value)
		}
	}
	return hosts, notFind, nil
}

// appendHost 添加 host 到 slice 中，并且避免重复
func appendHost(hostSlice []*model.Host, host *model.Host) []*model.Host {
	//TODO:是否要判断 nil？
	//if host == nil {
	//	return hostSlice
	//}
	find := false
	for _, value := range hostSlice {
		if value.ID == host.ID {
			find = true
			break
		}
	}

	if !find {
		return append(hostSlice, host)
	}
	return hostSlice
}

// appendHost 添加 host 到 slice 中，并且避免重复
func appendHosts(hostSlice []*model.Host, hosts ...*model.Host) []*model.Host {
	for _, value := range hosts {
		hostSlice = appendHost(hostSlice, value)
	}
	return hostSlice
}

// findHostByInterface 使用网卡信息查找主机
func (s *GringottsServer) findHostsByInterface(inf *message.RegisterRequest_NetInfo) ([]*model.Host, error) {
	var results []*model.Host

	if r := s.db.Preload("HostInterface").
		Table("host").
		Joins("left join host_interface on host_interface.host_id = host.id").
		Where("host_interface.hw_addr = ?", inf.MacAddress).
		Find(&results); r.Error != nil {

		return nil, errors.Wrapf(r.Error, "find hosts by mac address %q fail", inf.MacAddress)
	}
	return results, nil
}

func (s *GringottsServer) newHeartBeatResponse() *message.HeartBeatResponse {
	resp := &message.HeartBeatResponse{
		ServerId:      s.config.GetExternalAddress(),
		ConfigVersion: strconv.Itoa(time.Now().Minute()),
		MonitorInfo:   getAllTaskByAgentID(),
	}
	return resp

}

func getAllTaskByAgentID() *message.MonitorInfo {
	taskInfoOne := message.MonitorInfo{
		Items: []*message.MonitorItem{
			{
				TaskId:             1,
				ExecIntervalSecond: 10,
				Type:               message.MonitorItemType_SELF,
				SelfFunc:           message.SelfMonitorFunc_CPU,
			},
		},
	}

	taskInfoTwo := message.MonitorInfo{
		Items: []*message.MonitorItem{
			{
				TaskId:             2,
				ExecIntervalSecond: 10,
				Type:               message.MonitorItemType_SELF,
				SelfFunc:           message.SelfMonitorFunc_MEM,
			},
		},
	}

	if time.Now().Minute()%2 == 0 {
		return &taskInfoTwo
	}

	return &taskInfoOne
}
