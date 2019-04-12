package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/jinlingan/gringotts/message"
	"google.golang.org/grpc"
)

const (
	serverId = 99
)

type gringottsServer struct {
}

func (s *gringottsServer) HeartBeat(ctx context.Context, req *message.HeartBeatRequest) (*message.HeartBeatResponse, error) {
	log.Printf("get HeartBeat message from agent(id=%s,hostname=%s)", req.GetAgnetId(), req.GetHostName())
	return newHeartBeatResponse(), nil
}

func main() {
	lis, err := net.Listen("tcp", ":7777")
	if err != nil {
		log.Fatalf("can not listen in port :7777")
	}
	grpcServer := grpc.NewServer()
	message.RegisterGringottsServer(grpcServer, &gringottsServer{})
	grpcServer.Serve(lis)
}

func newHeartBeatResponse() *message.HeartBeatResponse {
	resp := &message.HeartBeatResponse{
		ServerId:      serverId,
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
