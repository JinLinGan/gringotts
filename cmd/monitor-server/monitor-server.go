package main

import (
	"context"
	"log"
	"net"

	"github.com/jinlingan/gringotts/message"

	"google.golang.org/grpc"
)

type gringottsServer struct {
}

func (s *gringottsServer) HeartBeat(context.Context, *message.HeartBeatRequest) (*message.HeartBeatResponse, error) {
	resp := message.HeartBeatResponse{}
	log.Println("called")
	return &resp, nil
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
