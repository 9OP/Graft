package cmd

import (
	"fmt"
	"net/netip"

	"graft/pkg"
	"graft/pkg/domain"

	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:     "stop [ip:port]",
	GroupID: "membership",
	Short:   "Stop cluster node",
	Args:    validateAddrArg,
	RunE: func(cmd *cobra.Command, args []string) error {
		host, _ := netip.ParseAddrPort(args[0])
		id := hashString(host.String())
		peer := domain.Peer{Id: id, Host: host}

		if err := pkg.RemoveClusterPeer(peer, cluster.String()); err != nil {
			return fmt.Errorf("failed remove peer\n%v", err.Error())
		}

		fmt.Print("removed from cluster ", peer.Host)

		if err := pkg.Shutdown(peer.Host.String()); err != nil {
			return fmt.Errorf("shutdown failed\n%v", err.Error())
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
