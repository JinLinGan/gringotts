package jobmanager

import (
	"fmt"

	"github.com/jinlingan/gringotts/gringotts-server/model"
	"github.com/pkg/errors"
)

type jobManagerService struct {
	jobStore  model.JobStore
	hostStore model.HostStore
}

func (j jobManagerService) CreateJobForAgent(agentID string, job *model.Job) (jobID string, err error) {

	jobID, err = j.jobStore.CreateJob(job)
	if err != nil {
		return "", errors.Wrap(err, "添加任务失败")
	}

	err = j.hostStore.UpdateConfigVersion(agentID)
	if err != nil {
		dErr := j.jobStore.DeleteJob(jobID)
		if dErr != nil {
			return "", errors.Wrap(err, fmt.Sprintf("更新 Agent (%q) 配置版本号失败，尝试删除刚刚新增的任务 (%q) 也失败了", agentID, jobID))
		}
		return "nil", errors.Wrap(err, fmt.Sprintf("更新 Agent (%q) 配置版本号失败，已经删除刚刚新增的任务 (%q) ", agentID, jobID))
	}
	return jobID, nil
}

func (j jobManagerService) DeleteJobForAgent(agentID string, jobID string) error {

	err := j.jobStore.DeleteJob(jobID)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("删除任务 (%q) 失败", jobID))
	}

	err = j.hostStore.UpdateConfigVersion(agentID)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("删除任务 (%q) 成功，但更新 Agent (%q) 配置版本号失败，当前操作不会生效直到 agent 配置版本号被更新", jobID, agentID))
	}
	return nil
}

func (j jobManagerService) UpdateJobConfigForAgent(agentID string, job *model.Job) error {

	err := j.jobStore.UpdateJobConfig(job)
	if err != nil {
		return errors.Wrap(err, "更新任务失败")
	}

	err = j.hostStore.UpdateConfigVersion(agentID)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("更新任务成功，但更新 Agent (%q) 配置版本号失败，当前操作不会生效直到 agent 配置版本号被更新", agentID))
	}
	return nil
}

func New(jobStore model.JobStore, hostStore model.HostStore) model.JobManagerService {
	return &jobManagerService{
		jobStore:  jobStore,
		hostStore: hostStore,
	}
}
