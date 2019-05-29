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
	"fmt"

	"github.com/jinlingan/gringotts/gringotts-server/server"
	"github.com/spf13/cobra"
)

func newStartCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "start Gringotts Server",
		Long:  `start Gringotts Server`,
		RunE:  start,
	}
}

const (
	serverID = "99"
)

func start(cmd *cobra.Command, args []string) error {
	address := ":7777"
	serverInst, err := server.NewServer(address, serverID)
	if err != nil {
		return fmt.Errorf("can not create new server in port %s : %s", address, err)
	}
	return serverInst.Serve()

}
