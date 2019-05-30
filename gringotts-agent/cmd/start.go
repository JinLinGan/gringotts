package cmd

import (
	"log"

	"github.com/jinlingan/gringotts/gringotts-agent/agent"
	"github.com/jinlingan/gringotts/gringotts-agent/config"

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
		log.Fatalf("get flag value of %s error: %s", workPathFlagName, err)
	}

	cfg, err := config.NewConfig(w)
	if err != nil {
		log.Printf("create agent with err: %s", err)
	}

	s, err := flags.GetString(serverAddressFlagName)
	if err != nil {
		log.Fatalf("get flag value of %s error: %s", serverAddressFlagName, err)
	}

	// 根据命令行参数设置服务端地址
	if s != "" {
		log.Printf("set server address to %s", s)
		cfg.SetServerAddress(s)
	}

	a := agent.NewAgent(cfg)

	return a.Start()
}
