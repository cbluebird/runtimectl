package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "runtimectl",
	Short: "runtimectl is a tool for managing runtimes for devbox",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		if rootCmd.SilenceErrors {
			fmt.Println(err)
		}
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(newSyncCmd())
	rootCmd.AddCommand(newGenCmd())
}
