package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"time"

	"github.com/jinlingan/gringotts/gringotts-server/service/jobmanager"
	"github.com/jinlingan/gringotts/gringotts-server/store/job"

	"github.com/jinlingan/gringotts/gringotts-server/store/host"

	"github.com/jinlingan/gringotts/gringotts-server/model"

	"github.com/jinlingan/gringotts/gringotts-server/config"
	"github.com/jinlingan/gringotts/pkg/log"
	"github.com/jinlingan/gringotts/pkg/message"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// GringottsServer 服务器
type GringottsServer struct {
	sync.RWMutex
	grServer   *grpc.Server
	config     *config.ServerConfig
	logger     log.Logger
	db         *sqlx.DB
	hostStore  model.HostStore
	jobStore   model.JobStore
	jobManager model.JobManagerService
}

func (s *GringottsServer) AddJob(ctx context.Context, req *message.AddJobRequest) (*message.AddJobResponse, error) {
	agent, err := s.hostStore.Find(req.AgentID)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("查找 Agent（%q）失败", req.AgentID))
	}
	if agent == nil {
		return nil, errors.Errorf("未找到 Agent（%q）", req.AgentID)
	}

	job, err := model.NewJobFromGRPC(req)
	if err != nil {
		return nil, err
	}

	jobID, err := s.jobManager.CreateJobForAgent(req.AgentID, job)
	if err != nil {
		return nil, err
	}

	return &message.AddJobResponse{
		JobID: jobID,
	}, nil

}

func (s *GringottsServer) DelJob(ctx context.Context, req *message.DelJobRequest) (*message.DelJobResponse, error) {

	agent, err := s.hostStore.Find(req.AgentID)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("查找 Agent（%q）失败", req.AgentID))
	}
	if agent == nil {
		return nil, errors.Wrap(err, fmt.Sprintf("未找到 Agent（%q）", req.AgentID))
	}

	//TODO:判断任务是否存在，是否属于对应 agent
	err = s.jobManager.DeleteJobForAgent(req.AgentID, req.JobID)
	if err != nil {
		return nil, err
	}
	return &message.DelJobResponse{
		Deleted: true,
	}, nil
}

func (s *GringottsServer) GetJobs(ctx context.Context, req *message.GetJobsRequest) (*message.GetJobsResponse, error) {
	agent, err := s.hostStore.Find(req.AgentID)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("查找 Agent（%q）失败", req.AgentID))
	}
	if agent == nil {
		return nil, errors.Wrap(err, fmt.Sprintf("未找到 Agent（%q）", req.AgentID))
	}

	jobs, err := s.jobStore.GetJobs(req.AgentID)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("查找 Agent（%q）任务失败", req.AgentID))
	}
	return NewGetJobsGRPCResp(jobs, agent.ConfigVersion), nil
}

func NewGetJobsGRPCResp(jobs []*model.Job, configVersion int64) *message.GetJobsResponse {
	return &message.GetJobsResponse{
		ConfigVersion: configVersion,
		Jobs:          model.NewGRPCJobs(jobs),
	}
}

//NewServer 新建 Server 对象
func NewServer(cfg *config.ServerConfig, logger log.Logger) (*GringottsServer, error) {
	//TODO:移动到配置文件中

	dataSourceName := "gringotts:gringotts@tcp(mysql)/gringotts?parseTime=true"
	db, err := sqlx.Connect("mysql", dataSourceName)
	if err != nil {
		return nil, errors.Wrap(err, "open database error")
	}
	h := host.New(db)
	j := job.New(db)
	jm := jobmanager.New(j, h)
	server := &GringottsServer{
		grServer:   grpc.NewServer(),
		config:     cfg,
		logger:     logger,
		db:         db,
		hostStore:  h,
		jobStore:   j,
		jobManager: jm,
	}
	message.RegisterGringottsServer(server.grServer, server)
	return server, nil

}

//Serve 开始提供服务
func (s *GringottsServer) Serve() error {
	lsP := s.config.GetListenerPort()
	lis, err := net.Listen("tcp", lsP)
	if err != nil {
		return errors.Wrapf(err, "can not listen in port 0.0.0.0%s", lsP)
	}
	s.logger.Infof("gringotts server listen in port 0.0.0.0%s", lsP)
	return s.grServer.Serve(lis)
}

//HeartBeat 接收心跳
func (s *GringottsServer) HeartBeat(ctx context.Context,
	req *message.HeartBeatRequest) (*message.HeartBeatResponse, error) {
	agentID := req.GetAgentID()
	if agentID == "" {
		return nil, errors.Errorf("agent ID 为 %q 不合法", agentID)
	}
	s.logger.Debugf("get HeartBeat message from agent(id=%s)", req.GetAgentID())

	ok, err := s.hostStore.Exist(agentID)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("查询 agent ID %q 是否存在失败 ", agentID))
	}

	if !ok {
		return nil, errors.Errorf("agent ID %q 不存在", agentID)
	}
	err = s.hostStore.UpdateHeartBeatTime(agentID, time.Now().Unix())

	configVersion, err := s.hostStore.GetConfigVersion(agentID)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("获取 Agent（%q) 配置版本号失败", agentID))
	}
	return &message.HeartBeatResponse{
		ServerId:      s.config.GetExternalAddress(),
		ConfigVersion: configVersion,
		//MonitorInfo:   getAllTaskByAgentID(),
	}, nil
}

