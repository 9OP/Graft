package cmd

import (
	"github.com/spf13/cobra"
)

var cluster ipAddr

var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Manage cluster",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(clusterCmd)
}
