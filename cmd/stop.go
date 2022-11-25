package cmd

import (
	"fmt"
	"net/netip"

	"graft/pkg"
	"graft/pkg/domain"

	"github.com/spf13/cobra"
)

var removeFromCluster bool

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop a cluster node",
	Long: `Stop a cluster node:

	This command allows to:
	- Shutdown a cluster node
	- Shutdown and remove node from cluster configuration
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

		if removeFromCluster {
			if err := pkg.RemoveClusterPeer(peer); err != nil {
				return fmt.Errorf("did not remove cluster peer %s: %v", host, err.Error())
			}
		} else {
			if err := pkg.Shutdown(peer); err != nil {
				return fmt.Errorf("did not shutdown cluster peer %s: %v", host, err.Error())
			}
		}

		return nil
	},
}

func init() {
	stopCmd.Flags().BoolVar(&removeFromCluster, "remove", false, "Remove node from cluster configuration")
	stopCmd.MarkPersistentFlagRequired("cluster")

	rootCmd.AddCommand(stopCmd)
}
