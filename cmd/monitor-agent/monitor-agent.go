package main

import (
	"context"
	"log"
	"time"

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

const (
	address     = "localhost:7777"
	defaultName = "world"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := message.NewGringottsClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.HeartBeat(ctx, &message.HeartBeatRequest{})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.String())
}
