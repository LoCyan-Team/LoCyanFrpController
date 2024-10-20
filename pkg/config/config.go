package config

import (
	"gopkg.in/ini.v1"
	"log"
)

type Config struct {
	Addr string
}

func ReadCfg() *Config {
	cfg, err := ini.Load("config.ini")

	if err != nil {
		log.Fatalf("Failed to read Config File! err: %v", err)
	}

	connectInfo := cfg.Section("connection")
	addr := connectInfo.Key("addr").String()

	config := new(Config)
	config.Addr = addr
	return config
}
