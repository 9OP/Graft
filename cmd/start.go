package cmd

import (
	"fmt"
	"net/netip"
	"time"

	"graft/pkg"
	"graft/pkg/domain"
	"graft/pkg/utils/log"

	"github.com/spf13/cobra"
)

var config string

var startCmd = &cobra.Command{
	Use:     "start [ip:port]",
	GroupID: "membership",
	Short:   "Start cluster node",
	Long: `Start cluster node:

	If specified cluster flag is equal to [ip:port], then starts
	a new cluster with first node [ip:port].
	`,
	Args: validateAddrArg,
	RunE: func(cmd *cobra.Command, args []string) error {
		host, _ := netip.ParseAddrPort(args[0])
		id := hashString(host.String())

		peer := domain.Peer{
			Id:   id,
			Host: host,
		}

		log.ConfigureLogger(id)

		// Start new cluster
		if host.String() == cluster.String() {
			cf, err := loadConfiguration(config)
			if err != nil {
				return fmt.Errorf("failed load configuration\n%v", err.Error())
			}

			quit := pkg.Start(
				id,
				host,
				domain.Peers{id: peer},
				cf.Fsm,
				cf.Timeouts.Election,
				cf.Timeouts.Heartbeat,
			)

			// wait for peer to upgrade leader
			time.Sleep(3 * time.Second)

			err = pkg.AddSelf(peer)
			if err != nil {
				return err
			}

			// wait
			<-quit
		} else {
			// Add new peer to cluster
			quit, err := pkg.AddClusterPeer(peer, cluster.String())
			if err != nil {
				return fmt.Errorf("failed add peer\n%v", err.Error())
			}

			// wait
			<-quit

		}

		return nil
	},
}

func init() {
	startCmd.Flags().StringVar(&config, "config", "conf/graft-config.yml", "Configuration file path")
	rootCmd.AddCommand(startCmd)
}
