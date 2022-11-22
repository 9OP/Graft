package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	version = "0.0.1"
	rootCmd = &cobra.Command{
		Use:     "graft",
		Version: version,
		Short:   "Graft is a simple implementation of Raft distributed consensus",
		Long:    `Makes any FSM distributed and resilient to single point of failure.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
)

func Execute() {
	rootCmd.Flags().SortFlags = false
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
