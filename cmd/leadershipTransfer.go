package cmd

import (
	"fmt"
	"net/netip"

	"graft/pkg"
	"graft/pkg/domain"

	"github.com/spf13/cobra"
)

var leadershipTransferCmd = &cobra.Command{
	Use:   "leader [ip:port]",
	Short: "Leadership transfer to peer",
	Args:  argAddrValidator,
	RunE: func(cmd *cobra.Command, args []string) error {
		host, _ := netip.ParseAddrPort(args[0])
		id := hashString(host.String())
		peer := domain.Peer{Id: id, Host: host}

		if err := pkg.LeadeshipTransfer(peer); err != nil {
			return fmt.Errorf("did not transfer leadership to peer %s: %v", host, err.Error())
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(leadershipTransferCmd)
}
