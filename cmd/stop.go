package cmd

import (
	"net/netip"

	"graft/pkg"
	"graft/pkg/domain"

	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop a cluster node",
	Long: `Stop a cluster node:

	- Mark node as inactive
	- Remove node from cluster configuration
	- Shutdown node
	`,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactArgs(1)(cmd, args); err != nil {
			return err
		}

		if _, err := netip.ParseAddrPort(args[0]); err != nil {
			return err
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		host, _ := netip.ParseAddrPort(args[0])
		id := hashString(host.String())
		peer := domain.Peer{Id: id, Host: host}

		return pkg.RemoveClusterPeer(peer)
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
