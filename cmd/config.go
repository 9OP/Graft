package cmd

import (
	"os"

	yaml "gopkg.in/yaml.v3"
)

type Config struct {
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

func GetConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
