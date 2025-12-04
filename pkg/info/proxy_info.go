package info

import (
	"lcf-controller/net/server"
	"lcf-controller/pkg/config"
	"strings"
)

func GetTunnelInfo(cfg *config.Config, tunnelType string) ([]map[string]any, error) {
	tunnels, err := server.GetProxyList(tunnelType)
	if err != nil {
		return nil, err
	}
	tunnelList := make([]map[string]any, 0)
	for _, p := range tunnels.Tunnels {
		tmp := make(map[string]any)
		tmp["node_id"] = cfg.ControllerConfig.NodeId
		tmp["tunnel_name"] = strings.Split(p.Name, ".")[1]
		tmp["inbound_traffic"] = p.TodayTrafficIn
		tmp["outbound_traffic"] = p.TodayTrafficOut
		tunnelList = append(tunnelList, tmp)
	}
	return tunnelList, nil
}
