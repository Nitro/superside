package main

import (
	"os"

	"github.com/BurntSushi/toml"
	log "github.com/Sirupsen/logrus"
)

type Config struct {
	Superside *ApiConfig       `toml:"superside"`
}

type ApiConfig struct {
	BindIP       string `toml:"bind_ip"`
	BindPort     int    `toml:"bind_port"`
	LoggingLevel string `toml:"logging_level"`
}

func parseConfig(path string) *Config {
	var config Config
	_, err := toml.DecodeFile(path, &config)
	if err != nil {
		exitWithError(err, "Failed to parse config file")
	}

	proxy := config.Superside
	if proxy == nil {
		log.Error("Missing 'haproxy' section of config file")
		os.Exit(1)
	}

	if config.Superside.BindIP == "" {
		config.Superside.BindIP = "0.0.0.0"
	}

	if config.Superside.BindPort == 0 {
		config.Superside.BindPort = 7778
	}

	configureLoggingLevel(config.Superside.LoggingLevel)

	return &config
}

func configureLoggingLevel(level string) {
	switch {
	case len(level) == 0:
		log.SetLevel(log.InfoLevel)
	case level == "info":
		log.SetLevel(log.InfoLevel)
	case level == "warn":
		log.SetLevel(log.WarnLevel)
	case level == "error":
		log.SetLevel(log.ErrorLevel)
	case level == "debug":
		log.SetLevel(log.DebugLevel)
	}
}
