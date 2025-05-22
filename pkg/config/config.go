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
	OpenGFWConfig    OpenGFWConfig
	MonitorConfig    MonitorConfig
}

type ControllerConfig struct {
	Enable       bool
	Addr         string
	SendDuration time.Duration
	NodeId       int
	NodeApiKey   string
}

type FrpServerConfig struct {
	Username     string
	Password     string
	AdminApiHost string
	AdminApiPort int
}

type OpenGFWConfig struct {
	Enable          bool
	ConfigFilePath  string
	RulesetFilePath string
}

type MonitorConfig struct {
	Name          string
	Enable        bool
	Addr          string
	NetworkDevice string
	AuthSecret    string
}

func ReadCfg() *Config {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		logger.Fatal("read config file failed", zap.Error(err))
	}

	config := new(Config)

	// Controller
	commonInfo := cfg.Section("common")
	controllerCfg := new(ControllerConfig)
	commonEnable, err := commonInfo.Key("enable").Bool()
	if err != nil {
		logger.Fatal("parse config file failed", zap.Error(err))
	}
	controllerCfg.Enable = commonEnable
	controllerCfg.Addr = commonInfo.Key("addr").String()
	sendDuration, err := commonInfo.Key("send_duration").Int()
	if err != nil {
		logger.Fatal("parse config file failed", zap.Error(err))
	}
	controllerCfg.SendDuration = time.Duration(sendDuration) * time.Second
	nodeId, err := commonInfo.Key("node_id").Int()
	if err != nil {
		logger.Fatal("parse config file failed", zap.Error(err))
	}
	controllerCfg.NodeId = nodeId
	controllerCfg.NodeApiKey = commonInfo.Key("node_api_key").String()

	// Frp Server
	frpsInfo := cfg.Section("frps")
	frpServerCfg := new(FrpServerConfig)
	frpServerCfg.Username = frpsInfo.Key("username").String()
	frpServerCfg.Password = frpsInfo.Key("password").String()
	frpServerCfg.AdminApiHost = frpsInfo.Key("admin_api_host").String()
	adminApiPort, err := frpsInfo.Key("admin_api_port").Int()
	if err != nil {
		logger.Fatal("parse config file failed", zap.Error(err))
	}
	frpServerCfg.AdminApiPort = adminApiPort

	// OpenGFW
	opengfwInfo := cfg.Section("opengfw")
	opengfwCfg := new(OpenGFWConfig)
	opengfwEnable, err := opengfwInfo.Key("enable").Bool()
	if err != nil {
		logger.Fatal("parse config file failed", zap.Error(err))
	}
	opengfwCfg.Enable = opengfwEnable
	opengfwCfg.ConfigFilePath = opengfwInfo.Key("config_file").String()
	opengfwCfg.RulesetFilePath = opengfwInfo.Key("ruleset_file").String()

	// Monitor
	monitorInfo := cfg.Section("monitor")
	monitorCfg := new(MonitorConfig)
	monitorCfg.Name = monitorInfo.Key("name").String()
	monitorEnable, err := monitorInfo.Key("enable").Bool()
	if err != nil {
		logger.Fatal("parse config file failed", zap.Error(err))
	}
	monitorCfg.Enable = monitorEnable
	monitorCfg.Addr = monitorInfo.Key("addr").String()
	monitorCfg.NetworkDevice = monitorInfo.Key("network_device").String()
	monitorCfg.AuthSecret = monitorInfo.Key("auth_secret").String()

	config.ControllerConfig = *controllerCfg
	config.FrpServerConfig = *frpServerCfg
	config.OpenGFWConfig = *opengfwCfg
	config.MonitorConfig = *monitorCfg
	return config
}
