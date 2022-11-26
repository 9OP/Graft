package cmd

import (
	"fmt"
	"net/netip"

	"graft/pkg"
	"graft/pkg/domain"

	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop [<ip>:<port>]",
	Short: "Shutdown node",
	Args:  argAddrValidator,
	RunE: func(cmd *cobra.Command, args []string) error {
		host, _ := netip.ParseAddrPort(args[0])
		id := hashString(host.String())
		peer := domain.Peer{Id: id, Host: host}

		if err := pkg.Shutdown(peer); err != nil {
			return fmt.Errorf("did not shutdown cluster peer %s: %v", host, err.Error())
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
