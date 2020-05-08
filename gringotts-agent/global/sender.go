package global

import (
	"github.com/jinlingan/gringotts/gringotts-agent/model"

	"github.com/jinlingan/gringotts/pkg/log"
)

type SenderPool interface {
	GetSender(string) Sender
}

type Sender interface {
	SendTelegraf(model.Metric)
}

var GlobalSenderPool *OutSenderPool

func NewGlobalSenderPool(out chan<- model.Metric, logger log.Logger) *OutSenderPool {
	return &OutSenderPool{
		logger:     logger,
		SenderPool: make(map[string]Sender),
		Out:        out,
	}
}

type OutSenderPool struct {
	logger     log.Logger
	SenderPool map[string]Sender
	Out        chan<- model.Metric
}

func (s *OutSenderPool) GetSender(jobID string) Sender {
	if s, ok := s.SenderPool[jobID]; ok {
		return s
	}

	newSender := NewDefaultSender(s.Out)
	s.SenderPool[jobID] = newSender

	return newSender

}

type DefaultSender struct {
	out chan<- model.Metric
}

func (s *DefaultSender) SendTelegraf(metric model.Metric) {
	s.out <- metric
}

// NewDefaultSender 创建默认发送器
func NewDefaultSender(out chan<- model.Metric) Sender {
	return &DefaultSender{out: out}
}

//TODO: 后续需要增加一次性发送器
