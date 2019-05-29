//Package cmd include all agent command
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Execute 执行
func Execute() {
	var rootCmd = &cobra.Command{
		Use:   "gringotts-agent",
		Short: "Gringotts Agent",
		Long:  `Gringotts Agent`,
	}
	rootCmd.AddCommand(newStartCmd())
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
