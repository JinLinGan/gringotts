package scheduler

import (
	"time"

	"github.com/jinlingan/gringotts/gringotts-agent/check"
	"github.com/jinlingan/gringotts/gringotts-agent/model"
)

//TODO:目前直接使用 Job 的结构，后续可能需要定义自己的结构方便使用

// Task 代表一个需要运行的任务,包含任务配置、任务状态、任务实例
type Task struct {
	JobConfig model.JobConfig
	JobState  model.JobState
	Check     check.Check
	stop      chan bool // to stop this queue
	stopped   chan bool // signals that this queue has stopped
}

func (t *Task) Run() {
	go func() {
		tick := time.NewTicker(time.Second * time.Duration(t.JobConfig.Interval))
		defer tick.Stop()

		// 先跑一次
		t.Check.Run()

		for {
			select {
			case <-tick.C:
				if t.Check != nil {
					t.Check.Run()
				}
			case <-t.stop:
				t.stopped <- true
				return
			}
		}
	}()
}

func (t *Task) Stop() {
	//TODO: 需要判断有没有运行，如果 Check 是 nil 外部不会调用 Run ，也无法 Stop
	t.stop <- true
	<-t.stopped
}
