package main

import (
	"context"
	"log"
	"math/rand"
	"os"
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

var (
	address = "localhost:7777"
	agentID = string(rand.Intn(9999))
)

var server message.GringottsClient
var configVersion int64

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	server = message.NewGringottsClient(conn)
	stop := make(chan int, 1)
	go func() {
		for {
			select {
			case <-time.Tick(5 * time.Second):
				ctx, cancle := context.WithTimeout(context.Background(), time.Second*3)
				defer cancle()
				r, err := server.HeartBeat(ctx, newHeartBeatRequest())
				if err != nil {
					log.Printf("send HeartBeat with err: %v", err)
					continue
				}
				if configVersion != r.ConfigVersion {
					log.Printf("get HeartBeat response from server(id=%d) with config version = %d", r.ServerId, r.ConfigVersion)
					log.Printf("not equal local version %d , reload", configVersion)
					processConfig(r.MonitorInfo)
					configVersion = r.ConfigVersion
				}
				// else {
				// 	log.Printf(" equal local version %d , skip", configVersion)
				// }
			}
		}
	}()

	<-stop
}

var allTaskStopSignal chan int

var doing = false

func processConfig(info *message.MonitorInfo) {
	if doing {
		close(allTaskStopSignal)
		log.Println("Close all running task")
	}
	allTaskStopSignal = startAllAgent(info)
}

func startAllAgent(infos *message.MonitorInfo) chan int {
	// wg := sync.WaitGroup{}
	// wg.Add(len(infos.GetItems()))
	doing = true
	stop := make(chan int, 1)
	for _, i := range infos.GetItems() {
		go func(t *message.MonitorItem) {
			for {
				select {
				case <-time.Tick(time.Duration(i.ExecIntervalSecond) * time.Second):
					log.Printf("Task %d , get %s info %d \n", t.TaskId, t.SelfFunc, time.Now().Second())
				case <-stop:
					// wg.Done()
					log.Printf("Stop Task ID %d", t.TaskId)
					return
				}
			}
		}(i)
	}
	// wg.Wait()
	return stop

}

func newHeartBeatRequest() *message.HeartBeatRequest {
	req := message.HeartBeatRequest{
		AgnetId: agentID,
		Time:    time.Now().UnixNano(),
	}
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unkonw"
		log.Printf("get hostname with err: %s", err)
	}
	req.HostName = hostname
	return &req
}
