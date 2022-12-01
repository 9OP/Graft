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
		Short:   "Raft distributed consensus",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
)

var (
	cluster         ipAddr
	adminGroup      = &cobra.Group{ID: "cluster", Title: "Cluster administration commands:"}
	membershipGroup = &cobra.Group{ID: "membership", Title: "Cluster membership commands:"}
	fsmGroup        = &cobra.Group{ID: "fsm", Title: "State machine commands:"}
)

func Execute() {
	rootCmd.AddGroup(membershipGroup)
	rootCmd.AddGroup(fsmGroup)
	rootCmd.AddGroup(adminGroup)

	rootCmd.PersistentFlags().VarP(&cluster, "cluster", "c", "Any live cluster peer")
	rootCmd.MarkPersistentFlagRequired("cluster")
	rootCmd.PersistentFlags().SortFlags = false
	rootCmd.SilenceUsage = true

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
