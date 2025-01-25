package frps

type ProxyInfoResponse struct {
	Name            string `json:"name"`
	Conf            *Conf  `json:"conf"`
	ClientVersion   string `json:"client_version"`
	TodayTrafficIn  int    `json:"today_traffic_in"`
	TodayTrafficOut int    `json:"today_traffic_out"`
	CurConns        int    `json:"cur_conns"`
	LastStartTime   string `json:"last_start_time"`
	LastCloseTime   string `json:"last_close_time"`
	Status          string `json:"status"`
}

type Proxy struct {
	Proxies []ProxyInfoResponse `json:"proxies"`
}

// 定义一个结构体来映射配置信息
type Conf struct {
	Name                 string      `json:"name"`
	Type                 string      `json:"type"`
	UseEncryption        bool        `json:"use_encryption"`
	UseCompression       bool        `json:"use_compression"`
	Group                string      `json:"group"`
	GroupKey             string      `json:"group_key"`
	ProxyProtocolVersion string      `json:"proxy_protocol_version"`
	BandwidthLimit       string      `json:"bandwidth_limit"`
	BandwidthLimitMode   string      `json:"bandwidth_limit_mode"`
	Metas                interface{} `json:"metas"`
	LocalIP              string      `json:"local_ip"`
	LocalPort            int         `json:"local_port"`
	Plugin               string      `json:"plugin"`
	PluginParams         interface{} `json:"PluginParams"`
	HealthCheckType      string      `json:"health_check_type"`
	HealthCheckTimeoutS  int         `json:"health_check_timeout_s"`
	HealthCheckMaxFailed int         `json:"health_check_max_failed"`
	HealthCheckIntervalS int         `json:"health_check_interval_s"`
	HealthCheckURL       string      `json:"health_check_url"`
	HealthCheckAddr      string      `json:"HealthCheckAddr"`
	RemotePort           int         `json:"remote_port"`
}
