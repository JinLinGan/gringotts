package server

import (
	"context"
	"io"
	"net"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/pkg/errors"

	"github.com/jinlingan/gringotts/gringotts-server/config"
	"github.com/jinlingan/gringotts/message"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// GringottsServer 服务器
type GringottsServer struct {
	grpcServer *grpc.Server
	config     *config.ServerConfig
}

//NewServer 新建 Server 对象
func NewServer(cfg *config.ServerConfig) (*GringottsServer, error) {

	server := &GringottsServer{
		grpcServer: grpc.NewServer(),
		config:     cfg,
	}
	message.RegisterGringottsServer(server.grpcServer, server)
	return server, nil

}

//Serve 开始提供服务
func (s *GringottsServer) Serve() error {
	lis, err := net.Listen("tcp", s.config.GetListenerPort())
	if err != nil {
		return errors.Errorf("can not listen in port 0.0.0.0%s", s.config.GetListenerPort())
	}
	return s.grpcServer.Serve(lis)
}

//HeartBeat 接收心跳
func (s *GringottsServer) HeartBeat(ctx context.Context,
	req *message.HeartBeatRequest) (*message.HeartBeatResponse, error) {
	log.Printf("get HeartBeat message from agent(id=%s,hostname=%s)", req.GetAgentId(), req.GetHostName())
	return s.newHeartBeatResponse(), nil
}

//DownloadFile 下载文件
func (s *GringottsServer) DownloadFile(f *message.File, fs message.Gringotts_DownloadFileServer) error {
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
func (s *GringottsServer) Register(
	ctx context.Context,
	req *message.RegisterRequest) (*message.RegisterResponse, error) {
	return nil, errors.Errorf("not implement")
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
