//Package cmd include all agent command
package cmd

import (
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
	rootCmd.AddCommand(newStartCmd(), newJobCmd())
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
