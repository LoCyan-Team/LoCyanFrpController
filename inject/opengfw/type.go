package opengfw

import "time"

type CliConfigRuleset struct {
	GeoIp   string `mapstructure:"geoip"`
	GeoSite string `mapstructure:"geosite"`
}

type CliConfig struct {
	IO      CliConfigIO      `mapstructure:"io"`
	Workers CliConfigWorkers `mapstructure:"workers"`
	Ruleset CliConfigRuleset `mapstructure:"ruleset"`
	Replay  CliConfigReplay  `mapstructure:"replay"`
}

type CliConfigIO struct {
	QueueSize      uint32  `mapstructure:"queueSize"`
	QueueNum       *uint16 `mapstructure:"queueNum"`
	Table          string  `mapstructure:"table"`
	ConnMarkAccept uint32  `mapstructure:"connMarkAccept"`
	ConnMarkDrop   uint32  `mapstructure:"connMarkDrop"`

	ReadBuffer  int  `mapstructure:"rcvBuf"`
	WriteBuffer int  `mapstructure:"sndBuf"`
	Local       bool `mapstructure:"local"`
	RST         bool `mapstructure:"rst"`
}

type CliConfigReplay struct {
	Realtime bool `mapstructure:"realtime"`
}

type CliConfigWorkers struct {
	Count                      int           `mapstructure:"count"`
	QueueSize                  int           `mapstructure:"queueSize"`
	TCPMaxBufferedPagesTotal   int           `mapstructure:"tcpMaxBufferedPagesTotal"`
	TCPMaxBufferedPagesPerConn int           `mapstructure:"tcpMaxBufferedPagesPerConn"`
	TCPTimeout                 time.Duration `mapstructure:"tcpTimeout"`
	UDPMaxStreams              int           `mapstructure:"udpMaxStreams"`
}
