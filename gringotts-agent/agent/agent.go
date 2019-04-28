package agent

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jinlingan/gringotts/gringotts-agent/communication"
	"github.com/jinlingan/gringotts/gringotts-agent/config"
)

// Agent Gringotts Agent
type Agent struct {
	cfg       *config.AgentConfig
	apiClient *communication.Client
}

// NewAgent 新建 Agent
func NewAgent(cfg *config.AgentConfig) *Agent {

	//新建客户端
	client, err := communication.NewClient(cfg)
	if err != nil {
		log.Printf("can not create communicate agent with server %s ,err is %s", cfg.GetServerAddress(), err)
	}
	return &Agent{
		cfg:       cfg,
		apiClient: client,
	}
}

// Start 启动 Agent
func (a *Agent) Start() error {

	stop := make(chan int, 1)

	// 如果 Agent 还没有注册
	if !a.cfg.IsRegistered() {

		// 启动注册流程
		if err := a.regist(); err != nil {
			return fmt.Errorf("regist agent to server %s error: %s", a.cfg.GetServerAddress(), err)
		}
	}

	// if err := downloadFile(client, "main", "aaaa", a.cfg.GetExecuterPath()); err != nil {
	// 	log.Printf("can not download file from  server %s ,err is %s", a.cfg.GetServerAddress(), err)
	// }
	//开始发送心跳
	go a.sendHeartBeat()
	<-stop
	return nil
}

// regist 注册 agent
func (a *Agent) regist() error {
	//TODO finish this
	_, err := a.apiClient.Regist("", nil)
	return err
}

// func (a *Agent) downloadFile(client *communication.Client, filename, sha1, destPath string) error {
// 	return client.DownloadFile(filename, sha1, destPath)
// }

func (a *Agent) sendHeartBeat() {

	ticker := time.NewTicker(5 * time.Second)
	for {

		// set timer
		start := time.Now()
		//send HeartBeat
		ctx, cancle := context.WithTimeout(context.Background(), time.Second*30)
		r, err := a.apiClient.HeartBeat(ctx, a.cfg.GetAgentID())
		//
		log.Printf("send HeartBeat (%s)", time.Since(start))
		if err != nil {
			log.Printf("send HeartBeat with err: %v", err)

		} else if a.cfg.GetConfigVersion() != r.ConfigVersion {

			log.Printf("get HeartBeat response from server(id=%s) with config version = %d", r.ServerId, r.ConfigVersion)
			log.Printf("not equal local version %d , reload", a.cfg.GetConfigVersion())
			// processConfig(r.MonitorInfo)
			if a.cfg.SetConfigVersion(r.ConfigVersion) != nil {
				log.Printf("set config version error: %s", err)
			}
			//TODO: stop agent

		}
		<-ticker.C
		cancle()

	}
}

// func grpcSpeedTest(rpcDurations *prometheus.SummaryVec) {
// 	count := 1000

// 	conn, err := grpc.Dial(address, grpc.WithInsecure())
// 	if err != nil {
// 		log.Fatalf("did not connect: %v", err)
// 	}
// 	server := message.NewGringottsa.apiClient(conn)
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
// 			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.91: 0.01,
//											0.92: 0.01, 0.93: 0.01, 0.94: 0.01,
//											0.95: 0.01, 0.96: 0.01, 0.97: 0.01,
//											0.98: 0.01, 0.99: 0.001},
// 		},
// 		[]string{"service"},
// 	)
// 	prometheus.MustRegister(rpcDurations)
// 	go func() {
// 		log.Fatal(http.ListenAndServe(":9999", nil))
// 	}()
// 	grpcSpeedTest(rpcDurations)
// }

// var allTaskStopSignal chan int

// var doing = false

// func processConfig(info *message.MonitorInfo) {
// 	if doing {
// 		close(allTaskStopSignal)
// 		log.Println("Close all running task")
// 	}
// 	allTaskStopSignal = startAllAgent(info)
// }

// func startAllAgent(infos *message.MonitorInfo) chan int {
// 	// wg := sync.WaitGroup{}
// 	// wg.Add(len(infos.GetItems()))
// 	doing = true
// 	stop := make(chan int, 1)
// 	for _, i := range infos.GetItems() {
// 		go func(t *message.MonitorItem) {
// 			tick := time.NewTicker(time.Duration(t.ExecIntervalSecond) * time.Second)
// 			for {
// 				select {
// 				case <-tick.C:
// 					log.Printf("Task %d , get %s info %d \n", t.TaskId, t.SelfFunc, time.Now().Second())
// 				case <-stop:
// 					// wg.Done()
// 					log.Printf("Stop Task ID %d", t.TaskId)
// 					return
// 				}
// 			}
// 		}(i)
// 	}
// 	// wg.Wait()
// 	return stop

// }
