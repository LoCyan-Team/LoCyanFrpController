package info

import (
	"lcf-controller/net/server"
	"strings"
)

func GetProxies(proxyType string) ([]map[string]any, error) {
	proxies, err := server.GetProxyList(proxyType)
	if err != nil {
		return nil, err
	}
	proxyList := make([]map[string]any, 0)
	for _, p := range proxies.Proxies {
		tmp := make(map[string]any)
		tmp["proxy_name"] = strings.Split(p.Name, ".")[1]
		tmp["inbound"] = p.TodayTrafficIn
		tmp["outbound"] = p.TodayTrafficOut
		proxyList = append(proxyList, tmp)
	}
	return proxyList, nil
}
