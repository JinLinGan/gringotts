package cmd

import (
	"context"
	"log"

	"github.com/jinlingan/gringotts/pkg/message"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func newJobCmd() *cobra.Command {
	JobCmd := &cobra.Command{
		Use:   "job",
		Short: "修改监控任务",
		Long:  `修改监控任务`,
		RunE:  job,
	}
	pf := JobCmd.PersistentFlags()

	pf.StringP(serverAddressFlag, "s", "", "服务器地址")

	pf.StringP(agentIDFlag, "a", "", "agent id")

	pf.Int32P(runningIntervalFlag, "i", 10, "运行间隔")

	pf.StringP(runningConfigFlag, "c", "", "运行配置")

	pf.StringP(runnerTypeFlag, "r", "telegraf", "执行器类型：telegraf、datadog")

	pf.StringP(moduleVersionFlag, "v", "", "执行器版本")

	pf.StringP(runnerModuleFlag, "m", "", "执行器名称")

	pf.Bool(listJobsFlag, false, "打印 Job")
	pf.Bool(removeFlag, false, "删除 Job")
	pf.String(jobIDFlag, "", "Job ID")

	return JobCmd
}

const (
	agentIDFlag         = "agent"
	runningIntervalFlag = "interval"
	runningConfigFlag   = "config"
	runnerTypeFlag      = "runner"
	moduleVersionFlag   = "version"
	runnerModuleFlag    = "module"
	serverAddressFlag   = "server"

	listJobsFlag = "list"

	removeFlag = "rm"
	jobIDFlag  = "job"
)

func job(cmd *cobra.Command, args []string) error {

	flags := cmd.Flags()

	serverAddress, err := flags.GetString(serverAddressFlag)
	if err != nil {
		return errors.Wrap(err, "获取服务器地址失败")
	}

	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	if err != nil {
		return errors.Wrap(err, "初始化 GRPC 客户端失败")
	}

	server := message.NewGringottsClient(conn)

	remove, err := flags.GetBool(removeFlag)
	if err != nil {
		return errors.Wrap(err, "removeFlag 获取失败")
	}

	list, err := flags.GetBool(listJobsFlag)
	if err != nil {
		return errors.Wrap(err, "listJobsFlag 获取失败")
	}

	if list {
		return listJob(cmd, server)
	}

	if remove {
		return removeJob(cmd, server)
	}

	return AddJob(cmd, server)

	return nil

}

func listJob(cmd *cobra.Command, server message.GringottsClient) error {
	return nil
}

func removeJob(cmd *cobra.Command, server message.GringottsClient) error {
	return nil
}

func AddJob(cmd *cobra.Command, server message.GringottsClient) error {
	flags := cmd.Flags()

	agentID, err := flags.GetString(agentIDFlag)
	if err != nil {
		return errors.Wrap(err, "agent ID 获取失败")
	}
	if agentID == "" {
		return errors.Errorf("添加任务时 agent ID 必填")
	}

	interval, err := flags.GetInt32(runningIntervalFlag)
	if err != nil {
		return errors.Wrap(err, "运行间隔获取失败")
	}

	config, err := flags.GetString(runningConfigFlag)
	if err != nil {
		return errors.Wrap(err, "配置文件获取失败")
	}

	runnerTypeArg, err := flags.GetString(runnerTypeFlag)
	if err != nil {
		return errors.Wrap(err, "执行器类型获取失败")
	}
	var runnerType message.JobRunner
	switch runnerTypeArg {
	case "telegraf":
		runnerType = message.JobRunner_Telegraf
	case "datadog":
		runnerType = message.JobRunner_Datadog
	default:
		return errors.Errorf("未知执行器类型 %q", runnerTypeArg)
	}

	version, err := flags.GetString(moduleVersionFlag)
	if err != nil {
		return errors.Wrap(err, "模块版本获取失败")
	}

	moduleName, err := flags.GetString(runnerModuleFlag)
	if err != nil {
		return errors.Wrap(err, "获取模块名称失败")
	}

	req := &message.AddJobRequest{
		AgentID:       agentID,
		RunnerType:    runnerType,
		RunnerModule:  moduleName,
		ModuleVersion: version,
		Interval:      interval,
		Config:        config,
	}

	resp, err := server.AddJob(context.Background(), req)
	if err != nil {
		return errors.Wrap(err, "添加任务失败")
	}
	log.Printf("添加任务成功，任务 ID 为 %q", resp.JobID)

	return nil
}
