package cmd

import (
	"log"

	"github.com/jinlingan/gringotts/gringotts-agent/agent"
	"github.com/jinlingan/gringotts/gringotts-agent/config"

	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start agent",
	// Long: `start agent`,
	RunE: praseFlagsAndStartAgent,
}

func initStartCmd() {
	rootCmd.AddCommand(startCmd)
	startCmd.PersistentFlags().StringP("workpath", "w", "",
		"work path used to save all program files (default \""+config.GetDefaultWorkPath()+"\")")
	startCmd.PersistentFlags().StringP("server", "s", "",
		"server address  (default \""+config.GetDefaultServerAddress()+"\")")
}

const (
	workPathFlagName      = "workpath"
	serverAddressFlagName = "server"
)

func praseFlagsAndStartAgent(cmd *cobra.Command, args []string) error {

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
