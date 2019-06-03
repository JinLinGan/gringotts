// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/pkg/errors"

	"github.com/jinlingan/gringotts/common/log"
	"github.com/jinlingan/gringotts/gringotts-server/config"
	"github.com/jinlingan/gringotts/gringotts-server/server"
	"github.com/spf13/cobra"
)

func newStartCmd() *cobra.Command {
	startCmd := &cobra.Command{
		Use:   "start",
		Short: "start Gringotts Server",
		Long:  `start Gringotts Server`,
		RunE:  start,
	}
	startCmd.PersistentFlags().StringP(listenerPortFlagName, "p", config.GetDefaultListenerPort(), "listener port")
	startCmd.PersistentFlags().StringP(
		externalAddressFlagName,
		"a", config.GetDefaultExternalAddress(), "external address")
	return startCmd
}

const (
	listenerPortFlagName    = "port"
	externalAddressFlagName = "address"
)

func start(cmd *cobra.Command, args []string) error {
	cfg := config.NewServerConfig()
	flags := cmd.Flags()
	//TODO 改变日志路径
	logger := log.NewStdAndFileLogger("/tmp/test.log")
	// 根据命令行参数设置监听端口
	p, err := flags.GetString(listenerPortFlagName)
	if err != nil {
		logger.Fatal(errors.Wrapf(err, "get flag value of %s fail", listenerPortFlagName))
	}

	cfg.SetListenerPort(p)

	a, err := flags.GetString(externalAddressFlagName)
	if err != nil {
		logger.Fatal(errors.Wrapf(err, "get flag value of %s fail", externalAddressFlagName))
	}
	cfg.SetExternalAddress(a)
	serverInst, err := server.NewServer(cfg, logger)

	if err != nil {
		return errors.Wrap(err, "can not create new server in port")
	}
	return serverInst.Serve()

}
