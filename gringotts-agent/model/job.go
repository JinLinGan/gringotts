package model

import (
	"github.com/jinlingan/gringotts/pkg/log"
	"github.com/jinlingan/gringotts/pkg/message"
)

type JobState struct {
	JobID           string
	RunningState    int32
	ErrorMsg        string
	LastRunningTime int64
}

type JobConfig struct {
	JobID         string
	RunnerType    string
	RunnerModule  string
	ModuleVersion string
	Interval      int32
	Config        string
}

const (
	RunnerTypeDatadog  = "datadog"
	RunnerTypeTelegraf = "telegraf"
)

const (
	RunningStateOK    = 0
	RunningStateError = 1
)

func GetJobsFromGRPC(jobs []*message.Job) []*JobConfig {
	jobConfigs := make([]*JobConfig, 0, len(jobs))
	for _, v := range jobs {
		jc := &JobConfig{
			JobID:         v.JobID,
			RunnerModule:  v.RunnerModule,
			ModuleVersion: v.ModuleVersion,
			Interval:      v.Interval,
			Config:        v.Config,
		}

		switch v.RunnerType {
		case message.JobRunner_Telegraf:
			jc.RunnerType = RunnerTypeTelegraf
		case message.JobRunner_Datadog:
			jc.RunnerType = RunnerTypeDatadog
		default:
			log.Errorf("收到未知任务类型 %s - %d 请求报文为 %s", v.RunnerType, v.RunnerType, v)
		}

		jobConfigs = append(jobConfigs, jc)

	}
	return jobConfigs
}
