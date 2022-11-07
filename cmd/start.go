package cmd

import (
	"fmt"
	"graft/pkg"
	"graft/pkg/domain/entity"
	"os"

	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v3"
)

type configuration struct {
	Fsm      string `yaml:"fsm"`
	Timeouts struct {
		Election  int `yaml:"election"`
		Heartbeat int `yaml:"heartbeat"`
	}
	Peers map[string]struct {
		Host  string `yaml:"host"`
		Ports struct {
			P2p string `yaml:"p2p"`
			Api string `yaml:"api"`
		} `yaml:"ports"`
	} `yaml:"peers"`
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

type logLevel string

const (
	DEBUG logLevel = "DEBUG"
	INFO  logLevel = "INFO"
	ERROR logLevel = "ERROR"
)

func (l *logLevel) String() string {
	return string(*l)
}
func (l *logLevel) Set(v string) error {
	switch v {
	case "DEBUG", "INFO", "ERROR":
		*l = logLevel(v)
		return nil
	default:
		return fmt.Errorf("must be one of 'DEBUG', 'INFO', or 'ERROR'")
	}
}
func (l *logLevel) Type() string {
	return "logLevel"
}

var (
	level      = INFO
	configPath string
	config     *configuration
	peers      entity.Peers
)

var startCmd = &cobra.Command{
	Use:   "start [peer]",
	Short: "Start a cluster node",
	Args: func(cmd *cobra.Command, args []string) error {
		// Custom args validator
		// ensures that passed `peer` is declared
		// in the configuration file.

		configPath, _ := cmd.Flags().GetString("config")
		if err := cobra.ExactArgs(1)(cmd, args); err != nil {
			return err
		}

		id := args[0]
		c, err := loadConfiguration(configPath)
		if err != nil {
			return err
		}

		config = c
		if _, ok := config.Peers[id]; ok {
			// Validate peers format
			validPeers := entity.Peers{}
			for peerId, peer := range config.Peers {
				validPeer, err := entity.NewPeer(peerId, peer.Host, peer.Ports.P2p, peer.Ports.Api)
				if err != nil {
					return err
				}
				if peerId == id {
					continue
				}
				validPeers[peerId] = *validPeer
			}
			peers = validPeers
			return nil
		}

		peers := make([]string, 0, len(config.Peers))
		for k := range config.Peers {
			peers = append(peers, k)
		}
		return fmt.Errorf(" invalid peer: %s\n\tauthorised peers: %s", id, peers)
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		persistentLocation := fmt.Sprintf("persistent_%s.json", id)

		pkg.Start(
			id,
			peers,
			persistentLocation,
			config.Timeouts.Election,
			config.Timeouts.Heartbeat,
			config.Peers[id].Ports.P2p,
			config.Peers[id].Ports.Api,
			config.Fsm,
			level.String(),
		)
	},
}

func init() {
	startCmd.Flags().VarP(&level, "log-level", "l", `log level. allowed: "DEBUG", "INFO", "ERROR"`)
	startCmd.Flags().StringVarP(&configPath, "config", "c", "graft-config.yml", "Configuration file path")
	rootCmd.AddCommand(startCmd)
}
