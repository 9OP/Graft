package cmd

import (
	"fmt"
	"net/netip"

	"graft/pkg"

	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop [ip:port]",
	Short: "Shutdown node",
	Args:  validateAddrArg,
	RunE: func(cmd *cobra.Command, args []string) error {
		host, _ := netip.ParseAddrPort(args[0])
		if err := pkg.Shutdown(host.String()); err != nil {
			return fmt.Errorf("failed shutdown\n%v", err.Error())
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
