package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var membershipCmd = &cobra.Command{
	Use:   "membership",
	Short: "Manage cluster memberships",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var addMemberCmd = &cobra.Command{
	Use:   "add [peer]",
	Short: "Add new node to cluster",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("add node")
	},
}

var removeMemberCmd = &cobra.Command{
	Use:   "remove [peer-id]",
	Short: "Remove old node from cluster",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("remove node")
	},
}

func init() {
	// startCmd.Flags().VarP(&level, "log-level", "l", `log level. allowed: "DEBUG", "INFO", "ERROR"`)
	// startCmd.Flags().StringVarP(&configPath, "config", "c", "conf/graft-config.yml", "Configuration file path")
	membershipCmd.AddCommand(addMemberCmd)
	membershipCmd.AddCommand(removeMemberCmd)
	rootCmd.AddCommand(membershipCmd)
}
