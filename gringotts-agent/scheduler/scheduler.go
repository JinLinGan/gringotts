package scheduler

import (
	"fmt"
	"sync"
	"time"

	"github.com/jinlingan/gringotts/gringotts-agent/global"

	check2 "github.com/jinlingan/gringotts/gringotts-agent/check"

	_ "github.com/jinlingan/gringotts/gringotts-agent/check/loaders"
	"github.com/jinlingan/gringotts/gringotts-agent/model"
	"github.com/jinlingan/gringotts/pkg/log"
)

// Scheduler 加载配置，获取 Check 并调度执行
type Scheduler interface {
	//Schedule([]*model.JobConfig) 增量接口暂不实现
	//UnSchedule([]*model.JobConfig) 增量接口暂不实现
	ReloadConfig([]*model.JobConfig)
	GetJobstate() []*model.JobState
	//GetSender(jobID string) global.Sender
	GetChecksByNameForConfigs()
	Stop()
}

//var GlobalScheduler Scheduler

type JobScheduler struct {
	sync.RWMutex
	logger log.Logger
	Tasks  map[string]*Task
	//Out    chan<- model.Metric
	//SenderPool map[string]Sender
}

//func NewCheckScheduler(out chan<- model.Metric, log log.Logger) Scheduler {
func NewJobScheduler(log log.Logger) Scheduler {
	return &JobScheduler{
		RWMutex: sync.RWMutex{},
		logger:  log,
		Tasks:   make(map[string]*Task),
		//Out:     out,
		//SenderPool: make(map[string]Sender),
	}
}

func (c *JobScheduler) deleteJobs(jobs []*model.JobConfig) {

	jobID := make(map[string]bool, len(jobs))

	// 取出所有的 JobID
	for _, j := range jobs {
		jobID[j.JobID] = true
	}

	for id, t := range c.Tasks {
		if !jobID[id] {
			c.logger.Infof("正在停止任务 %q", id)
			t.Stop()
			c.logger.Infof("停止任务 %q 成功", id)
		}

		// 也可以直接删除不判断
		_, ok := global.GlobalSenderPool.SenderPool[id]
		if ok {
			delete(global.GlobalSenderPool.SenderPool, id)
		}
	}
}
func (c *JobScheduler) addJobs(jobs []*model.JobConfig) {
	for _, j := range jobs {
		if _, ok := c.Tasks[j.JobID]; !ok {
			t := c.loadCheck(j)

			// TODO:是否只需要加正常的
			c.Tasks[j.JobID] = t

			if t.JobState.RunningState == model.RunningStateOK {
				// 运行
				//newSender := NewDefaultSender(c.Out)
				//c.SenderPool[j.JobID] = newSender
				t.Run()
			}
		}
	}
}

func (c *JobScheduler) loadCheck(job *model.JobConfig) *Task {
	loader, ok := check2.Loaders[job.RunnerType]
	if !ok {
		c.logger.Errorf("任务 ID %q 任务类型为 %q 但是未找到与之对应的 Loader", job.JobID, job.RunnerType)
		return createErrorCheckTask(fmt.Sprintf("未找到任务类型 %q 对应的 Loader", job.RunnerType), job)
	}

	check, err := loader.Loade(job)
	if err != nil {
		c.logger.Errorf("任务 ID %q 加载配置时出错：%v", job.JobID, err)
		return createErrorCheckTask(fmt.Sprintf("加载配置时出错：%v", err), job)
	}

	return &Task{
		JobConfig: *job,
		JobState: model.JobState{
			JobID: job.JobID,
			//TODO:细分状态，初始化成功、失败？
			RunningState: model.RunningStateOK,
		},
		Check: check,
	}

}

func createErrorCheckTask(errMsg string, job *model.JobConfig) *Task {
	return &Task{
		JobConfig: *job,
		JobState: model.JobState{
			JobID:           job.JobID,
			RunningState:    model.RunningStateError,
			ErrorMsg:        errMsg,
			LastRunningTime: time.Now().Unix(),
		},
		Check: nil,
	}
}

func (c *JobScheduler) ReloadConfig(jobs []*model.JobConfig) {
	c.deleteJobs(jobs)
	c.addJobs(jobs)
}

func (c *JobScheduler) GetJobstate() []*model.JobState {
	panic("implement me")
}

func (c *JobScheduler) GetChecksByNameForConfigs() {
	panic("implement me")
}

func (c *JobScheduler) Stop() {

	//q.stop <- true
	//<-q.stopped

	//加锁
	//stop all checker
}

func NewCoreScheduler() Scheduler {
	return &JobScheduler{}
}
