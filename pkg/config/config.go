package config

import (
	"gopkg.in/ini.v1"
	"log"
	"time"
)

type Config struct {
	Addr         string
	AdminPort    int
	Username     string
	Password     string
	SendDuration time.Duration
	NodeId       int
	NodeApiKey   string
}

func ReadCfg() *Config {
	cfg, err := ini.Load("config.ini")

	if err != nil {
		log.Fatalf("Failed to read Config File! err: %v", err)
	}

	commonInfo := cfg.Section("common")
	sendDuration, err := commonInfo.Key("send_duration").Int()
	if err != nil {
		log.Fatalf("Parse config file failed!, err: %s", err)
	}
	nodeId, err := commonInfo.Key("nodeId").Int()
	if err != nil {
		log.Fatalf("Parse config file failed!, err: %s", err)
	}
	nodeApiKey := commonInfo.Key("nodeApiKey").String()

	connectInfo := cfg.Section("connection")
	addr := connectInfo.Key("addr").String()
	username := connectInfo.Key("username").String()
	password := connectInfo.Key("password").String()
	adminPort, err := connectInfo.Key("admin_port").Int()
	if err != nil {
		log.Fatalf("Parse config file failed!, err: %s", err)
	}

	config := new(Config)
	config.Addr = addr
	config.Username = username
	config.Password = password
	config.AdminPort = adminPort
	config.SendDuration = time.Duration(sendDuration)
	config.NodeId = nodeId
	config.NodeApiKey = nodeApiKey
	return config
}
