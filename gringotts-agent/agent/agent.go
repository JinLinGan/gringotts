package agent

import (
	"context"
	"log"
	"time"

	"github.com/pkg/errors"

	"github.com/jinlingan/gringotts/gringotts-agent/model"

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
		if err := a.register(); err != nil {
			return errors.Wrapf(err, "register agent to server %s fail", a.cfg.GetServerAddress())
		}
	}

	//开始发送心跳
	go a.sendHeartBeat()
	<-stop
	return nil
}

// register 注册 agent
func (a *Agent) register() error {
	//TODO finish this
	_, err := a.apiClient.Register("aaaa", &model.NetInfos{
		"eth0": {},
	})
	return err
}

func (a *Agent) sendHeartBeat() {

	ticker := time.NewTicker(5 * time.Second)
	for {

		// set timer
		start := time.Now()
		//send HeartBeat
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
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
		cancel()

	}
}
