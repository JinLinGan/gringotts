//Package cmd include all agent command
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gringotts-agent",
	Short: "Gringotts Agent",
	Long:  `Gringotts Agent`,
}

// Execute 执行
func Execute() {
	initStartCmd()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
