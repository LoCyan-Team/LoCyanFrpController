package config

import (
	"go.uber.org/zap"
	"gopkg.in/ini.v1"
	"lcf-controller/logger"
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
		logger.Logger.Fatal("read config file failed", zap.Error(err))
	}

	commonInfo := cfg.Section("common")
	sendDuration, err := commonInfo.Key("send_duration").Int()
	if err != nil {
		logger.Logger.Fatal("parse config file failed", zap.Error(err))
	}
	nodeId, err := commonInfo.Key("node_id").Int()
	if err != nil {
		logger.Logger.Fatal("parse config file failed", zap.Error(err))
	}
	nodeApiKey := commonInfo.Key("node_api_key").String()

	connectInfo := cfg.Section("connection")
	addr := connectInfo.Key("addr").String()
	username := connectInfo.Key("username").String()
	password := connectInfo.Key("password").String()
	adminPort, err := connectInfo.Key("admin_port").Int()
	if err != nil {
		logger.Logger.Fatal("parse config file failed", zap.Error(err))
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
