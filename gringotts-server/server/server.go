package server

import (
	"context"
	"io"
	"net"
	"os"
	"time"

	"github.com/jinlingan/gringotts/gringotts-server/model"

	"github.com/jinlingan/gringotts/common/log"
	"github.com/jinlingan/gringotts/common/message"
	"github.com/jinlingan/gringotts/gringotts-server/config"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// GringottsServer 服务器
type GringottsServer struct {
	grServer *grpc.Server
	config   *config.ServerConfig
	logger   log.Logger
	db       *gorm.DB
}

//NewServer 新建 Server 对象
func NewServer(cfg *config.ServerConfig, logger log.Logger) (*GringottsServer, error) {
	dataSourceName := "gringotts:gringotts@tcp(127.0.0.1)/gringotts?parseTime=true"
	db, err := gorm.Open("mysql", dataSourceName)
	if err != nil {
		return nil, errors.Wrapf(err, "can not connect to database %s", dataSourceName)
	}
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

	return nil, errors.Errorf("not implement")
}

// findHost 使用网卡信息查找主机，返回主机对象以及未找到匹配项的网卡信息
func (s *GringottsServer) findHost(net []*message.RegisterRequest_NetInfo) ([]*model.Host, []*message.RegisterRequest_NetInfo, error) {
	notFind := make([]*message.RegisterRequest_NetInfo, len(net))
	hosts := make([]*model.Host, 5)

	for _, value := range net {
		h, err := s.findHostByInterface(value)
		if err != nil {
			return nil, nil, err
		}
		if h != nil {
			hosts = appendHost(hosts, h)
		} else {
			notFind = append(notFind, value)
		}
	}
	return hosts, notFind, nil
}

// appendHost 添加 host 到 slice 中，并且避免重复
func appendHost(hosts []*model.Host, host *model.Host) []*model.Host {
	find := false
	for _, value := range hosts {
		if value.ID == host.ID {
			find = true
			break
		}
	}

	if find == false {
		return append(hosts, host)
	}
}

// findHostByInterface 使用网卡信息查找主机
func (s *GringottsServer) findHostByInterface(inf *message.RegisterRequest_NetInfo) ([]*model.Host, error) {
	count := 0
	//TODO:写实现
	//s.db.Model(&model.HostInterface{}).Where("HWAddr = ?", inf.MacAddress).Count(&count)
	return nil, nil
}

func (s *GringottsServer) newHeartBeatResponse() *message.HeartBeatResponse {
	resp := &message.HeartBeatResponse{
		ServerId:      s.config.GetExternalAddress(),
		ConfigVersion: int64(time.Now().Minute()),
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
