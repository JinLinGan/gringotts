package check

import (
	"github.com/jinlingan/gringotts/gringotts-agent/model"
)

// Loader 用来加载配置生成实际要运行的 Check
type Loader interface {

	// 是否在加载时就确定 Sender
	//Loade(*model.JobConfig, sender scheduler.Sender) (Check, error)

	Loade(*model.JobConfig) (Check, error)

	//TODO: get jobID 是不是需要，这样可以提醒用户需要获取 JoBID 不然无法发送数据。

	//GetJobID() string
}

var Loaders = map[string]Loader{}

func RegisterLoader(name string, l Loader) {
	Loaders[name] = l
}
