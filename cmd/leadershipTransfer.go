package cmd

import (
	"fmt"

	"graft/pkg"

	"github.com/spf13/cobra"
)

var leadershipTransferCmd = &cobra.Command{
	Use:   "leader [ip:port]",
	Short: "Leadership transfer to peer",
	Args:  validateAddrArg,
	RunE: func(cmd *cobra.Command, args []string) error {
		host := args[0]

		if err := pkg.LeadeshipTransfer(host); err != nil {
			return fmt.Errorf("failed leadership transfer\n%v", err.Error())
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(leadershipTransferCmd)
}
