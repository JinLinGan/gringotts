package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/jinlingan/gringotts/communication"
	"github.com/jinlingan/gringotts/message"
)

var (
	address = "localhost:7777"
	agentID = string(rand.Intn(9999))
)
var configVersion int64

// func grpcSpeedTest(rpcDurations *prometheus.SummaryVec) {
// 	count := 1000

// 	conn, err := grpc.Dial(address, grpc.WithInsecure())
// 	if err != nil {
// 		log.Fatalf("did not connect: %v", err)
// 	}
// 	server := message.NewGringottsClient(conn)
// 	for i := 1; i <= count; i++ {
// 		go func() {
// 			for true {
// 				start := time.Now()
// 				server.HeartBeat(context.Background())
// 				rpcDurations.WithLabelValues("normal").Observe(float64(time.Since(start).Nanoseconds()))
// 			}
// 		}()
// 	}
// }

// func performanceTest() {
// 	http.Handle("/metrics", promhttp.Handler())
// 	rpcDurations := prometheus.NewSummaryVec(
// 		prometheus.SummaryOpts{
// 			Name:       "grpc_durations_seconds",
// 			Help:       "GRPC latency distributions.",
// 			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.91: 0.01, 0.92: 0.01, 0.93: 0.01, 0.94: 0.01, 0.95: 0.01, 0.96: 0.01, 0.97: 0.01, 0.98: 0.01, 0.99: 0.001},
// 		},
// 		[]string{"service"},
// 	)
// 	prometheus.MustRegister(rpcDurations)
// 	go func() {
// 		log.Fatal(http.ListenAndServe(":9999", nil))
// 	}()
// 	grpcSpeedTest(rpcDurations)
// }

func main() {

	// stop := make(chan int, 1)
	// client, err := communication.NewClient(address)
	// if err != nil {
	// 	log.Printf("can not communicate with server %s ,err is %s", address, err)
	// }
	// go sendHeartBeat(client)

	// <-stop
	client, err := communication.NewClient(address)
	if err != nil {
		log.Printf("can not communicate with server %s ,err is %s", address, err)
	}

	if err := downloadFile(client, "vgo", "aaaa"); err != nil {
		log.Print(err)
	}

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
				case <-time.Tick(time.Duration(t.ExecIntervalSecond) * time.Second):
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

func sendHeartBeat(client *communication.Client) {
	ctx, cancle := context.WithTimeout(context.Background(), time.Second*3)
	defer cancle()
	for {

		<-time.Tick(5 * time.Second)

		// set timer
		start := time.Now()
		//send HeartBeat
		r, err := client.HeartBeat(ctx, agentID)
		//
		log.Printf("send HeartBeat (%s)", time.Since(start))
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
func downloadFile(client *communication.Client, filename, sha1 string) error {
	return client.DownloadFile(filename, sha1, "./depend")
}
