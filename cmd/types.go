package cmd

import (
	"fmt"
	"net/netip"
)

type ipAddr struct{ netip.AddrPort }

func (i ipAddr) String() string {
	return fmt.Sprintf("%v:%v", i.Addr(), i.Port())
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
	return "ipAddr"
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
	return "logLevel"
}
