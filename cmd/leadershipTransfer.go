package cmd

import (
	"fmt"

	"graft/pkg"

	"github.com/spf13/cobra"
)

var leadershipTransferCmd = &cobra.Command{
	Use:     "leader [ip:port]",
	GroupID: "cluster",
	Short:   "Transfer cluster leadership",
	Args:    validateAddrArg,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := pkg.LeadeshipTransfer(args[0]); err != nil {
			return fmt.Errorf("failed leadership transfer\n%v", err.Error())
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(leadershipTransferCmd)
}
