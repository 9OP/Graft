package cmd

import (
	"fmt"
	"net/netip"

	"graft/pkg"
	"graft/pkg/domain"
	"graft/pkg/utils/log"

	"github.com/spf13/cobra"
)

var config string

var startCmd = &cobra.Command{
	Use:   "start [ip:port]",
	Short: "Start a new cluster",
	Args:  validateAddrArg,
	RunE: func(cmd *cobra.Command, args []string) error {
		host, _ := netip.ParseAddrPort(args[0])
		id := hashString(host.String())
		log.ConfigureLogger(id)

		cf, err := loadConfiguration(config)
		if err != nil {
			return fmt.Errorf("failed load configuration\n%v", err.Error())
		}

		peers := domain.Peers{
			id: domain.Peer{
				Id:     id,
				Host:   host,
				Active: true,
			},
		}

		quit := pkg.Start(
			id,
			host,
			peers,
			cf.Fsm,
			cf.Timeouts.Election,
			cf.Timeouts.Heartbeat,
		)

		// wait
		<-quit
		return nil
	},
}

func init() {
	startCmd.Flags().StringVarP(&config, "config", "c", "conf/graft-config.yml", "Configuration file path")
	rootCmd.AddCommand(startCmd)
}
