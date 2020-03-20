package agent

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/jinlingan/gringotts/gringotts-agent/communication"
	"github.com/jinlingan/gringotts/gringotts-agent/config"
	"github.com/jinlingan/gringotts/gringotts-agent/model"
	"github.com/jinlingan/gringotts/pkg/log"
	"github.com/jinlingan/gringotts/pkg/metadata/host"
	"github.com/pkg/errors"
)

// Agent Gringotts Agent
type Agent struct {
	sync.RWMutex
	cfg          *config.AgentConfig
	apiClient    *communication.Client
	logger       log.Logger
	isRegistered bool
	agentInfo    agentRunningInfo
}

//var _ io.Writer = &Agent{}

// 注册超时时间
const registerTimeOut = time.Second * 60

// 注册间隔
const registerInterval = time.Second * 5

// agentRunningInfo 保存了 agentID 和 配置版本
type agentRunningInfo struct {
	AgentID       string
	ConfigVersion string
}

// NewAgent 新建 Agent
func NewAgent(cfg *config.AgentConfig, logger log.Logger) *Agent {

	//新建客户端
	client, err := communication.NewClient(cfg, logger)
	if err != nil {
		logger.Warne(err, "can not create communicate agent with server %s", cfg.GetServerAddress())
	}
	return &Agent{
		cfg:       cfg,
		apiClient: client,
		logger:    logger,
	}
}

func (a *Agent) register() (*model.RegisterResp, error) {
	info := host.GetHostInfo()

	t := time.After(registerTimeOut)
	for {
		select {
		case <-t:
			return nil, errors.Errorf("注册超时，当前超时时间为 %.f 秒", registerTimeOut.Seconds())
		default:
			resp, err := a.apiClient.Register(a.agentInfo.AgentID, info)
			if err == nil {
				return resp, err
			}

			a.logger.Errorf("注册失败，等待 %.f 秒后重试: %s", registerInterval.Seconds(), err)
			time.Sleep(registerInterval)
		}
	}
}

// Start 启动 Agent
func (a *Agent) Start() error {
	stop := make(chan int, 1)
	agentInfo, err := a.getAgentIDFormWorkDir()

	if err != nil {
		a.logger.Infof("read agent info failed so set state unregistered. Caused by: %v", err)
		a.isRegistered = false
	} else {
		a.isRegistered = true
		a.agentInfo = *agentInfo
	}
	resp, err := a.register()
	if err != nil {
		a.logger.Errorf("register agent to server %s fail", a.cfg.GetServerAddress())
		return errors.New("start agent fail")
	}

	a.logger.Infof("获取到 AgentID=%s，ConfigVersion=%s", resp.AgentID, resp.ConfigVersion)

	err = a.saveRegisterInfo(resp)

	if err != nil {
		a.logger.Errorf("保存注册信息失败. Caused by: %v", err)
		return errors.New("start agent fail")
	}
	//开始发送心跳
	go a.sendHeartBeat()
	<-stop
	return nil
}

func (a *Agent) saveRegisterInfo(info *model.RegisterResp) error {
	path := a.cfg.GetAgentRunningInfoFilePath()

	agentInfo := agentRunningInfo{
		AgentID:       info.AgentID,
		ConfigVersion: info.ConfigVersion,
	}
	b, err := json.Marshal(agentInfo)
	if err != nil {
		return errors.Wrap(err, "can not marshal agentInfo")
	}
	path = filepath.Clean(path)

	err = ioutil.WriteFile(path, b, 0644)
	if err != nil {
		return errors.Wrapf(err, "can not write file %s", path)
	}

	a.agentInfo = agentInfo
	return nil
}

func (a *Agent) sendHeartBeat() {

	ticker := time.NewTicker(5 * time.Second)
	for {

		// set timer
		start := time.Now()
		//send HeartBeat
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		r, err := a.apiClient.HeartBeat(ctx, a.GetAgentID())
		cancel()

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

	}
}

func (a *Agent) getAgentIDFormWorkDir() (*agentRunningInfo, error) {

	path := a.cfg.GetAgentRunningInfoFilePath()

	agentInfo := &agentRunningInfo{}

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
	return a.agentInfo.AgentID
}

// GetConfigVersion 获取配置版本
func (a *Agent) GetConfigVersion() string {
	a.RLock()
	defer a.RUnlock()
	return a.agentInfo.ConfigVersion
}

// SetConfigVersion 设置配置版本
func (a *Agent) SetConfigVersion(v string) error {

	a.Lock()
	a.agentInfo.ConfigVersion = v
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
