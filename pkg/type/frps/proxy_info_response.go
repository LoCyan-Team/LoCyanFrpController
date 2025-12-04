package frps

type TunnelInfoResponse struct {
	Name            string `json:"name"`
	Conf            *Conf  `json:"conf"`
	ClientVersion   string `json:"clientVersion"`   // 修正：clientVersion
	TodayTrafficIn  int    `json:"todayTrafficIn"`  // 修正：todayTrafficIn
	TodayTrafficOut int    `json:"todayTrafficOut"` // 修正：todayTrafficOut
	CurConns        int    `json:"curConns"`        // 修正：curConns
	LastStartTime   string `json:"lastStartTime"`   // 修正：lastStartTime
	LastCloseTime   string `json:"lastCloseTime"`   // 修正：lastCloseTime
	Status          string `json:"status"`          // 修正：status
}

type Tunnel struct {
	Tunnels []TunnelInfoResponse `json:"proxies"`
}

// 定义一个结构体来映射配置信息
type Conf struct {
	Name         string        `json:"name"`
	Type         string        `json:"type"`
	Transport    *Transport    `json:"transport"`    // 新增：transport 字段
	LoadBalancer *LoadBalancer `json:"loadBalancer"` // 新增：loadBalancer 字段
	HealthCheck  *HealthCheck  `json:"healthCheck"`  // 新增：healthCheck 字段
	LocalIP      string        `json:"localIP"`
	Plugin       interface{}   `json:"plugin"`
	RemotePort   int           `json:"remotePort"`

	// 可选字段（根据实际需要保留）
	UseEncryption        bool        `json:"use_encryption,omitempty"`
	UseCompression       bool        `json:"use_compression,omitempty"`
	Group                string      `json:"group,omitempty"`
	GroupKey             string      `json:"group_key,omitempty"`
	ProxyProtocolVersion string      `json:"proxy_protocol_version,omitempty"`
	BandwidthLimit       string      `json:"bandwidth_limit,omitempty"`
	BandwidthLimitMode   string      `json:"bandwidth_limit_mode,omitempty"`
	Metas                interface{} `json:"metas,omitempty"`
	LocalPort            int         `json:"local_port,omitempty"`
	PluginParams         interface{} `json:"plugin_params,omitempty"`
	HealthCheckType      string      `json:"health_check_type,omitempty"`
	HealthCheckTimeoutS  int         `json:"health_check_timeout_s,omitempty"`
	HealthCheckMaxFailed int         `json:"health_check_max_failed,omitempty"`
	HealthCheckIntervalS int         `json:"health_check_interval_s,omitempty"`
	HealthCheckURL       string      `json:"health_check_url,omitempty"`
	HealthCheckAddr      string      `json:"health_check_addr,omitempty"`
}

// Transport 新增：Transport 结构体
type Transport struct {
	BandwidthLimit     string `json:"bandwidthLimit"`
	BandwidthLimitMode string `json:"bandwidthLimitMode"`
}

// LoadBalancer 新增：LoadBalancer 结构体
type LoadBalancer struct {
	Group string `json:"group"`
}

// HealthCheck 新增：HealthCheck 结构体
type HealthCheck struct {
	Type            string `json:"type"`
	IntervalSeconds int    `json:"intervalSeconds"`
}
