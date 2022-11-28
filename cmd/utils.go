package cmd

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net/netip"
	"os"

	yaml "gopkg.in/yaml.v3"
)

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

type execType string

func (e execType) String() string {
	return "<nil>"
}

func (e *execType) Set(v string) error {
	switch v {
	case "QUERY", "COMMAND":
		*e = execType(v)
		return nil
	default:
		return fmt.Errorf("\n\tmust be one of 'COMMAND' or 'QUERY'")
	}
}

func (e execType) Type() string {
	return `COMMAND|QUERY`
}

type ipAddr struct{ netip.AddrPort }

func (i ipAddr) String() string {
	return "<nil>"
}

func (i *ipAddr) Set(v string) error {
	ip, err := netip.ParseAddrPort(v)
	if err != nil {
		return err
	}
	i.AddrPort = ip
	return nil
}

func (i ipAddr) Type() string {
	return "<ip>:<port>"
}

type logLevel string

const (
	DEBUG logLevel = "DEBUG"
	INFO  logLevel = "INFO"
	ERROR logLevel = "ERROR"
)

func (l logLevel) String() string {
	return string(l)
}

func (l *logLevel) Set(v string) error {
	switch v {
	case "DEBUG", "INFO", "ERROR":
		*l = logLevel(v)
		return nil
	default:
		return fmt.Errorf("\n\tmust be one of 'DEBUG', 'INFO', or 'ERROR'")
	}
}

func (l logLevel) Type() string {
	return "[level]"
}
