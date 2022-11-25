package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var cluster ipAddr

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

/*
Commands; todo
- start / add node
- shutdown node
- print config
*/

func Execute() {
	rootCmd.PersistentFlags().Var(&cluster, "cluster", "A live cluster peer addr for issuing commands")
	rootCmd.MarkPersistentFlagRequired("cluster")
	rootCmd.Flag("cluster").DefValue = "<nil>"

	rootCmd.Flags().SortFlags = false
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
