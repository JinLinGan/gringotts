package check

// Check 代表了没一个监控任务实例

// Check is an interface for types capable to run checks
type Check interface {
	Run() error // run the check
	Stop()      // stop the check if it's running
	//String() string                // provide a printable version of the check name
	//Configure(config string) error // configure the check from the outside
	//Interval() time.Duration                                            // return the interval time for the check
	//ID() ID                                                             // provide a unique identifier for every check instance
	//GetWarnings() []error                                               // return the last warning registered by the check
	//GetMetricStats() (map[string]int64, error) // get metric stats from the sender
	//Version() string                           // return the version of the check if available
	//ConfigSource() string                      // return the configuration source of the check
	//IsTelemetryEnabled() bool                                           // return if telemetry is enabled for this check

	//TODO:考虑要在什么阶段设置 Sender ？ Checker ，Scheduler，loader？
	//目前放到 Scheduler，并由 Scheduler 负责维护一个 SenderPool
	//Datagod 是由具体的 Checker 负责向 Aggregate 要一个 Sender 因为不是所有的 Metrics 都会经过 Checker
	//python 的 Checker 会直接向暴露的函数发送监控数据，所以好像有必要实现一个全局的数据接收器。
	//如果以后需要试跑 Python 怎么通知 Python 监控数据接收方法要额外处理监控数据？
	//Check 内自己做判断获取一个调试用的 Sender，全局的数据接收器也需要知道这个 job 不需要发送 ？ Scheduler 也需要知道它是一个调试用的 Sender。所以 Scheduler 好像要是全局的。

	//// SetSender 设置发送 channel
	//SetSender(chan<-)
}
