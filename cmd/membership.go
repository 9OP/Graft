package cmd

import (
	"fmt"
	"net/netip"

	"graft/pkg"
	"graft/pkg/domain"
	"graft/pkg/utils"

	"github.com/spf13/cobra"
)

func argAddrValidator(cmd *cobra.Command, args []string) error {
	if err := cobra.ExactArgs(1)(cmd, args); err != nil {
		return err
	}

	if _, err := netip.ParseAddrPort(args[0]); err != nil {
		return err
	}

	return nil
}

var addNodeCmd = &cobra.Command{
	Use:   "add [<ip>:<port>]",
	Short: "Add new node to live cluster",
	Long: `Add new node to live cluster:

1. Fetch cluster configuration from leader
2. Update configuration with new node as inactive
3. Start new node
4. Update configuration with new node as active
	`,
	Args: argAddrValidator,
	RunE: func(cmd *cobra.Command, args []string) error {
		utils.ConfigureLogger(level.String())

		host, _ := netip.ParseAddrPort(args[0])
		id := hashString(host.String())
		newPeer := domain.Peer{Id: id, Host: host, Active: false}
		clusterPeer := domain.Peer{Host: cluster.AddrPort}

		quit, err := pkg.AddClusterPeer(newPeer, clusterPeer)
		if err != nil {
			return err
		}

		// wait
		<-quit
		return nil
	},
}

var shutdown bool

var removeNodeCmd = &cobra.Command{
	Use:   "remove [<ip>:<port>]",
	Short: "Remove node from live cluster",
	Long: `Remove node from live cluster:

1. Fetch cluster leader
2. Update configuration with old node as inactive
3. Update configuration with old node removed
4. Shutdown old node
	`,
	Args: argAddrValidator,
	RunE: func(cmd *cobra.Command, args []string) error {
		host, _ := netip.ParseAddrPort(args[0])
		id := hashString(host.String())
		peer := domain.Peer{Id: id, Host: host}
		clusterPeer := domain.Peer{Host: cluster.AddrPort}

		if err := pkg.RemoveClusterPeer(peer, clusterPeer); err != nil {
			return fmt.Errorf("did not remove cluster peer %s: %v", host, err.Error())
		}

		fmt.Printf("Node %v removed from cluster\n", peer.Host)

		if shutdown {
			if err := pkg.Shutdown(peer); err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	addNodeCmd.Flags().Var(&level, "log", `log level. allowed: "DEBUG", "INFO", "ERROR"`)
	removeNodeCmd.Flags().BoolVarP(&shutdown, "shutdown", "s", false, "Shutdown node once removed")

	clusterCmd.AddCommand(addNodeCmd)
	clusterCmd.AddCommand(removeNodeCmd)
}
