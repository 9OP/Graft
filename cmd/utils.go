package cmd

import (
	"crypto/sha1"
	"encoding/hex"
	"net/netip"
	"os"

	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v3"
)

func validateAddrArg(cmd *cobra.Command, args []string) error {
	if err := cobra.ExactArgs(1)(cmd, args); err != nil {
		return err
	}
	if _, err := netip.ParseAddrPort(args[0]); err != nil {
		return err
	}
	return nil
}

type configuration struct {
	Fsm      string `yaml:"fsm"`
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

type ipAddr string

func (i ipAddr) String() string {
	return string(i)
}

func (i *ipAddr) Set(v string) error {
	ip, err := netip.ParseAddrPort(v)
	if err != nil {
		return err
	}
	*i = ipAddr(ip.String())
	return nil
}

func (i ipAddr) Type() string {
	return "<ip>:<port>"
}
