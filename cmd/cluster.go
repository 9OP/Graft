package cmd

import (
	"fmt"
	"net/netip"

	"graft/pkg"
	"graft/pkg/domain"

	"github.com/spf13/cobra"
)

var cluster ipAddr

var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Manage cluster",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var configurationCmd = &cobra.Command{
	Use:   "configuration",
	Short: "Print cluster configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		clusterPeer := domain.Peer{Host: cluster.AddrPort}
		configuration, err := pkg.ClusterConfiguration(clusterPeer)
		if err != nil {
			return err
		}

		c, err := configuration.ToJSON()
		if err != nil {
			return err
		}

		fmt.Println(string(c))

		return nil
	},
}

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
	clusterCmd.PersistentFlags().Var(&cluster, "cluster", "Live cluster peer for sending commands")
	clusterCmd.Flag("cluster").DefValue = "<nil>"
	clusterCmd.MarkPersistentFlagRequired("cluster")

	clusterCmd.AddCommand(configurationCmd)
	rootCmd.AddCommand(clusterCmd)
	rootCmd.AddCommand(leadershipTransferCmd)
}
