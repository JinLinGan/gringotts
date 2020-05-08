package model

type JobManagerService interface {
	CreateJobForAgent(agentID string, job *Job) (jobID string, err error)
	DeleteJobForAgent(agentID string, jobID string) error
	UpdateJobConfigForAgent(agentID string, job *Job) error
}
