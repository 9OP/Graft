package utils

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func ConfigureLogger(logLevel string) {
	var level log.Level = log.InfoLevel
	switch logLevel {
	case "DEBUG":
		level = log.DebugLevel
	case "INFO":
		level = log.InfoLevel
	case "ERROR":
		level = log.ErrorLevel
	}

	log.SetLevel(level)
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.TextFormatter{
		ForceColors:     true,
		DisableColors:   false,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05.000",
	})
}
