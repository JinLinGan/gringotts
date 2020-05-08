package collector

import (
	"sync"

	"github.com/jinlingan/gringotts/gringotts-agent/check"
	"github.com/jinlingan/gringotts/gringotts-agent/scheduler"
)

// Collector 实际执行 Checker
type Collector struct {
	checkInstances int64
	state          uint32

	scheduler *scheduler.Scheduler
	//runner    *runner.Runner
	checks map[string]check.Check

	m sync.RWMutex
}
