package info

import (
	"lcf-controller/net/server"
	"lcf-controller/pkg/config"
	"strings"
)

func GetProxies(cfg *config.Config, proxyType string) ([]map[string]any, error) {
	proxies, err := server.GetProxyList(proxyType)
	if err != nil {
		return nil, err
	}
	proxyList := make([]map[string]any, 0)
	for _, p := range proxies.Proxies {
		tmp := make(map[string]any)
		tmp["node_id"] = cfg.ControllerConfig.NodeId
		tmp["proxy_name"] = strings.Split(p.Name, ".")[1]
		tmp["in_bound_traffic"] = p.TodayTrafficIn
		tmp["out_bound_traffic"] = p.TodayTrafficOut
		proxyList = append(proxyList, tmp)
	}
	return proxyList, nil
}
