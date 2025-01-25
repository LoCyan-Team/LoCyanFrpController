package config

import (
	"go.uber.org/zap"
	"gopkg.in/ini.v1"
	"lcf-controller/logger"
	"time"
)

type Config struct {
	ControllerConfig ControllerConfig
	FrpServerConfig  FrpServerConfig
	MonitorConfig    MonitorConfig
}

type ControllerConfig struct {
	Addr         string
	SendDuration time.Duration
	NodeId       int
	NodeApiKey   string
}

type FrpServerConfig struct {
	Host      string
	AdminPort int
	Username  string
	Password  string
}

type MonitorConfig struct {
	Name          string
	Addr          string
	NetworkDevice string
	AuthSecret    string
}

func ReadCfg() *Config {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		logger.Logger.Fatal("read config file failed", zap.Error(err))
	}

	config := new(Config)

	// Controller
	commonInfo := cfg.Section("common")
	controllerCfg := new(ControllerConfig)
	controllerCfg.Addr = commonInfo.Key("addr").String()
	sendDuration, err := commonInfo.Key("send_duration").Int()
	if err != nil {
		logger.Logger.Fatal("parse config file failed", zap.Error(err))
	}
	controllerCfg.SendDuration = time.Duration(sendDuration) * time.Second
	nodeId, err := commonInfo.Key("node_id").Int()
	if err != nil {
		logger.Logger.Fatal("parse config file failed", zap.Error(err))
	}
	controllerCfg.NodeId = nodeId
	controllerCfg.NodeApiKey = commonInfo.Key("node_api_key").String()

	// Frp Server
	frpsInfo := cfg.Section("frps")
	frpServerCfg := new(FrpServerConfig)
	adminPort, err := frpsInfo.Key("admin_port").Int()
	if err != nil {
		logger.Logger.Fatal("parse config file failed", zap.Error(err))
	}
	frpServerCfg.Host = frpsInfo.Key("host").String()
	frpServerCfg.AdminPort = adminPort
	frpServerCfg.Username = frpsInfo.Key("username").String()
	frpServerCfg.Password = frpsInfo.Key("password").String()

	// Monitor
	monitorInfo := cfg.Section("monitor")
	monitorCfg := new(MonitorConfig)
	monitorCfg.Name = monitorInfo.Key("name").String()
	monitorCfg.Addr = monitorInfo.Key("addr").String()
	monitorCfg.NetworkDevice = monitorInfo.Key("network_device").String()
	monitorCfg.AuthSecret = monitorInfo.Key("auth_secret").String()

	config.ControllerConfig = *controllerCfg
	config.FrpServerConfig = *frpServerCfg
	config.MonitorConfig = *monitorCfg
	return config
}
