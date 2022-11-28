package cmd

import (
	"fmt"

	"graft/pkg"
	"graft/pkg/domain"

	"github.com/spf13/cobra"
)

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

func init() {
	clusterCmd.AddCommand(configurationCmd)
}
