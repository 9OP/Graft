package cmd

import (
	"fmt"
	"net/netip"

	"graft/pkg"
	"graft/pkg/domain"
	"graft/pkg/utils/log"

	"github.com/spf13/cobra"
)

var addNodeCmd = &cobra.Command{
	Use:   "add [ip:port]",
	Short: "Add new node to live cluster",
	Long: `Add new node to live cluster:

1. Fetch cluster configuration from leader
2. Update configuration with new node as inactive
3. Start new node
4. Update configuration with new node as active
	`,
	Args: validateAddrArg,
	RunE: func(cmd *cobra.Command, args []string) error {
		host, _ := netip.ParseAddrPort(args[0])
		id := hashString(host.String())
		log.ConfigureLogger(id)
		newPeer := domain.Peer{Id: id, Host: host, Active: false}

		quit, err := pkg.AddClusterPeer(newPeer, cluster.String())
		if err != nil {
			return fmt.Errorf("failed add peer\n%v", err.Error())
		}

		// wait
		<-quit
		return nil
	},
}

var shutdown bool

var removeNodeCmd = &cobra.Command{
	Use:   "remove [ip:port]",
	Short: "Remove node from live cluster",
	Long: `Remove node from live cluster:

1. Fetch cluster leader
2. Update configuration with old node as inactive
3. Update configuration with old node removed
4. Shutdown old node
	`,
	Args: validateAddrArg,
	RunE: func(cmd *cobra.Command, args []string) error {
		host, _ := netip.ParseAddrPort(args[0])
		id := hashString(host.String())
		peer := domain.Peer{Id: id, Host: host}

		if err := pkg.RemoveClusterPeer(peer, cluster.String()); err != nil {
			return fmt.Errorf("failed remove peer\n%v", err.Error())
		}

		fmt.Printf("Node %v removed from cluster\n", peer.Host)

		if shutdown {
			if err := pkg.Shutdown(peer.Host.String()); err != nil {
				return fmt.Errorf("failed shutdown peer\n%v", err.Error())
			}
		}

		return nil
	},
}

func init() {
	removeNodeCmd.Flags().BoolVarP(&shutdown, "shutdown", "s", false, "Shutdown node once removed")
	clusterCmd.AddCommand(addNodeCmd)
	clusterCmd.AddCommand(removeNodeCmd)
}
