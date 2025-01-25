package frps

type ServerInfoResponse struct {
	Version               string         `json:"version"`
	BindPort              int            `json:"bind_port"`
	BindUDPPort           int            `json:"bind_udp_port"`
	VhostHTTPPort         int            `json:"vhost_http_port"`
	VhostHTTPSPort        int            `json:"vhost_https_port"`
	TCPMUXHTTPConnectPort int            `json:"tcpmux_httpconnect_port"`
	KCPBindPort           int            `json:"kcp_bind_port"`
	QUICBindPort          int            `json:"quic_bind_port"`
	SubdomainHost         string         `json:"subdomain_host"`
	MaxPoolCount          int            `json:"max_pool_count"`
	MaxPortsPerClient     int            `json:"max_ports_per_client"`
	HeartBeatTimeout      int            `json:"heart_beat_timeout"`
	TotalTrafficIn        int64          `json:"total_traffic_in"`
	TotalTrafficOut       int64          `json:"total_traffic_out"`
	CurConns              int            `json:"cur_conns"`
	ClientCounts          int            `json:"client_counts"`
	ProxyTypeCount        ProxyTypeCount `json:"proxy_type_count"`
}

type ProxyTypeCount struct {
	TCP   int `json:"tcp"`
	UDP   int `json:"udp"`
	HTTP  int `json:"http"`
	HTTPS int `json:"https"`
	XTCP  int `json:"xtcp"`
	STCP  int `json:"stcp"`
}
