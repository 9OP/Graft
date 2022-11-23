package cmd

import (
	"fmt"
	"net"
	"net/netip"

	"github.com/spf13/cobra"
)

// rename start

/*
usage:

# Start a new cluster
graft start [host] --rpc-port <rpc-port> --api-port <api-port> --config <config-path>

# Add node to existing cluster
graft cluster add [host] --rpc-port <rpc-port> --api-port <api-port> --node


*/

var (
	rpcPort uint16
	cfPath  string
	cluster ipAddr
)

type ipAddr struct{ netip.AddrPort }

func (i ipAddr) String() string {
	return fmt.Sprintf("%v:%v", i.Addr(), i.Port())
}

func (i ipAddr) Set(v string) error {
	_, err := netip.ParseAddrPort(v)
	return err
}

func (i ipAddr) Type() string {
	return "ipAddr"
}

var mainCmd = &cobra.Command{
	Use:   "main",
	Short: "Start a cluster node",
	Long: `Start a cluster node:

	This command allows to:
	- starts a new cluster (start the first node)
	- start a new node and append an existing cluster config
	`,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactArgs(1)(cmd, args); err != nil {
			return err
		}

		if ip := net.ParseIP(args[0]); ip == nil {
			return fmt.Errorf("invalid host (IPv4|IPv6) %s", args[0])
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		// If cluster is given:
		// start as node and ignore config

		cmd.Help()
	},
}

func init() {
	mainCmd.Flags().VarP(&cluster, "cluster", "l", "Cluster addr")
	mainCmd.Flags().Uint16VarP(&rpcPort, "port", "p", 8080, "Node port")
	mainCmd.Flags().StringVarP(&cfPath, "config", "c", "conf/graft-config.yml", "Configuration file path")

	// Add log-level

	rootCmd.AddCommand(mainCmd)
}
