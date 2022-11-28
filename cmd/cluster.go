package cmd

import (
	"fmt"
	"net/netip"

	"graft/pkg"
	"graft/pkg/domain"

	"github.com/spf13/cobra"
)

var (
	cluster ipAddr
	exType  execType
)

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

var executeCmd = &cobra.Command{
	Use:   "execute [command]",
	Short: "Execute command on FSM",
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactArgs(1)(cmd, args); err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		clusterPeer := domain.Peer{Host: cluster.AddrPort}
		entry := args[0]
		var logType domain.LogType

		switch exType {
		case "COMMAND":
			logType = domain.LogCommand
		case "QUERY":
			logType = domain.LogQuery
		}

		res, err := pkg.Execute(entry, logType, clusterPeer)
		if err != nil {
			return err
		}

		fmt.Println(string(res.Out))

		return nil
	},
}

func init() {
	clusterCmd.PersistentFlags().Var(&cluster, "cluster", "Live cluster peer for sending commands")
	clusterCmd.MarkPersistentFlagRequired("cluster")

	executeCmd.Flags().Var(&exType, "type", `Execute type`)
	executeCmd.MarkFlagRequired("type")

	clusterCmd.AddCommand(configurationCmd)
	clusterCmd.AddCommand(executeCmd)
	rootCmd.AddCommand(clusterCmd)
	rootCmd.AddCommand(leadershipTransferCmd)
}
