package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var version = "0.0.1"
var rootCmd = &cobra.Command{
	Use:     "graft",
	Version: version,
	Short:   "Graft is a simple implementation of Raft distributed consensus",
	Long:    `Makes any FSM distributed and resilient to single point of failure.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		// fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
