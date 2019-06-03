package agent

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/jinlingan/gringotts/gringotts-agent/communication"
	"github.com/jinlingan/gringotts/gringotts-agent/config"
	"github.com/jinlingan/gringotts/gringotts-agent/log"
	"github.com/jinlingan/gringotts/gringotts-agent/model"
)

// Agent Gringotts Agent
type Agent struct {
	sync.RWMutex
	cfg          *config.AgentConfig
	apiClient    *communication.Client
	logger       log.Logger
	isRegistered bool
	agentInfo    *agentRunningInfo
}

// agentRunningInfo 保存了 agentID 和 配置版本
type agentRunningInfo struct {
	agentID       string
	configVersion int64
}

// NewAgent 新建 Agent
func NewAgent(cfg *config.AgentConfig, logger log.Logger) *Agent {

	//新建客户端
	client, err := communication.NewClient(cfg, logger)
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
	//TODO:移动到启动时判断
	agentInfo, err := a.getAgentIDFormWorkdir()
	if err != nil {
		a.logger.Info("read agent info failed so set state unregistered")
		a.isRegistered = false
	} else {
		a.isRegistered = true
		a.agentInfo = agentInfo
	}

	// 如果 Agent 还没有注册
	if !a.isRegistered {

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
		r, err := a.apiClient.HeartBeat(ctx, a.GetAgentID())
		//
		a.logger.Debugf("send HeartBeat (%s)", time.Since(start))
		if err != nil {
			a.logger.Errorf("send HeartBeat with err: %v", err)

		} else if a.GetConfigVersion() != r.ConfigVersion {

			a.logger.Infof("get HeartBeat response from server(id=%s) with config version = %d", r.ServerId, r.ConfigVersion)
			a.logger.Infof("not equal local version %d , reload", a.GetConfigVersion())
			// processConfig(r.MonitorInfo)
			if a.SetConfigVersion(r.ConfigVersion) != nil {
				a.logger.Errorf("set config version error: %s", err)
			}
			//TODO: stop agent

		}
		<-ticker.C
		cancel()

	}
}

func (a *Agent) getAgentIDFormWorkdir() (*agentRunningInfo, error) {

	path := a.cfg.GetAgentRunningInfoFilePath()

	agentInfo := new(agentRunningInfo)

	b, err := ioutil.ReadFile(filepath.Clean(path))
	if err != nil {
		// 如果文件不存在
		if os.IsNotExist(err) {
			return nil, errors.Wrapf(err, "agent info can not find")
		}
		return nil, errors.Wrapf(err, "read agent info failed")
	}
	if err := json.Unmarshal(b, agentInfo); err != nil {
		return nil, errors.Wrapf(err, "decode agent info file %s fail", path)
	}
	return agentInfo, nil
}

// GetAgentID 获取 AgentID
func (a *Agent) GetAgentID() string {
	a.RLock()
	defer a.RUnlock()
	return a.agentInfo.agentID
}

// GetConfigVersion 获取配置版本
func (a *Agent) GetConfigVersion() int64 {
	a.RLock()
	defer a.RUnlock()
	return a.agentInfo.configVersion
}

// SetConfigVersion 设置配置版本
func (a *Agent) SetConfigVersion(v int64) error {

	a.Lock()
	a.agentInfo.configVersion = v
	a.Unlock()

	b, err := json.Marshal(a.agentInfo)
	if err != nil {
		return errors.Wrapf(err, "encode agent config %+v fail", a.agentInfo)
	}
	path := a.cfg.GetAgentRunningInfoFilePath()
	if err := ioutil.WriteFile(path, b, config.PermissionMode); err != nil {
		return errors.Wrapf(err, "write agent config file %s fail", path)
	}
	return nil
}

// IsRegistered Agent 是否成功注册
func (a *Agent) IsRegistered() bool {
	a.RLock()
	defer a.RUnlock()
	return a.isRegistered
}
