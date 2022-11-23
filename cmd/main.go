package cmd

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"graft/pkg"
	"graft/pkg/domain"
	"net/netip"
	"os"

	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v3"
)

/*
usage:

# Start a new cluster
graft start [host] --rpc-port <rpc-port> --api-port <api-port> --config <config-path>

# Add node to existing cluster
graft cluster add [host] --rpc-port <rpc-port> --api-port <api-port> --node


*/

var (
	level   = INFO
	config  string
	cluster ipAddr
)

type configuration struct {
	Fsm struct {
		Eval string `yaml:"eval"`
		Init string `yaml:"init"`
	}
	Timeouts struct {
		Election  int `yaml:"election"`
		Heartbeat int `yaml:"heartbeat"`
	}
}

func loadConfiguration(path string) (*configuration, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg configuration
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

var startCmd = &cobra.Command{
	Use:   "start",
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

		if _, err := netip.ParseAddrPort(args[0]); err != nil {
			return err
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		isClusterProvided := cmd.Flags().Changed("cluster")

		host, err := netip.ParseAddrPort(args[0])
		if err != nil {
			return err
		}

		if isClusterProvided {
			fmt.Println("add node to cluster")
			// Get cluster config:
			// - timeouts
			// - peers
			// Update cluster config
			// Start peer
		} else {
			config, _ := cmd.Flags().GetString("config")

			cf, err := loadConfiguration(config)
			if err != nil {
				return err
			}

			hasher := sha1.New()
			hasher.Write([]byte(host.String()))
			id := hex.EncodeToString(hasher.Sum(nil))

			pkg.Start(
				id,
				host,
				domain.Peers{},
				fmt.Sprintf("conf/%s.json", id), // persistent location
				cf.Timeouts.Election,
				cf.Timeouts.Heartbeat,
				level.String(),
			)
		}

		return nil
	},
}

func init() {
	startCmd.Flags().Var(&cluster, "cluster", "Cluster addr")
	startCmd.Flags().StringVarP(&config, "config", "c", "conf/graft-config.yml", "Configuration file path")
	startCmd.Flags().Var(&level, "log", `log level. allowed: "DEBUG", "INFO", "ERROR"`)

	rootCmd.AddCommand(startCmd)
}
