package agent

import (
	"context"
	"time"

	"github.com/jinlingan/gringotts/gringotts-agent/log"

	"github.com/pkg/errors"

	"github.com/jinlingan/gringotts/gringotts-agent/model"

	"github.com/jinlingan/gringotts/gringotts-agent/communication"
	"github.com/jinlingan/gringotts/gringotts-agent/config"
)

// Agent Gringotts Agent
type Agent struct {
	cfg       *config.AgentConfig
	apiClient *communication.Client
	logger    log.Logger
}

// NewAgent 新建 Agent
func NewAgent(cfg *config.AgentConfig, logger log.Logger) *Agent {

	//新建客户端
	client, err := communication.NewClient(cfg)
	if err != nil {
		logger.Warn(errors.Wrapf(err, "can not create communicate agent with server %s", cfg.GetServerAddress()))
	}
	return &Agent{
		cfg:       cfg,
		apiClient: client,
		logger:    logger,
	}
}

// Start 启动 Agent
func (a *Agent) Start() error {

	stop := make(chan int, 1)

	// 如果 Agent 还没有注册
	if !a.cfg.IsRegistered() {

		// 启动注册流程
		//TODO:重试 N 次
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
	//TODO 写代码
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
		a.logger.Debugf("send HeartBeat (%s)", time.Since(start))
		if err != nil {
			a.logger.Errorf("send HeartBeat with err: %v", err)

		} else if a.cfg.GetConfigVersion() != r.ConfigVersion {

			a.logger.Infof("get HeartBeat response from server(id=%s) with config version = %d", r.ServerId, r.ConfigVersion)
			a.logger.Infof("not equal local version %d , reload", a.cfg.GetConfigVersion())
			// processConfig(r.MonitorInfo)
			if a.cfg.SetConfigVersion(r.ConfigVersion) != nil {
				a.logger.Errorf("set config version error: %s", err)
			}
			//TODO: stop agent

		}
		<-ticker.C
		cancel()

	}
}
