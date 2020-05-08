package forwarder

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/jinlingan/gringotts/gringotts-agent/model"
)

type Forwarder struct {
	in chan model.Metric
}

const forwarderInputBuffer = 1000

func NewForwarder() *Forwarder {
	return &Forwarder{in: make(chan model.Metric, forwarderInputBuffer)}
}

//GetInputChannel 获取 Forwarder 的 Input Channel 用于给 Sender 写入数据
func (f *Forwarder) GetInputChannel() chan<- model.Metric {
	return f.in
}

func (f *Forwarder) Run() {
	for m := range f.in {
		spew.Dump(m)
	}
}
