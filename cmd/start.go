package cmd

import (
	"crypto/sha1"
	"encoding/hex"
	"net/netip"
	"os"

	"graft/pkg"
	"graft/pkg/domain"
	"graft/pkg/utils"

	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v3"
)

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

func hashString(str string) string {
	hasher := sha1.New()
	hasher.Write([]byte(str))
	return hex.EncodeToString(hasher.Sum(nil))
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
		utils.ConfigureLogger(level.String())
		host, _ := netip.ParseAddrPort(args[0])
		id := hashString(host.String())
		isClusterProvided := cmd.Flags().Changed("cluster")

		if isClusterProvided {
			clusterPeer := domain.Peer{Host: cluster.AddrPort}
			newPeer := domain.Peer{Id: id, Host: host, Active: false}

			quit, err := pkg.AddClusterPeer(newPeer, clusterPeer)
			if err != nil {
				return err
			}

			// wait
			<-quit
			return nil

		} else {
			cf, err := loadConfiguration(config)
			if err != nil {
				return err
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
				cf.Timeouts.Election,
				cf.Timeouts.Heartbeat,
			)

			// wait
			<-quit
			return nil
		}
	},
}

func init() {
	startCmd.Flags().Var(&cluster, "cluster", "Cluster addr")
	startCmd.Flags().StringVarP(&config, "config", "c", "conf/graft-config.yml", "Configuration file path")
	startCmd.Flags().Var(&level, "log", `log level. allowed: "DEBUG", "INFO", "ERROR"`)

	rootCmd.AddCommand(startCmd)
}
