package model

import (
	"github.com/jinlingan/gringotts/pkg/message"
	"github.com/pkg/errors"
)

type Job struct {
	JobID           string `protobuf:"bytes,1,opt,name=jobID,proto3" json:"ID,omitempty"`
	AgentID         string
	RunnerType      JobRunner `protobuf:"varint,2,opt,name=runnerType,proto3,enum=JobRunner" json:"runnerType,omitempty"`
	RunnerModule    string    `protobuf:"bytes,3,opt,name=runnerModule,proto3" json:"runnerModule,omitempty"`
	ModuleVersion   string    `protobuf:"bytes,4,opt,name=moduleVersion,proto3" json:"moduleVersion,omitempty"`
	RunningInterval int32     `protobuf:"varint,5,opt,name=interval,proto3" json:"interval,omitempty"`
	Config          string    `protobuf:"bytes,6,opt,name=config,proto3" json:"config,omitempty"`
	CreateTime      int32     `protobuf:"varint,7,opt,name=createTime,proto3" json:"createTime,omitempty"`
	UpdateTime      int32     `protobuf:"varint,8,opt,name=updateTime,proto3" json:"updateTime,omitempty"`

	RunningState    JobRunningState `protobuf:"varint,2,opt,name=state,proto3,enum=JobRunningState" json:"state,omitempty"`
	ErrorMsg        string          `protobuf:"bytes,3,opt,name=errorMsg,proto3" json:"errorMsg,omitempty"`
	LastRunningTime int64           `protobuf:"varint,4,opt,name=lastRunningTime,proto3" json:"lastRunningTime,omitempty"`
	LastReportTime  int32
}

type JobStore interface {
	CreateJob(job *Job) (jobID string, err error)
	DeleteJob(jobID string) error
	UpdateJobConfig(job *Job) error

	GetJobs(agentID string) ([]*Job, error)
	UpdateJobRunningState(job *Job) error
}

type JobRunner int32

const (
	JobRunnerTelegraf JobRunner = 0
	JobRunnerDatadog  JobRunner = 1
)

type JobRunningState int32

const (
	JobRunningStateWait  JobRunningState = 0
	JobRunningStateOk    JobRunningState = 1
	JobRunningStateError JobRunningState = 2
	JobRunningStateUndef JobRunningState = 3
)

func NewJobRunningInfoFromGRPC(req *message.HeartBeatRequest) []*Job {
	jobsRunningInfo := make([]*Job, 0, len(req.Jobs))

	for _, j := range req.Jobs {
		n := &Job{
			JobID:           j.JobID,
			ErrorMsg:        j.ErrorMsg,
			LastRunningTime: j.LastRunningTime,
		}

		switch j.State {
		case message.JobRunningState_OK:
			n.RunningState = JobRunningStateOk
		case message.JobRunningState_Error:
			n.RunningState = JobRunningStateError
		default:
			n.RunningState = JobRunningStateUndef
		}
		jobsRunningInfo = append(jobsRunningInfo, n)
	}

	return jobsRunningInfo
}

func NewGRPCJobs(jobs []*Job) []*message.Job {
	grpcJobs := make([]*message.Job, 0, len(jobs))

	for _, j := range jobs {
		n := &message.Job{
			JobID:         j.JobID,
			RunnerModule:  j.RunnerModule,
			ModuleVersion: j.ModuleVersion,
			Interval:      j.RunningInterval,
			Config:        j.Config,
			CreateTime:    j.CreateTime,
			UpdateTime:    j.UpdateTime,
		}

		switch j.RunnerType {
		case JobRunnerTelegraf:
			n.RunnerType = message.JobRunner_Telegraf
		case JobRunnerDatadog:
			n.RunnerType = message.JobRunner_Datadog
		default:
			continue
		}
		grpcJobs = append(grpcJobs, n)
	}
	return grpcJobs
}

func NewJobFromGRPC(req *message.AddJobRequest) (*Job, error) {
	j := &Job{
		AgentID:         req.AgentID,
		RunnerModule:    req.RunnerModule,
		ModuleVersion:   req.ModuleVersion,
		RunningInterval: req.Interval,
		Config:          req.Config,
	}
	switch req.RunnerType {
	case message.JobRunner_Telegraf:
		j.RunnerType = JobRunnerTelegraf
	case message.JobRunner_Datadog:
		j.RunnerType = JobRunnerDatadog
	default:
		return nil, errors.Errorf("未知的执行器类型 %v", req.RunnerType)

	}

	return j, nil

}
