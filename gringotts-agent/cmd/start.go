package cmd

import (
	"os"

	"github.com/pkg/errors"

	"github.com/jinlingan/gringotts/gringotts-agent/agent"
	"github.com/jinlingan/gringotts/gringotts-agent/config"
	"github.com/jinlingan/gringotts/pkg/log"
	"github.com/spf13/cobra"
)

// startCmd represents the start command

func newStartCmd() *cobra.Command {
	startCmd := &cobra.Command{
		Use:          "start",
		Short:        "start agent",
		SilenceUsage: true,
		// Long: `start agent`,
		RunE: parseFlagsAndStartAgent,
	}
	startCmd.PersistentFlags().StringP(workPathFlagName, "w", config.GetDefaultWorkPath(),
		"work path used to save all program files")
	startCmd.PersistentFlags().StringP(serverAddressFlagName, "s", config.GetDefaultServerAddress(),
		"server address")
	return startCmd
}

const (
	workPathFlagName      = "workpath"
	serverAddressFlagName = "server"
)

func parseFlagsAndStartAgent(cmd *cobra.Command, args []string) error {

	flags := cmd.Flags()

	// 根据命令行参数设置工作目录
	w, err := flags.GetString(workPathFlagName)
	if err != nil {
		return errors.Wrapf(err, "get flag value of %s", workPathFlagName)
	}

	cfg, err := config.NewConfig(w)
	if err != nil {
		return errors.Wrapf(err, "create agent config fail")
	}

	//初始化 logger
	stdLogger := log.NewStdoutLogger()
	p := cfg.GetWorkPath() + string(os.PathSeparator) + "logs" +
		string(os.PathSeparator) + "gringotts-agent.log"
	stdLogger.Infof("set log file path to %s", p)
	logger := log.NewStdAndFileLogger(p)

	// 使用新的 logger 替换
	log.SetDefaultLogger(logger)
	cfg.SetLogger(logger)

	s, err := flags.GetString(serverAddressFlagName)
	if err != nil {
		logger.Fatal(err, "get flag value of %s", serverAddressFlagName)
	}

	// 根据命令行参数设置服务端地址
	if s != config.GetDefaultServerAddress() {
		logger.Infof("set server address to %s", s)
		cfg.SetServerAddress(s)
	}

	a := agent.NewAgent(cfg, logger)

	return a.Start()
}
