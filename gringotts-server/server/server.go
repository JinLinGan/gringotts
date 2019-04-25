package server

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	"github.com/jinlingan/gringotts/message"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type GringottsServer struct {
	listenAddress string
	grpcServer    *grpc.Server
	serverID      string
}

func NewServer(address string, serverID string) (*GringottsServer, error) {

	server := &GringottsServer{
		listenAddress: address,
		grpcServer:    grpc.NewServer(),
	}
	message.RegisterGringottsServer(server.grpcServer, server)
	return server, nil

}

func (s *GringottsServer) Serve() error {
	lis, err := net.Listen("tcp", s.listenAddress)
	if err != nil {
		return fmt.Errorf("can not listen in port %s", s.listenAddress)
	}
	return s.grpcServer.Serve(lis)
}
func (s *GringottsServer) HeartBeat(ctx context.Context, req *message.HeartBeatRequest) (*message.HeartBeatResponse, error) {
	log.Printf("get HeartBeat message from agent(id=%s,hostname=%s)", req.GetAgnetId(), req.GetHostName())
	return s.newHeartBeatResponse(), nil
}
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

func (s *GringottsServer) newHeartBeatResponse() *message.HeartBeatResponse {
	resp := &message.HeartBeatResponse{
		ServerId:      s.serverID,
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