//DownloadFile 下载文件
func (s *GringottsServer) DownloadFile(f *message.File, fs message.Gringotts_DownloadFileServer) error {
	//TODO:改变文件路径
	rf, err := os.Open("/Users/jinlin/code/golang/src/github.com/jinlingan/gringotts/testfile/" + f.GetFileName())

	if err != nil {
		return err
	}
	buff := make([]byte, 500)
	for {
		n, err := rf.Read(buff)
		if err != nil && err != io.EOF {
			return status.Errorf(500, "read file error %s", err)
		}
		exit := false
		if err == io.EOF {
			exit = true
		}

		err = fs.Send(&message.FileChunk{
			Data: buff[:n],
		})
		if err != nil {
			return status.Errorf(500, "send file Chunk error %s", err)
		}

		if exit {
			break
		}
	}
	if rf.Close() != nil {
		return err
	}
	return nil
}
func (s *GringottsServer) getHostByID(id string) (*[]model.Host, error) {
	var hs []model.Host
	err := s.db.Get(&hs, "SELECT * FROM agent WHERE agent_id=?", 1)

	return &hs, err
}

func (s *GringottsServer) addNewAgent(h *model.Host) (string, error) {
	return s.hostStore.Create(h)
}

// Register agent 注册，说明文档 https://docs.google.com/drawings/d/1jFwqKoWa-JNRlh52ZKseSTheOajAmV2WIoVD6dV5p9Q/edit?usp=sharing
func (s *GringottsServer) Register(ctx context.Context,
	req *message.RegisterRequest) (*message.RegisterResponse, error) {

	reqJSON, _ := json.Marshal(req)
	h := model.NewHostFromGRPC(req)

	s.logger.Debugf("收到注册消息：%s", reqJSON)

	agent, err := s.hostStore.Find(h.AgentID)
	if err != nil {
		s.logger.Errorf("使用 AgentID %q 查找已注册 Agent 失败，注册信息为 %s : %s", req.AgentID, reqJSON, err)

		return nil, errors.Wrap(err, "注册 Agent 失败")
	}

	if agent == nil {
		s.logger.Infof("AgentID %q 未找到已注册 Agent", req.AgentID)
	}
	// 如果AgentID 不存在
	if req.AgentID == "" || agent == nil {
		newAgentID, err := s.addNewAgent(h)
		if err != nil {
			s.logger.Errorf("注册新 Agent 失败，注册信息为 %s : %s", reqJSON, err)
			return nil, errors.Wrap(err, "注册 Agent 失败")
		}

		return &message.RegisterResponse{
			AgentId:       newAgentID,
			ConfigVersion: 0,
		}, nil
	}

	allSame, acceptable, msg := model.CheckHostChanceAcceptable(agent, h)
	if allSame {
		s.logger.Debug("agent 信息没有任何变更")
		v, err := s.hostStore.GetConfigVersion(h.AgentID)

		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("获取 agent(%q) 版本错误", h.AgentID))
		}
		return &message.RegisterResponse{
			AgentId:       h.AgentID,
			ConfigVersion: v,
		}, nil
	}

	if !acceptable {

		//TODO:对于未注册成功的 Agent 可以记录信息，支持在线手动迁移，类似有一个半成功的状态 Agent 可以一直尝试注册
		s.logger.Errorf("Agent %q 注册失败，Agent 提交的信息变更太多：%s", req.AgentID, msg)

		return nil, errors.Errorf("Agent %q 注册失败，Agent 提交的信息变更太多：%s", req.AgentID, msg)
	}

	// 更新信息

	s.logger.Debug("agent 信息有变更但是可接收，记录变更信息")

	err = s.hostStore.UpdateAgentInfo(h)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("发现agent(%q)配置更新并且可以接受，但是更新配置时发生错误", h.AgentID))
	}

	v, err := s.hostStore.GetConfigVersion(h.AgentID)

	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("获取 agent(%q) 版本错误", h.AgentID))
	}

	return &message.RegisterResponse{
		AgentId:       h.AgentID,
		ConfigVersion: v,
	}, nil

}

//func (s *GringottsServer) newHeartBeatResponse() *message.HeartBeatResponse {
//	resp := &message.HeartBeatResponse{
//		ServerId:      s.config.GetExternalAddress(),
//		ConfigVersion: ,
//		//MonitorInfo:   getAllTaskByAgentID(),
//	}
//	return resp
//
//}
