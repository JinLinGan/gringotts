package cmd

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/jinlingan/gringotts/communication"
	"github.com/jinlingan/gringotts/config"
	"github.com/jinlingan/gringotts/message"

	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start agent",
	// Long: `start agent`,
	RunE: startAgent,
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.PersistentFlags().StringP(WORK_DIR_FLAG_NAME, "w", "/usr/local/gringotts", "workdir used to save all program files")
	startCmd.PersistentFlags().StringP(SERVER_ADDRESS_FLAG_NAME, "s", "127.0.0.1:7777", "server address")
}

const (
	WORK_DIR_FLAG_NAME       = "workdir"
	SERVER_ADDRESS_FLAG_NAME = "server"
)

func startAgent(cmd *cobra.Command, args []string) error {

	stop := make(chan int, 1)
	// 根据命令行参数设置工作目录
	w, err := cmd.Flags().GetString(WORK_DIR_FLAG_NAME)
	if err != nil {
		return err
	}
	if err := config.SetWorkDir(w); err != nil {
		return err
	}

	// 根据命令行参数设置服务端地址
	s, err := cmd.Flags().GetString(SERVER_ADDRESS_FLAG_NAME)
	if err != nil {
		return err
	}
	config.ServerAddress = s

	//新建客户端
	client, err := communication.NewClient(config.ServerAddress)
	if err != nil {
		log.Printf("can not communicate with server %s ,err is %s", config.ServerAddress, err)
	}

	if err := downloadFile(client, "main", "aaaa"); err != nil {
		log.Printf("can not download file from  server %s ,err is %s", config.ServerAddress, err)
	}

	//开始发送心跳
	go sendHeartBeat(client)
	<-stop
	return nil
}

func downloadFile(client *communication.Client, filename, sha1 string) error {
	return client.DownloadFile(filename, sha1, "./depend")
}

func sendHeartBeat(client *communication.Client) {

	ticker := time.NewTicker(5 * time.Second)
	for {

		// set timer
		start := time.Now()
		//send HeartBeat
		ctx, cancle := context.WithTimeout(context.Background(), time.Second*30)
		r, err := client.HeartBeat(ctx, agentID)
		//
		log.Printf("send HeartBeat (%s)", time.Since(start))
		if err != nil {
			log.Printf("send HeartBeat with err: %v", err)

		} else {

			if configVersion != r.ConfigVersion {
				log.Printf("get HeartBeat response from server(id=%s) with config version = %d", r.ServerId, r.ConfigVersion)
				log.Printf("not equal local version %d , reload", configVersion)
				processConfig(r.MonitorInfo)
				configVersion = r.ConfigVersion
			}
		}
		<-ticker.C
		cancle()

	}
}

var (
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
			tick := time.NewTicker(time.Duration(t.ExecIntervalSecond) * time.Second)
			for {
				select {
				case <-tick.C:
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
